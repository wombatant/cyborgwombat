build:
	go build -o modelmaker
install: build
	cp modelmaker $(GOPATH)/bin
fmt:
	go fmt
