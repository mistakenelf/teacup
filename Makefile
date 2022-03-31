.PHONY: example-filetree
example-filetree:
	@go run ./examples/filetree/filetree.go

.PHONY: example-help
example-help:
	@go run ./examples/help/help.go

.PHONY: example-code
example-code:
	@go run ./examples/code/code.go

.PHONY: example-markdown
example-markdown:
	@go run ./examples/markdown/markdown.go

.PHONY: example-pdf
example-pdf:
	@go run ./examples/pdf/pdf.go

.PHONY: example-statusbar
example-statusbar:
	@go run ./examples/statusbar/statusbar.go

.PHONY: example-iamge
example-image:
	@go run ./examples/image/image.go