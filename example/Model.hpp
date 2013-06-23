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

		string Field1;
		modelmaker::unknown Field2;
		vector< int > Field3;
		vector< vector< string > > Field4;
		map< string, string > Field5;
};

}


#endif