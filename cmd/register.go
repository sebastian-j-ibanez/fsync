/*
Copyright © 2024 Sebastian Ibanez <sebas.ibanez219@gmail.com>
*/
package cmd

import (
	"fmt"
	"net"
	"os"
	"strconv"

	prot "github.com/sebastian-j-ibanez/fsync/protocol"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a peer",
	Long: `Start process to register a peer client.
	Program will ask for IP and then port.`,
	Run: func(cmd *cobra.Command, args []string) {
		var ip string
		var port uint16
		fmt.Println("please enter the peer IP address: ")
		fmt.Scanln(&ip)
		if net.ParseIP(ip) == nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: invalid ip\n")
			os.Exit(-1)
		}

		fmt.Println("please enter the peer port number: ")
		fmt.Scan(&port)
		peer := prot.Peer{
			IP:   ip,
			Port: strconv.FormatUint(uint64(port), 10),
		}
		err := prot.RegisterPeer(peer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fsync: ERROR: unable to register peer\n")
			os.Exit(-1)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
