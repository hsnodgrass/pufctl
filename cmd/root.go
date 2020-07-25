package cmd

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hsnodgrass/pufctl/config"
	"github.com/hsnodgrass/pufctl/internal/util"
)

var (
	cfgFile string
	path    string

	rootCmd = &cobra.Command{
		Use:   "pufctl",
		Short: "Puppetfile Swiss Army Knife",
		Long:  `A tool for doing everything you need to with Puppetfiles.`,
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.ExitErr(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "Path to your Puppetfile")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to config file config file (default is $HOME/.pufctl.yaml or \".\")")
	rootCmd.SetVersionTemplate(config.Version())
}

func initConfig() {
	viper.SetDefault("puppetfile", "./Puppetfile")
	viper.SetDefault("forge_url", "https://forgeapi.puppetlabs.com")
	viper.SetDefault("forge_proxy", "")
	viper.SetDefault("git_private_key", "")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".pufctl")
		viper.SetConfigType("yaml")
		home, err := homedir.Dir()
		if err != nil {
			viper.AddConfigPath(".")
		} else {
			viper.AddConfigPath(home)
		}
	}
	if path != "" {
		viper.Set("puppetfile", path)
	}
	err := util.ValidatePath(viper.GetString("puppetfile"))
	if err != nil {
		util.ExitErr(err)
	}
	viper.SetEnvPrefix("pufctl")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
