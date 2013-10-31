build: model
	c++ main.cpp Model.cpp -o main -ljansson
model: model.cy
	../cyborgbear -i model.cy -o Model
clean:
	rm -f main Model.cpp Model.hpp
