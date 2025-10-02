.PHONY: run
run:
	go run cmd/demo-cli/*.go
	open tax50tawi-stamped.pdf

.PHONY: run-with-args
run-with-args:
	go run cmd/demo-cli/*.go -output tax50tawi-stamped.pdf -signature cmd/demo-cli/signature-1024x278.png -logo cmd/demo-cli/logo-410x361.png
	open tax50tawi-stamped.pdf

.PHONY: clean
clean:
	rm -rf tax50tawi-stamped.pdf
