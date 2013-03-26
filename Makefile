build:
	go build
fmt:
	go fmt
run: build
	./modelmaker image.txt
