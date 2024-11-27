package docprocessing

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	//"time"
	//"fmt"
	//cmd "drbaz.com/invoice/cmd/cobraui"
	"drbaz.com/invoice/internal"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/nguyenthenguyen/docx"
	//"github.com/thecodingmachine/gotenberg-go-client/v7"
	//"github.com/andybrewer/mack"
)

// GetTemplate ...
func GetTemplate() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	wordTemplate := os.Getenv("WORD_TEMPLATE")
	log.Debug(wordTemplate)
	log.Debug("docprocessing Package")
	return wordTemplate
}

// ReplaceDocument ...
func ReplaceDocument(invoiceDate string, invoiceNumber string, invoiceHours string, invoiceStart string, invoiceEnd string, hourlyRate string) string {
	//TODO get template from user home documents folder
	logger := internal.DefineLogger(log.InfoLevel)
	logger.Info("Starting Invoice Creation")
	docDir, _ := os.UserHomeDir()
	docDir = docDir + "/Documents/"

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	template := os.Getenv("WORD_TEMPLATE")
	template = docDir + template
	newInvoice := os.Getenv("REPLACED_WORD_TEMPLATE")
	
	// Read from docx file
	r, err := docx.ReadDocxFile(template)
	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)

	// Or read from a filesystem object:
	// r, err := docx.ReadDocxFromFS(file string, fs fs.FS)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// calculate total invoice amount invoiceHours x hourlyRate
	iHours, err := strconv.Atoi(invoiceHours)
	if err != nil {
		logger.Fatal(err)
	}

	iRate, err := strconv.Atoi(hourlyRate)
	if err != nil {
		logger.Fatal(err)
	}
	totalHours := strconv.Itoa(iHours * iRate)

	//now := time.Now()

	//logger.Debug(now.Format("01 January 2006"))

	docx1 := r.Editable()
	// replace data in invoice [invoicedate] [number] [hours] [rate] [total] [start] [end]

	err = docx1.Replace("[invoicedate]", invoiceDate, -1)
	if err != nil {
		logger.Debug("Error", "replacing invoicedate", err.Error())
	}
	docx1.Replace("[number]", invoiceNumber, -1)
	docx1.Replace("[hours]", invoiceHours, -1)
	docx1.Replace("[rate]", hourlyRate, -1)
	docx1.Replace("[total]", totalHours, -1)
	docx1.Replace("[start]", invoiceStart, -1)
	docx1.Replace("[end]", invoiceEnd, -1)

	//invoice number to end of document name
	fileName := strings.TrimSuffix(newInvoice, filepath.Ext(newInvoice))
	if invoiceNumber == "" {
		invoiceNumber = "XXXXXX"
	}
	// HS_Invoice
	//TODO Replace XX with invoiceNumber
	newInvoice = docDir + fileName + invoiceNumber + ".docx"
	//logger.Debug(newInvoice)

	err = docx1.WriteToFile(newInvoice)
	if err != nil {
		panic(err)
	}
	r.Close()
	logger.Debug("Invoice created successfully", "file", newInvoice)
	return newInvoice
}

// CreatePDF converts a Word document to PDF using LibreOffice
func CreatePDF(newInvoice string,PrintFlag bool) string {
	logger := internal.DefineLogger(log.InfoLevel)
	logger.Debug("CreatePDF", "file", newInvoice)
	logger.Debug("Print PDF", "flag", PrintFlag)
	// Get the directory of the input file
	dir := filepath.Dir(newInvoice)
	
	// Construct the command to convert to PDF
	cmd := exec.Command("soffice",
		"--headless",
		"--convert-to", "pdf",
		"--outdir", dir,
		newInvoice)

	// Run the conversion
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Failed to convert to PDF", "error", err, "output", string(output))
		return ""
	}

	// Get the PDF filename (same as input file but with .pdf extension)
	pdfFile := strings.TrimSuffix(newInvoice, filepath.Ext(newInvoice)) + ".pdf"

	logger.Info("PDF created successfully", "pdf", pdfFile)

	// List available printers
	printers, err := ListPrinters()
	if err == nil {
		for _, printer := range printers {
			logger.Info("Available printer:", "name", printer)
		}
	}

	// Automatically print one copy to default printer
	err = PrintPDF(pdfFile, "", 1,PrintFlag)
	if err != nil {
		logger.Error("Failed to print PDF", "error", err)
	}

	return pdfFile
}

// PrintPDF prints a PDF file using the system's default printer or a specified printer
func PrintPDF(pdfPath string, printerName string, copies int,print bool) error {
	logger := internal.DefineLogger(log.InfoLevel)

	// Verify the PDF file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF file not found: %s", pdfPath)
	}

	// Build the print command
	args := []string{pdfPath}

	// Add printer name if specified
	if printerName != "" {
		args = append([]string{"-d", printerName}, args...)
	}

	// Add number of copies if more than 1
	if copies > 1 {
		args = append([]string{"-n", strconv.Itoa(copies)}, args...)
	}
	//print = false
	logger.Debug("Pring PDF","flag",print)
	if print{
		// Execute the lp command
		cmd := exec.Command("lp", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.Error("Failed to print PDF", "error", err, "output", string(output))
			return fmt.Errorf("failed to print PDF: %v - %s", err, string(output))
		}
		logger.Info("PDF printed successfully", "file", pdfPath, "printer", printerName, "copies", copies)
		return nil
	}
		logger.Info("Do not print PDF, --print=false", "file", pdfPath, "printer", printerName, "print", print)
		return nil
	}	

// ListPrinters returns a list of available printers on the system
func ListPrinters() ([]string, error) {
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list printers: %v", err)
	}

	// Parse the output to get printer names
	lines := strings.Split(string(output), "\n")
	printers := make([]string, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "printer") {
			// Extract printer name (format: "printer PrinterName is idle")
			parts := strings.Fields(line)
			if len(parts) > 1 {
				printers = append(printers, parts[1])
			}
		}
	}

	return printers, nil
}

// TestPDFCreation tests the PDF creation and printing functionality
func TestPDFCreation() error {
	logger := internal.DefineLogger(log.DebugLevel)

	// Get user's Documents directory
	docDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}
	docDir = docDir + "/Documents/"

	// Load environment variables
	err = godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Create a test invoice
	logger.Info("Creating test invoice...")
	testInvoice := ReplaceDocument(
		"2024-03-20", // Invoice date
		"TEST001",    // Invoice number
		"10",         // Hours
		"2024-03-01", // Start date
		"2024-03-15", // End date
		"20",         // Hourly rate
	)

	if testInvoice == "" {
		return fmt.Errorf("failed to create test invoice")
	}
	logger.Info("Test invoice created", "path", testInvoice)

	// Convert to PDF and print
	logger.Info("Converting to PDF...")
	pdfPath := CreatePDF(testInvoice,false)
	if pdfPath == "" {
		return fmt.Errorf("failed to create PDF")
	}
	logger.Info("PDF created successfully", "path", pdfPath)

	return nil
}

//PrintDocument ...
// func PrintDocument(doc string) []byte {
// 	s := "ls"
// 	cmd1 := exec.Command(s)
// 	error := cmd1.Run()
// 	if error !=nil {
// 		log.Printf(error.Error())
// 	}
// 	log.Info("%s\n",b)
// 	return "ok"
// }
