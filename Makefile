help:
	@echo "Use one of the following sub-commands with make:"
	@echo "  - lint         To run the linter against the go services."

#
# Make commands.
#

lint:
	golangci-lint run ./...
