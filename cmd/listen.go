/*
Copyright © 2024 Sebastian Ibanez <sebas.ibanez219@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sebastian-j-ibanez/fsync/client"
	dir "github.com/sebastian-j-ibanez/fsync/directory"
	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for sync requests",
	Long:  `Listens for a sync request over a socket connection on port 2000.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		portFlag, _ := cmd.Flags().GetString("port")
		scanFlag, _ := cmd.Flags().GetBool("scan")

		// Handle port flag
		port := 8080
		var err error
		if portFlag != "" {
			port, err = strconv.Atoi(portFlag)
			if port <= 0 {
				fmt.Fprintf(os.Stderr, "error: port must be greater than 0\n")
				os.Exit(-1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: cannot convert arg to port\n")
				os.Exit(-1)
			}
		}

		// Broadcast MDNS service
		endBroadcast := make(chan bool)
		if scanFlag {
			go func() {
				if err := client.BroadcastMDNSService(port, endBroadcast); err != nil {
					fmt.Fprintf(os.Stderr, "error: %s", err.Error())
					os.Exit(-1)
				}
			}()
		}

		// Init dir manager
		path, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(-1)
		}

		d, err := dir.NewDirManager(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(-1)
		}

		// Init client
		c := client.Client{
			DirMan: *d,
		}

		// Await sync
		err = c.AwaitSync(port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(-1)
		}

		// End MDNS service broadcast
		if scanFlag {
			endBroadcast <- true
		}

		fmt.Printf("Sync completed successfully!\n")
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.PersistentFlags().BoolP("scan", "s", false, "scan network for peer")
	listenCmd.PersistentFlags().StringP("port", "p", "2000", "specify the port")
}
