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
	Short: "Listen over a socket connection",
	Long: `Listen over a socket connection for a sync request.
Listens over all network interfaces on port 2000 by default.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Handle port flag
		var port int
		var err error
		portFlag, _ := cmd.Flags().GetString("port")
		if portFlag != "" {
			port, err = strconv.Atoi(portFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: Invalid port: %s\n", portFlag)
				os.Exit(-1)
			}
		} else {
			port = 2000
		}

		// Init dir manager and client
		path, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
			os.Exit(-1)
		}

		d, err := dir.NewDirManager(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
			os.Exit(-1)
		}

		c := client.Client{
			DirMan: *d,
		}

		// Await sync
		fmt.Printf("listening over port %d...\n", port)
		err = c.AwaitSync()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
			os.Exit(-1)
		}
		fmt.Printf("sync completed successfully!\n")
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.PersistentFlags().String("port", "", "specify the port (default: 2000)")
}