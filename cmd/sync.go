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
		// Init directory manager and client
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
		c := client.Client{
			DirMan: *d,
		}

		addrFlag, _ := cmd.Flags().GetString("address")
		peersFlag, _ := cmd.Flags().GetString("peers")

		// Flag cases:
		// 1. Peer flag (use specific ip)
		// 2. Peers flag (use peer list saved in JSON)
		// 3. Neither (scan local network for peer)
		if addrFlag != "" {
			// Validate ip argument
			var peer prot.Peer
			if ip := net.ParseIP(addrFlag); ip != nil {
				peer.IP = ip.String()
				peer.Port = "8001"
			} else {
				fmt.Fprintf(os.Stderr, "error: invalid peer ip: %s\n", addrFlag)
				os.Exit(-1)
			}
			c.Peers = []prot.Peer{peer}
		} else if peersFlag != "" {
			c.Peers, err = prot.GetPeers()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: unable to get peers: %v\n", err)
				os.Exit(-1)
			}
		} else {
			// Scan for peer
			peer, err := client.DiscoverMDNSService()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(-1)
			}
			c.Peers = append(c.Peers, peer)
		}

		// Init sync
		filePattern := []string{}
		if len(args) > 0 {
			filePattern = args
		}

		err = c.InitSync(filePattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(-1)
		}

		fmt.Println("Sync completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.PersistentFlags().BoolP("scan", "s", false, "scan network for peer and sync")
	syncCmd.PersistentFlags().StringP("address", "a", "", "sync with specific IP:PORT")
	syncCmd.PersistentFlags().BoolP("peers", "p", false, "sync with registered peers")
	syncCmd.MarkFlagsMutuallyExclusive("address", "scan", "peers")
}
