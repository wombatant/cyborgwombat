build:
	go build -o modelmaker
clean:
	rm -f modelmaker
fmt:
	go fmt
run: build
	./modelmaker image.txt
