build:
	go build -o modelmaker
fmt:
	go fmt
run: build
	./modelmaker image.txt
