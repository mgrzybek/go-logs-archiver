BINARY      = go-logs-archiver
MODULE      = github.com/mgrzybek/go-logs-archiver
GO_VERSION  = 1.19
LAST_COMMIT = $(shell git rev-parse HEAD)

##############################################################################
# Go: usual targets

${BINARY}: go.mod ## Build the program (default target)
	go build -o ${BINARY} main.go

.PHONY: all
all: ${BINARY} ## Create the program

.PHONY: clean
clean: ## Clean the created artifacts
	rm -f ${BINARY} c.out

.PHONY: test
test: ## Run go test
	go test ./...

##############################################################################
# Go: packaging targets

.PHONY: get
get: ## Download required modules
	go get ./...

go.mod:
	go mod init ${MODULE}
	go mod tidy

##############################################################################
# Go: quality targets

c.out: ## Create the coverage file
	go test ./... -coverprofile=c.out

.PHONY: coverage
coverage: c.out ## Show the coverage ratio per function
	go tool cover -func=c.out

.PHONY: coverage-code
coverage-code: c.out ## Show the covered code in a browser
	go tool cover -html=c.out

##############################################################################
# Help

.PHONY: help
help: ## This help message
	@awk -F: \
		'/^([a-z\.-]+): *.* ## (.+)$$/ {gsub(/: .*?\s*##/, "\t");print}' \
		Makefile \
	| expand -t20 \
	| sort
