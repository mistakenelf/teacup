.PHONY: example-filetree
example-filetree:
	@go run ./examples/filetree/filetree.go

.PHONY: example-help
example-help:
	@go run ./examples/help/help.go

.PHONY: example-code
example-code:
	@go run ./examples/code/code.go
