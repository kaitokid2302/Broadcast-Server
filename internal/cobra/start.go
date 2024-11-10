/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kaitokid2302/broadcast-server/internal/utils"
	"github.com/spf13/cobra"
)

type connection struct {
	con   *websocket.Conn
	mutex sync.Mutex
}

var g map[*connection]bool = make(map[*connection]bool)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if utils.PortInUse(":8080") {
			fmt.Printf("Server already started")
			return
		}
		server := gin.Default()
		server.GET("/", handleWebsocket)
		server.Run(":8080")
	},
}

func handleWebsocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Chấp nhận tất cả các kết nối
		},
	}
	con, er := upgrader.Upgrade(c.Writer, c.Request, nil)
	if er != nil {
		return
	}
	a := &connection{con: con}
	g[a] = true
	go read(a)
}

func read(a *connection) {
	defer func() {
		a.con.Close()
		delete(g, a)
	}()
	for {
		_, message, er := a.con.ReadMessage()
		if er != nil {
			return
		}
		broadcast(string(message))
	}
}

func broadcast(message string) {
	for x := range g {
		write(x, message)
	}
}

func write(a *connection, message string) {
	a.mutex.Lock()
	er := a.con.WriteMessage(websocket.TextMessage, []byte(message))
	if er != nil {
		delete(g, a)
		a.con.Close()
	}
	a.mutex.Unlock()
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
