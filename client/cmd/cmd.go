package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
  "github.com/spf13/viper"
  
  "github.com/piratey7007/rediss/client/connections"
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
      Password: viper.GetString("password"),
      Db: viper.GetInt("db"),
    }
    connections.ConnectToServer(options)
	},
  
}

func init() {
  helpFlag := false
  rootCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringP("host", "h", "localhost", "The host to bind to")
  rootCmd.PersistentFlags().IntP("port", "p", 6379, "The port to bind to")
  rootCmd.PersistentFlags().StringP("password", "a", "", "The password to use when connecting to the server")
  rootCmd.PersistentFlags().IntP("db", "n", 0, "The database to use when connecting to the server")
  viper.BindPFlag("bind", rootCmd.PersistentFlags().Lookup("host"))
  viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
  viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
  viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
}

func initConfig() {
  viper.SetConfigName("config")
  viper.AddConfigPath(".")
  viper.SetConfigType("conf")

  viper.SetDefault("bind", "localhost")
  viper.SetDefault("port", 6379)
  viper.SetDefault("password", "")
  viper.SetDefault("db", 0)

  if err := viper.ReadInConfig(); err != nil {
    fmt.Println("Error reading config file:", err)
  }
}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}

