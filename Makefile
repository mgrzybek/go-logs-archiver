BINARY      = go-logs-archiver
MODULE      = github.com/mgrzybek/go-logs-archiver
GO_VERSION  = 1.19
LAST_COMMIT = $(shell git rev-parse HEAD)

.PHONY: help
help: ## This help message
	@awk -F: \
		'/^([a-z\.-]+): *.* ## (.+)$$/ {gsub(/: .*?\s*##/, "\t");print}' \
		Makefile \
	| expand -t20 \
	| sort

##############################################################################
# Go

.PHONY: test
test: ## Run go test
	go test ./...

.PHONY: get
get: ## Download required modules
	go get ./...

go.mod:
	go mod init ${MODULE}
	go mod tidy

${BINARY}: go.mod test ## Test and build the program
	go build -o ${BINARY} main.go

c.out: ## Create the coverage file
	go test ./... -coverprofile=c.out

.PHONY: coverage
coverage: c.out ## Show the coverage ratio per function
	go tool cover -func=c.out

.PHONY: all
all: ${BINARY} ## Create the program

.PHONY: clean
clean: ## Clean the created artifacts
	rm -f ${BINARY}
