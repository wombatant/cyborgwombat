//Generated Code

#include <string>
#include <sstream>
#include <vector>
#include <map>
#include <json/json.h>
#include "modelmakerdefs.hpp"


using std::string;
using std::vector;
using std::map;

namespace models {

class Model1: public Model {

	public:

		Model1();

		void load(string text);

		string write();

		bool load(json_object *obj);

		json_object* buildJsonObj();

		string field1;
		unknown field2;
		vector<int > field3;
		vector<vector<string > > field4;
		map<string, string > field5;
};

}

