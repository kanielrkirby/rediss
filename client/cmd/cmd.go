package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piratey7007/rediss/client/connections"
)

var rootCmd = &cobra.Command{
	Use:   "rediss-cli",
	Short: "Redis CLI interacts with a Redis server",
	Long: `A custom, simplified CLI to interact with Rediss server that takes user commands,
converts them to the Redis Serialization Protocol (RESP), and forwards them to the Rediss server.`,
	Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Connecting to server...")
    options := connections.ConnectionOptions{
      Host: cmd.Flag("host").Value.String(),
      Port: cmd.Flag("port").Value.String(),
    }
    connections.ConnectToServer(options)
	},
  
}

func init() {
  helpFlag := false
  rootCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")
  rootCmd.PersistentFlags().StringP("host", "h", "localhost", "The host to bind to")
  rootCmd.PersistentFlags().StringP("port", "p", "6379", "The port to bind to")
}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}

