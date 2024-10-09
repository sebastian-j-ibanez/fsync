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
		// Get flags
		portFlag, _ := cmd.Flags().GetString("port")
		scanFlag, _ := cmd.Flags().GetBool("scan")

		// Handle port flag
		port := 2000
		var err error
		if portFlag != "" {
			port, err = strconv.Atoi(portFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: cannot convert arg to port\n")
				os.Exit(-1)
			}
		}

		// Handle scan flag
		endBroadcast := make(chan bool)
		if scanFlag {
			go func() {
				if err := client.BroadcastMDNSService(port, endBroadcast); err != nil {
					fmt.Fprintf(os.Stderr, "fsync: ERROR: %s", err.Error())
					os.Exit(-1)
				}
			}()
		}

		// Init dir manager
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

		// Init client
		c := client.Client{
			DirMan: *d,
		}

		// Await sync
		fmt.Printf("listening over port %d...\n", port)
		err = c.AwaitSync(port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
			os.Exit(-1)
		}

		// End broadcast
		endBroadcast <- true

		fmt.Printf("sync completed successfully!\n")
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.PersistentFlags().String("port", "", "specify the port (default: 2000)")
	listenCmd.PersistentFlags().Bool("scan", false, "scan network for peer")
}
