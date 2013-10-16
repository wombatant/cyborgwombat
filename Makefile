build:
	go build -o cyborgbear
install: build
	mkdir -p $(GOPATH)/bin
	cp cyborgbear $(GOPATH)/bin
fmt:
	go fmt
