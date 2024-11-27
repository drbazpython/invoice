//Package cmd ...
package cmd

import (
	"fmt"
	"os"
	"strings"
	"drbaz.com/invoice/cmd/docprocessing"
	"github.com/glebarez/sqlite"
	"strconv"
	//"runtime"
	"drbaz.com/invoice/internal"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)
//PrintFlag A var to save the cmd print 
var PrintFlag bool

//! bubbletea
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type createModel struct {
	table table.Model
}

//InvoiceRow ...
var InvoiceRow table.Row

func (m createModel) Init() tea.Cmd { return nil }

func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		// case "enter":
		// 	return m, tea.Batch(
		//  		tea.Printf("Selected Invoice %s", m.table.SelectedRow()),
		//  	)
		}
		
	}
		InvoiceRow =  m.table.SelectedRow()
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	
}

func (m createModel) View() string {
	//TODO add a flag to create to print = true or false
	body := strings.Builder{}
	body.WriteString("\n\n")
	body.WriteString("\nA table of all invoices not yet completed \n")
	body.WriteString("Press enter to select a row, q or ctrl+c to continue and create pdf invoice\n")

	// _, fileName, lineNum, _ := runtime.Caller(0)
  	// fmt.Printf("%s: %d\n", fileName, lineNum)
	// body.WriteString(fileName)
	return baseStyle.Render("\n" + body.String() + "\n" + m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

//! end of bubbletea

// createCmd represents the add command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an invoice",
	Long:  `Get database data to create an doc 1stinvoice`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := internal.DefineLogger(log.InfoLevel)
		logger.Info("Creating invoice", "function", "create")
		PrintFlag, _ := cmd.Flags().GetBool("print")
		logger.Debug("PrintFlag", "value", PrintFlag)
		myStyle := internal.MyLipglossStyle("Create Invoice")
		fmt.Println(myStyle.Render("Create Invoice"))
		
		db, _ := gorm.Open(sqlite.Open("invoice.db"))
		invoices := []Invoice{}
		result := db.Find(&invoices)
		numOfInvoices := result.RowsAffected
		logger.Debug("Number of invoices", "in database",numOfInvoices)
		db.Model(Invoice{}).Find(&invoices)
		completed := "no"
		//var options []string
		var rows = []table.Row{}

		for _, invoice := range invoices {
			ID := strconv.FormatUint(uint64(invoice.ID), 10)
			if invoice.Completed { 
			completed = "yes"
			}
			rows = append(rows, []string{ID,invoice.InvoiceDate,invoice.Name,invoice.InvoiceNumber,invoice.StartDate,invoice.EndDate,invoice.HoursWorked,invoice.HourlyRate,completed})
		}
	

		columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Date", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "Inv No.", Width: 10},
		{Title: "Start",Width: 10},
		{Title: "End", Width: 10},
		{Title: "Hours", Width: 5},
		{Title: "Rate", Width: 5},
		{Title: "Completed", Width: 10},
		{Title: "PDF", Width: 20},
		}

		t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5), //TODO change to num of rows in database
	)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			BorderTop(true).
			Bold(true)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		m := createModel{t}
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		invoiceID := InvoiceRow[0]

		logger.Debug("Invoice","Row", InvoiceRow)

		invoiceDate := InvoiceRow[1]
		invoiceNum := InvoiceRow[3]
		invoiceHours := InvoiceRow[6]
		invoiceRate := InvoiceRow[7]
		invoiceStart := InvoiceRow[4]
		invoiceEnd := InvoiceRow[5]

		if PrintFlag {
			logger.Info("Invoice Details",
				"ID", invoiceID,
				"Date", invoiceDate,
				"Number", invoiceNum,
				"Hours", invoiceHours,
				"Rate", invoiceRate,
				"Start", invoiceStart,
				"End", invoiceEnd)
		}

		//(invoiceDate string,invoiceNumber string, invoiceHours string, invoiceStart string, invoiceEnd string, hourlyRate string
		newInvoice := docprocessing.ReplaceDocument(invoiceDate,invoiceNum,invoiceHours,invoiceStart,invoiceEnd,invoiceRate)
		if newInvoice != "" {
		 	logger.Debug("Replaced","new invoice",newInvoice)
		 } else {
		 	logger.Info("Not Replaced","function","create")
		}

		//Create PDF
		// newInvoice is the full path 
		pdfInvoice := docprocessing.CreatePDF(newInvoice,PrintFlag)
		//logger.Debug("pdf invoice","Alias",pdfInvoice)

		//* Mark invoice as completed
		logger.Debug("Marking Invoice completed","ID",invoiceID)
		pdf1 := strings.Split(pdfInvoice, ":")
		//logger.Debug(pdf1)
		pdf := pdf1[len(pdf1)-1]
		logger.Info("PDF","File Name",pdf)
		// Update with completed and pdf
		db.Model(&Invoice{}).Where("ID = ?", invoiceID).Updates(Invoice{Completed: true, Pdf: pdf})
		// db.Model(&Invoice{}).Where("ID = ?", invoiceID).Update("Completed",true)

		
	},
	
	
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().Bool("print", false, "Print invoice details")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
