// Package cmd ...
package cmd

/*
Copyright © 2024 Dr Barrie Smith <drbazgo@gmail.com>
*/
import (
	"fmt"
	//"drbaz.com/invoice/cmd/docprocessing"
	"drbaz.com/invoice/internal"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	//"drbaz.com/invoice/cmd/logging"
	//"database/sql"
	//_ "github.com/marcboeker/go-duckdb"
)

//Document ...
type Document interface {
	Replace(template string) bool
}

//TODO add remining fields

//Invoice ...
type Invoice struct {
	gorm.Model
	InvoiceDate string
	Name      string
	InvoiceNumber    string
	HoursWorked string
	StartDate string
	EndDate string
	HourlyRate string
	Completed bool
	Pdf string

}
//Replace ...
// func (i Invoice) Replace(template string) string {
// 	newDocument := docprocessing.ReplaceDocument("HS1234","23","start","end","20")
//     return newDocument
//}
// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an invoice",
	Long:  `Get user input to create an invoice and save to database`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := internal.DefineLogger(log.DebugLevel)
		logger.Debug("Invoice ","function","add")

		myStyle := internal.MyLipglossStyle("Add Invoice")
		fmt.Println(myStyle.Render("Add Invoice"))
		//log := logging.DefineLogger(log.DebugLevel)
		//log.Info("add invoice called")
		// newHeader := pterm.HeaderPrinter{
		// 	TextStyle:       pterm.NewStyle(pterm.FgBlack),
		// 	BackgroundStyle: pterm.NewStyle(pterm.BgRed),
		// 	Margin:          20,
		// }
		// newHeader.Println("Add Invoice")
		//pterm.DefaultHeader.Println("Add an invoice")
		//pterm.Println()

		//* Invoice Date
		invoiceDate := pterm.DefaultInteractiveTextInput
		invoiceDate.DefaultText = "Please enter the invoice date"
		result0, _ := invoiceDate.Show()
		pterm.Println()

		//* Name
		name := pterm.DefaultInteractiveTextInput
		name.DefaultText = "Please enter your name"
		result1, _ := name.Show()
		pterm.Println()

		//* Invoice Number
		invNum := pterm.DefaultInteractiveTextInput
		invNum.DefaultText = "Please enter the invoice number"
		result2, _ := invNum.Show()
		pterm.Println()

		//* Start Date
		startDate := pterm.DefaultInteractiveTextInput
		startDate.DefaultText = "Please enter the start date"
		result3, _ := startDate.Show()
		pterm.Println()

		//* End Date
		endDate := pterm.DefaultInteractiveTextInput
		endDate.DefaultText = "Please enter the end date"
		result4, _ := endDate.Show()
		pterm.Println()

		//* Hours Worked
		hoursWorked := pterm.DefaultInteractiveTextInput
		hoursWorked.DefaultText = "Please enter the hours worked this period"
		result5, _ := hoursWorked.Show()
		pterm.Println()

		//* Hourly Rate - get from .env
		hourlyRate := internal.GetHourlyRate()


		result, _ := pterm.DefaultInteractiveConfirm.Show()
		// Print a blank line for better readability.
		pterm.Println()
		// Print the user's answer in a formatted way.
		logger.Debug("you answered: ", "answer", result)
		pterm.Printfln("You answered: %s", BoolToText(result))
		if result {
			db, _ := gorm.Open(sqlite.Open("invoice.db"))
			db.AutoMigrate(&Invoice{})
			db.Create(&Invoice{
				InvoiceDate: result0,
				Name:      result1,
				InvoiceNumber:    result2,
				StartDate: result3,
				EndDate: result4,
				HoursWorked: result5,
				HourlyRate: hourlyRate,
				Completed: false,
				Pdf: "not ready",
			})
			logger.Debug("saved invoice to database", "function", "add")
		} else {
			logger.Debug("invoice aborted", "function", "add")
			pterm.Printfln("Invoice aborted")
			
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
	//BoolToText ...
	func BoolToText(b bool) string {
	if b {
		return pterm.Green("Yes")
	}
	return pterm.Red("No")
}
