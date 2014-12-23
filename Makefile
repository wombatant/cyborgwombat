build:
	go build -o cyborgwombat
install: build
	mkdir -p $(GOPATH)/bin
	cp cyborgwombat $(GOPATH)/bin
fmt:
	make -C parser
	go fmt
