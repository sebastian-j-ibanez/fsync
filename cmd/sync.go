/*
Copyright © 2024 Sebastian Ibanez <sebas.ibanez219@gmail.com>
*/
package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/sebastian-j-ibanez/fsync/client"
	dir "github.com/sebastian-j-ibanez/fsync/directory"
	prot "github.com/sebastian-j-ibanez/fsync/protocol"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync files with peer clients",
	Long: `Send a sync request to peer clients.
Uses list of peers unless port flag is specified.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init directory manager amd client
		path, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
			os.Exit(-1)
		}
		d, err := dir.NewDirManager(path + "/img")
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
			os.Exit(-1)
		}
		c := client.Client{
			DirMan: *d,
		}

		peerFlag, _ := cmd.Flags().GetString("peer")
		scanFlag, _ := cmd.Flags().GetString("scan")

		// Flag cases:
		// 1. Peer flag (use specific ip)
		// 2. Scan flag (scan local network for peer)
		// 3. Neither (use client list saved in JSON)
		if peerFlag != "" {
			// Validate ip argument
			var peer prot.Peer
			if ip := net.ParseIP(peerFlag); ip != nil {
				peer.IP = ip.String()
				peer.Port = "2000"
			} else {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: invalid peer ip: %s\n", peerFlag)
				os.Exit(-1)
			}

			// Set client peers
			c.Peers = []prot.Peer{peer}
		} else if scanFlag != "" {
			// Scan for peer
			peer, err := client.FindService()
			if err != nil {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
				os.Exit(-1)
			}

			c.Peers = append(c.Peers, peer)
		} else {
			c.Peers, err = prot.GetPeers()
			if err != nil {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: unable to get peers: %v\n", err)
				os.Exit(-1)
			}
		}

		// Init sync
		err = c.InitSync()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
			os.Exit(-1)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.PersistentFlags().String("peer", "", "specify IP for sync")
	syncCmd.PersistentFlags().String("scan", "", "scan network for listening peer")
	syncCmd.MarkFlagsMutuallyExclusive("peer", "scan")
}
