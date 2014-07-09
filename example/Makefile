cpp:
	../cyborgbear -i model.cy -o Model -cpp-boost
	c++ -std=c++11 main.cpp Model.cpp -o main -ljansson -lboost_serialization
go:
	../cyborgbear -i model.cy -o models.go -t go
	go build models.go
