package cmd

import (
	"github.com/spf13/cobra"
  "github.com/spf13/viper"
  
  "github.com/piratey7007/rediss/server/connections"
)

var rootCmd = &cobra.Command{
	Use:   "rediss-cli",
	Short: "Redis CLI interacts with a Redis server",
	Long: `A custom, simplified CLI to interact with Rediss server that takes user commands,
converts them to the Redis Serialization Protocol (RESP), and forwards them to the Rediss server.`,
	Run: func(cmd *cobra.Command, args []string) {
    options := connections.ConnectionOptions{
      Host: viper.GetString("bind"),
      Port: viper.GetInt("port"),
    }
    connections.StartServer(options)
	},
  
}

func init() {
  helpFlag := false
  rootCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")
  rootCmd.PersistentFlags().StringP("bind", "b", "localhost", "The host to bind to")
  rootCmd.PersistentFlags().IntP("port", "p", 6379, "The port to bind to")
}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}

