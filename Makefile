.PHONY: example-filetree
example-filetree:
	@go run ./examples/filetree/filetree.go

.PHONY: example-help
example-help:
	@go run ./examples/help/help.go

.PHONY: example-sourcecode
example-sourcecode:
	@go run ./examples/sourcecode/sourcecode.go
