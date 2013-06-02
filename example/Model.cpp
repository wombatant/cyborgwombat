//Generated Code

#include "Model.hpp"

using namespace models;
using std::stringstream;

Model1::Model1() {
	this->field1 = "";
}

void Model1::load(string json) {
	json_object *obj = json_tokener_parse(json.c_str());
	load(obj);
	json_object_put(obj);
}

bool Model1::load(json_object *in) {
	{
		json_object *obj0 = json_object_object_get(in, "Field1");
		if (json_object_get_type(obj0) == json_type_string) {
			this->field1 = json_object_get_string(obj0);
		}
	}
	{
		json_object *obj0 = json_object_object_get(in, "Field2");
		this->field2.load(obj0);
	}
	{
		json_object *obj0 = json_object_object_get(in, "Field3");
		if (obj0 != NULL && json_object_get_type(obj0) == json_type_array) {
			unsigned int size = json_object_array_length(obj0);
			this->field3.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				json_object *obj1 = json_object_array_get_idx(obj0, i);
				{
					json_object *obj0 = json_object_object_get(in, "Field3");
					if (obj0 != NULL) {
						if (json_object_get_type(obj1) == json_type_int) {
							this->field3[i] = json_object_get_int(obj1);
						}
					}
				}
			}
		}
	}
	{
		json_object *obj0 = json_object_object_get(in, "Field4");
		if (obj0 != NULL && json_object_get_type(obj0) == json_type_array) {
			unsigned int size = json_object_array_length(obj0);
			this->field4.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				json_object *obj1 = json_object_array_get_idx(obj0, i);
				if (obj1 != NULL && json_object_get_type(obj1) == json_type_array) {
					unsigned int size = json_object_array_length(obj1);
					this->field4[i].resize(size);
					for (unsigned int ii = 0; ii < size; ii++) {
						json_object *obj2 = json_object_array_get_idx(obj1, ii);
						if (json_object_get_type(obj2) == json_type_string) {
							this->field4[i][ii] = json_object_get_string(obj2);
						}
					}
				}
			}
		}
	}
	{
		json_object *obj0 = json_object_object_get(in, "Field5");
		if (obj0 != NULL && json_object_get_type(obj0) == json_type_object) {
			json_object_object_foreach(obj0, key, obj1) {
				string i;
				{
					std::stringstream s;
					s << key;
					s >> i;
				}
				{
					json_object *obj0 = json_object_object_get(in, "Field5");
					if (obj0 != NULL) {
						if (json_object_get_type(obj1) == json_type_string) {
							this->field5[i] = json_object_get_string(obj1);
						}
					}
				}
			}
		}
	}
	return true;
}

string Model1::write() {
	json_object *obj = buildJsonObj();
	string out = json_object_to_json_string(obj);
	json_object_put(obj);
	return out;
}

json_object* Model1::buildJsonObj() {
	json_object *obj = json_object_new_object();
	{
		json_object *out0 = json_object_new_string(this->field1.c_str());
		json_object_object_add(obj, "Field1", out0);
	}
	{
		json_object *out0 = this->field2.buildJsonObj();
		json_object_object_add(obj, "Field2", out0);
	}
	{
		json_object *out1 = json_object_new_array();
		for (unsigned int i = 0; i < this->field3.size(); i++) {
			json_object *out0 = json_object_new_int(this->field3[i]);
			json_object_array_add(out1, out0);
		}
		json_object_object_add(obj, "Field3", out1);
	}
	{
		json_object *out2 = json_object_new_array();
		for (unsigned int i = 0; i < this->field4.size(); i++) {
			json_object *out1 = json_object_new_array();
			for (unsigned int ii = 0; ii < this->field4[i].size(); ii++) {
				json_object *out0 = json_object_new_string(this->field4[i][ii].c_str());
				json_object_array_add(out1, out0);
			}
			json_object_array_add(out2, out1);
		}
		json_object_object_add(obj, "Field4", out2);
	}
	{
		json_object *out1 = json_object_new_object();
		for (map<string, string >::iterator n = this->field5.begin(); n != this->field5.end(); n++) {
			std::stringstream s;
			string key;
			s << n->first;
			s >> key;
			json_object *out0 = json_object_new_string(this->field5[n->first].c_str());
			json_object_object_add(out1, key.c_str(), out0);
		}
		json_object_object_add(obj, "Field5", out1);
	}
	return obj;
}
