# Invoice

#### This project aims to get a user time sheet input via a CLI, 
#### replace key fields in a word template,
#### save the new word doc as a pdf,
#### and update a database

Usage

- invoice add
  - user prompted for invoice details. After confirmation, invoice saved to database
- invoice list
  - list of invoices from database presented as a table
- invoice create
  - a list of invoices is presented as a table. User selects an invoice by moving
  up/down in table. The required invoice is selected by pressing return and "q".
  - the invoice template is transformed by replacing the [ ] with the data from the database.
  - a pdf is created and saved in the Documents folder.
  The invoice is not printed as the default --print flag is false
- invoice create --print=true
  - as above, but the invoice is printed to the default printer

Installation

- create a folder 'invoice' <project folder> and cd into it
- Run 'go mod init drbaz.com/invoice'
- Run 'gofolders'
- Enter project name as invoice
- Run 'go mod tidy'
- go build -o ./bin drbaz.com/invoice/cmd/invoice
- go install drbaz.com/invoice/cmd/invoice
- Then you can run invoice
- Also can run makefile buildInvoice

**Specification**

- Required modules:

  - github.com/charmbracelet/bubbletea
  - github.com/nguyenthenguyen/docx
  - strconv
  - github.com/charmbracelet/log
  - github.com/glebarez/sqlite
  - gorm.io/gorm
  - github.com/pterm/pterm
  - os
  - time

**Code Modules**

main.go - main code workflow

cobraui.go - get user input via cobra, pterm and bubbletea

documentProcessing.go - replace text in doc

database.go - save to database

```
go install github.com/spf13/cobra-cli@latest

 cobra-cli init
```

### TODO's

- [] TODO edit: Edit a row in the database
- [] TODO delete: delete a row in the database
- [] TODO add loglevel from .env using gotdotenv

