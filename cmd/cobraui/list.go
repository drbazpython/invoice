// Package cmd ...
package cmd

import (
	"strconv"
	"fmt"
	"os"
	"strings"
	"drbaz.com/invoice/internal"
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	//"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	
)

type listModel struct {
	table table.Model
}

func (m listModel) View() string {
	body := strings.Builder{}
	body.WriteString("\n")
	body.WriteString("\nA table of all invoices \n")
	body.WriteString("Press enter, q, ctrl+c  or esc to quit\n")
	return baseStyle.Render("\n" + body.String() + "\n" + m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func (m listModel) Init() tea.Cmd { return nil }

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "esc":
		// 	if m.table.Focused() {
		// 		m.table.Blur()
		// 	} else {
		// 		m.table.Focus()
		// 	}
		case "q", "ctrl+c","esc","enter":
			return m, tea.Quit
		// case "enter":
		// 	return m, tea.Batch(
		//  		tea.Printf("Selected Invoice %s", m.table.SelectedRow()),
		//  	)
		}
		
	}
		// InvoiceRow =  m.table.SelectedRow()
		// m.table, cmd = m.table.Update(msg)
		return m, cmd
	
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List invoices",
	Long:  `List all invoice in database in a table`,
	Run: func(cmd *cobra.Command, args []string) {
		
		logger := internal.DefineLogger(log.DebugLevel)

		myStyle := internal.MyLipglossStyle("List Invoices")
		fmt.Println(myStyle.Render("List Invoices"))

		logger.Debug("list database called", "function", "list")
		/*tableData := pterm.TableData{
			{"id", "Date","Name", "Inv Num","Start","End","Hours","Rate","Completed","PDF"},
		}*/
		db, _ := gorm.Open(sqlite.Open("invoice.db"),&gorm.Config{
			//Logger: log.Default.LogMode(logger),
		})
		invoices := []Invoice{}
		
		result := db.Model(Invoice{}).Find(&invoices)
		if result.Error != nil{
			logger.Error(result.Error)
		}
		//invoices[5].Completed = true
		completed := "no"
		var rows = []table.Row{}
		for _, invoice := range invoices {
			ID := strconv.FormatUint(uint64(invoice.ID), 10)
			if invoice.Completed {
				completed = "yes"
			}
			//options = append(options, ID+" "+invoice.InvoiceNumber)
			rows = append(rows, []string{ID, invoice.InvoiceDate,invoice.Name, invoice.InvoiceNumber, invoice.StartDate,invoice.EndDate,invoice.HoursWorked,invoice.HourlyRate, completed,invoice.Pdf})
		}

		columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Date", Width: 20},
		{Title: "Name", Width: 10},
		{Title: "Inv No.", Width: 10},
		{Title: "Start",Width: 20},
		{Title: "End", Width: 20},
		{Title: "Hours", Width: 5},
		{Title: "Rate", Width: 5},
		{Title: "Completed", Width: 10},
		{Title: "PDF", Width: 20},
		}

			t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(5), //TODO change to num of rows in database
	)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			BorderBottom(true).
			BorderTop(true).
			Bold(true)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("00")).
			Background(lipgloss.Color("#FFFFFF")).
			Bold(false)
		t.SetStyles(s)

		m := listModel{t}
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		//pterm.DefaultTable.WithBoxed().WithHasHeader().WithData(tableData).Render()

		// 4
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
