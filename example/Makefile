cpp:
	../cyborgwombat -i model.cy -o Model -n "wombat::models"
	c++ -std=c++11 main.cpp Model.cpp -o main -ljansson
go:
	../cyborgwombat -i model.cy -o models.go -t go
	go build models.go
