//Generated Code
#ifndef MODEL_HPP
#define MODEL_HPP
#include <string>
#include <sstream>
#include <vector>
#include <map>
#include "modelmakerdefs.hpp"


using std::vector;
using std::map;

namespace models {

using modelmaker::string;

class Model1: public modelmaker::Model {

	public:

		Model1();

		bool loadJsonObj(modelmaker::JsonVal obj);

		modelmaker::JsonValOut buildJsonObj();

		string field1;
		modelmaker::unknown field2;
		vector< int > field3;
		vector< vector< string > > field4;
		map< string, string > field5;
};

}


#endif