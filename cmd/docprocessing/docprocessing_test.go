package docprocessing_test

import (
	//"os"
	//"fmt"
	"testing"
	"drbaz.com/invoice/cmd/docprocessing"
	//"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	
)

func TestDocProcessingTemplateExists(t *testing.T){
	wordTemplate := docprocessing.GetTemplate()
	assert.FileExists(t,wordTemplate,"Template %v does not exist",wordTemplate)
} 

func TestReplaceDocument(t *testing.T){
	newFile := docprocessing.ReplaceDocument("1st January 2024","XXXXX","20","1st Augut 2024","31st August 2024","20")
	assert.FileExists(t, newFile,"new_invoice_XXXXX.docx")
}


// func TestPrintDocument(t *testing.T){
// 	b := docprocessing.PrintDocument("xxxx")
// 	fmt.Printf("%s\n",b)
// 	assert.NotNil(t,b)
// }
