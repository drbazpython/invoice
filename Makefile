 
 
 hello:
	echo "Hello to the Invoice Package"

run:
	go run main.go

buildInvoice:
	@echo "Building .... "
	@go build -o ./bin drbaz.com/invoice/cmd/invoice
	@echo "Installing ...."
	@go install drbaz.com/invoice/cmd/invoice
	@echo "invoice built to bin and installed"
	@echo "Running invoice ...."
	@invoice

runalltests:
	go test -v ./...

rundbtests:
	go test -v ./cmd/database

