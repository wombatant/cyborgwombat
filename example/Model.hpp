//Generated Code
#ifndef MODEL_HPP
#define MODEL_HPP
#include <string>
#include <sstream>
#include <vector>
#include <map>
#include <jansson.h>
#include "modelmakerdefs.hpp"


using std::string;
using std::vector;
using std::map;

namespace models {

class Model1: public modelmaker::Model {

	public:

		Model1();

		bool load_json_t(json_t *obj);

		json_t* buildJsonObj();

		string field1;
		modelmaker::unknown field2;
		vector< int > field3;
		vector< vector< string > > field4;
		map<string, string > field5;
};

}


#endif