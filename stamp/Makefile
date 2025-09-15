.PHONY: run
run:
	go run stamp.go
	open tax50tawi-stamped.pdf

.PHONY: run-with-args
run-with-args:
	go run stamp.go -input tax50tawi.pdf -output tax50tawi-stamped.pdf -signature signature-1024x278.png
	open tax50tawi-stamped.pdf

.PHONY: clean
clean:
	rm -rf tax50tawi-stamped.pdf
