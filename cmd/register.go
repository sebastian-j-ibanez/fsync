/*
Copyright © 2024 Sebastian Ibanez <sebas.ibanez219@gmail.com>
*/
package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

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
		var port string

		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "error: expected <IP:PORT>\n")
			os.Exit(-1)
		}

		// Split the port from the IP (IP:PORT)
		argParts := strings.Split(args[0], ":")
		if len(argParts) < 2 || argParts[0] == "" || argParts[1] == "" {
			fmt.Fprintf(os.Stderr, "error: expected <IP:PORT>\n")
			os.Exit(-1)
		}
		ip = argParts[0]
		port = argParts[1]

		if net.ParseIP(ip) == nil {
			fmt.Fprintf(os.Stderr, "error: invalid ip\n")
			os.Exit(-1)
		}

		peer := prot.Peer{
			IP:   ip,
			Port: port,
		}
		err := prot.RegisterPeer(peer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to register peer\n")
			os.Exit(-1)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
