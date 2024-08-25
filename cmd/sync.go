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
		// Handle peer flag
		var peer prot.Peer
		var err error
		peerFlag, _ := cmd.Flags().GetString("peer")
		if peerFlag != "" {
			if ip := net.ParseIP(peerFlag); ip != nil {
				peer.IP = ip.String()
				peer.Port = "2000"
			} else {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: invalid peer ip: %s\n", peerFlag)
				os.Exit(-1)
			}
		}

		// Init directory manager
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

		// Init client
		c := client.Client{
			DirMan: *d,
			Peers:  []prot.Peer{peer},
		}

		// If peer is initialized, init sync with peer.
		// Otherwise use resgistered peer list.
		if peer == (prot.Peer{}) {
			err = c.InitSync()
			if err != nil {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
				os.Exit(-1)
			}
		} else {
			err = c.InitSyncWithPeer(peer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fsync: ERROR: %v\n", err)
				os.Exit(-1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.PersistentFlags().String("peer", "", "specify IP for sync")
}
