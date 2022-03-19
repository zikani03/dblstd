package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sethbonnie/dblstd/shape"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dblstd --shapeFile <path> <repo_path>",
	Short: "checks if a repo conforms to a given standard",
	Long: `dblstd - short for DoubleStandards - checks if a project repo
conforms to a given standard (in the form of a "shape" file).`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		shapeFile, err := cmd.Flags().GetString("shape-file")
		if err != nil {
			return err
		}
		shapeData, err := ioutil.ReadFile(shapeFile)
		if err != nil {
			return err
		}

		missing, err := shape.Missing(args[0], shapeData, 0)
		if err != nil {
			return err
		}
		for path, isDir := range missing {
			if isDir {
				fmt.Fprintf(os.Stderr, "Missing required directory: %s\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "Missing required file: %s\n", path)
			}
		}
		return nil
	},
	Version: "0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dblstd.yaml)")

	rootCmd.PersistentFlags().StringP("shape-file", "s", "", "file containing expected shape of repo")

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dblstd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dblstd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}