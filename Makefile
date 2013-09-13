build:
	go build -o cyborgbear
install: build
	cp cyborgbear $(GOPATH)/bin
fmt:
	go fmt
