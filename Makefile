BINARY      = go-logs-archiver
MODULE      = github.com/mgrzybek/go-logs-archiver
GO_VERSION  = 1.19
LAST_COMMIT = $(shell git rev-parse HEAD)

.PHONY: help
help: ## This help message
	@awk -F: \
		'/^([a-z-]+): [a-z/- ]*## (.+)$$/ {gsub(/: .*?\s*##/, "\t");print}' \
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

${BINARY}: go.mod ## Test and build the program
	go build -o ${BINARY} main.go