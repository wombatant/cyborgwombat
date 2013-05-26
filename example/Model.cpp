//Generated Code

#include "Model.hpp"

using namespace models;

Model1::Model1() {
	this->field1 = "";
	this->field2 = 0;
}

void Model1::load(string json) {
	json_object *obj = json_tokener_parse(json.c_str());
	load(obj);
	json_object_put(obj);
}

bool Model1::load(json_object *in) {
	{
		json_object *obj0 = json_object_object_get(in, "Field1");
		if (obj0 != NULL) {
			string out0;
			if (json_object_get_type(obj0) == json_type_string) {
				out0 = json_object_get_string(obj0);
			}
			this->field1 = out0;
		}
	}
	{
		json_object *obj0 = json_object_object_get(in, "Field2");
		if (obj0 != NULL) {
			int out0;
			if (json_object_get_type(obj0) == json_type_int) {
				out0 = json_object_get_int(obj0);
			}
			this->field2 = out0;
		}
	}
	{
		json_object *obj1 = json_object_object_get(in, "Field3");
		if (obj1 != NULL) {
			if (json_object_get_type(obj1) != json_type_array) {
				return false;
			}
			int size = json_object_array_length(obj1);
			vector<int > out1;
			for (int i = 0; i < size; i++) {
				json_object *obj0 = json_object_array_get_idx(obj1, i);
				int out0;
				if (json_object_get_type(obj0) == json_type_int) {
					out0 = json_object_get_int(obj0);
				}
				this->field3.push_back(out0);
			}
		}
	}
	{
		json_object *obj2 = json_object_object_get(in, "Field4");
		if (obj2 != NULL) {
			if (json_object_get_type(obj2) != json_type_array) {
				return false;
			}
			int size = json_object_array_length(obj2);
			vector<vector<string > > out2;
			for (int i = 0; i < size; i++) {
				json_object *obj1 = json_object_array_get_idx(obj2, i);
				vector<string > out1;
				{
					if (obj1 != NULL) {
						if (json_object_get_type(obj1) != json_type_array) {
							return false;
						}
						int size = json_object_array_length(obj1);
						for (int i = 0; i < size; i++) {
							json_object *obj0 = json_object_array_get_idx(obj1, i);
							string out0;
							if (json_object_get_type(obj0) == json_type_string) {
								out0 = json_object_get_string(obj0);
							}
							out1.push_back(out0);
						}
					}
				}
				this->field4.push_back(out1);
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
		json_object *out0 = json_object_new_int(this->field2);
		json_object_object_add(obj, "Field2", out0);
	}
	{
		json_object *array0 = json_object_new_array();
		for (int i = 0; i < this->field3.size(); i++) {
			json_object *out0 = json_object_new_int(this->field3[i]);
			json_object_array_add(array0, out0);
		}
		json_object_object_add(obj, "Field3", array0);
	}
	{
		json_object *array0 = json_object_new_array();
		for (int i = 0; i < this->field4.size(); i++) {
			json_object *array1 = json_object_new_array();
			for (int ii = 0; ii < this->field4[i].size(); ii++) {
				json_object *out0 = json_object_new_string(this->field4[i][ii].c_str());
				json_object_array_add(array1, out0);
			}
			json_object_array_add(array0, array1);
		}
		json_object_object_add(obj, "Field4", array0);
	}
	return obj;
}
