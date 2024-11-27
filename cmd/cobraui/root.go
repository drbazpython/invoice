// Package cmd ...
package cmd

import (
	"fmt"
	"os"

	//"drbaz.com/invoice/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pterminvoice", // name of command
	Short: "A CLI App to create an invoice",
	Long: `A CLI App to create an invoice by entering data via the terminal and replacing text in a word template. For example:
	invoice add
 	invoice create
 	invoice list`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		toggle, _ := cmd.Flags().GetBool("toggle")
		if toggle {
			fmt.Println("toggle is true")
		}
		//fmt.Println(internal.DefineMyLogo())
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pterminvoice.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
