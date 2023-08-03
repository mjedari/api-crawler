package cmd

import (
	"github.com/mjedari/vgang-project/app/configs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var (
	configFile string

	rootCmd = &cobra.Command{
		Use:   "vgang-project",
		Short: "short description",
		Long:  `long description`,
	}
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "Mahdi Jedari", "i.jedari@gmail.com")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("vgang")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	viper.Unmarshal(&configs.Config)
	logrus.Info("configuration initialized! (Notice: configurations may be initialised from OS ENV)")
}
