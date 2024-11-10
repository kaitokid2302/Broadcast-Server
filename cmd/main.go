/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kaitokid2302/broadcast-server/internal/cobra"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	godotenv.Load()
	cobra.Execute()
	wg.Wait()
}
