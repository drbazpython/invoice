// Package cmd ...
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/pterm/pterm"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
)

// invoiceCmd represents the invoice command
var invoiceCmd = &cobra.Command{
	Use:   "invoice",
	Short: "Create an invoice from a row in the database",
	Long: `this command lists the rows in the database, the user selects the 
	id of the appropraite row and an invoice is created by replacing text in a template`,
	Run: func(cmd *cobra.Command, args []string) {
		
		logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
		logger.Info("creating invoice",logger.Args("supper","cool"))
		
		newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgBlack),
		BackgroundStyle: pterm.NewStyle(pterm.BgRed),
		Margin:          20,
	}
		newHeader.Println("Create Invoice - deprecated")
		//run list command
		listCmd.Run(cmd, []string{})
		//TODO implement a tablewith selectable items?
		//user enters the id
		id := pterm.DefaultInteractiveTextInput
		id.DefaultText ="Please enter the invoice id"
		id1, _ := id.Show()
		result, _ := pterm.DefaultInteractiveConfirm.Show()
		// Print a blank line for better readability.
		pterm.Println()
		// Print the user's answer in a formatted way.
		logger.Info("you answered: ",logger.Args("answer",BoolToText(result)))
		
		logger.Info("Database ",logger.Args("id", id1))

		// confirm that this is correct invoice data
		var invoice Invoice
		db, _ := gorm.Open(sqlite.Open("test.db"))

		db.First(&invoice,"2")
		logger.Info(invoice.Name,logger.Args("priority","high"))
	},
}

func init() {
	rootCmd.AddCommand(invoiceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// invoiceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// invoiceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func boolToText1(b bool) string {
	if b {
		return pterm.Green("Yes")
	}
	return pterm.Red("No")
}
