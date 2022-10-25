/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"os"

	"github.com/Lanly109/lan/utils"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	Version string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lan",
	Short: "A collection of common tools for code collection",
	Long: `A collection of tools that is suitable for various requirements in the code collection of OI competition, including 
    - cleaning up irrelevant files
    - checking lists 
    - checking file time 
    - generating md5 codes
    - converting file formats
    ......`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose = viper.GetBool("Debug")

		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug Mode")
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	log.SetFormatter(&nested.Formatter{
		TimestampFormat: utils.TimeTemplate,
		FieldsOrder:     []string{"StartTime", "EndTime"},
	})

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "more info for debug")
	viper.BindPFlag("Debug", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.Version = Version
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find working directory.
		home, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Warn(err)
		log.Warn("Use commandline values.")
	}
}
