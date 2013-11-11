build:
	go build -o cyborgbear
install: build
	mkdir -p $(GOPATH)/bin
	cp cyborgbear $(GOPATH)/bin
fmt:
	make -C parser
	go fmt
