.PHONY: run
run:
	go run cmd/cli/*.go
	open tax50tawi-stamped.pdf

.PHONY: run-with-args
run-with-args:
	go run cmd/demo-cli/*.go -output tax50tawi-stamped.pdf -signature cmd/demo-cli/demo-signature-1280x720-rectangle.png -seal cmd/demo-cli/demo-logo-1024x1024-square.png
	open tax50tawi-stamped.pdf

.PHONY: clean
clean:
	rm -rf tax50tawi-stamped.pdf
