//Generated Code

#include <string>
#include <vector>
#include <json/json.h>
#include "modelmakerdefs.hpp"


using std::string;
using std::vector;

namespace models {

class Model1: public Model {

	public:

		Model1();

		void load(string text);

		string write();

		bool load(json_object *obj);

		json_object* buildJsonObj();

		string field1;
		int field2;
		vector<int > field3;
		vector<vector<string > > field4;
};

};

