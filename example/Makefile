cpp:
	../cyborgbear -i model.cy -o Model
	c++ -std=c++11 main.cpp Model.cpp -o main -ljansson
go:
	../cyborgbear -i model.cy -o models.go -t go
	go build models.go
