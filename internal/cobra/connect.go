/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cobra

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kaitokid2302/broadcast-server/internal/utils"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.PortInUse(":8080") {
			fmt.Printf("\"Server not started yet\": %v\n", "Server not started yet")
			return
		}
		url := "ws://localhost:8080"
		con, _, er := websocket.DefaultDialer.Dial(url, nil)
		if er != nil {
			return
		}
		// read
		go func() {
			for {
				_, message, er := con.ReadMessage()
				if er != nil {
					fmt.Println("Server close, sorry, service ended!")
					os.Exit(0)
				}
				fmt.Printf("message: %v\n", string(message))
			}
		}()
		// write
		go func() {
			cin := bufio.NewScanner(os.Stdin)
			for {
				// fmt.Println("What do you want to send?:")
				cin.Scan()
				var message string = cin.Text()
				if message == "exit" {
					time.Sleep(time.Second)
					fmt.Println("Ok bye then.")
					os.Exit(0)
				}
				con.WriteMessage(1, []byte(message))
			}
		}()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
