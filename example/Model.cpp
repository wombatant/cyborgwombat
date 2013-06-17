//Generated Code

#include "Model.hpp"

using namespace models;
using std::stringstream;

Model1::Model1() {
	this->field1 = "";
}

bool Model1::load_json_t(json_t *in) {
	{
		json_t *obj0 = json_object_get(in, "Field1");
		{
			if (json_is_string(obj0)) {
				this->field1 = json_string_value(obj0);
			}
		}
		json_decref(obj0);
	}
	{
		json_t *obj0 = json_object_get(in, "Field2");
		{
			this->field2.load_json_t(obj0);
		}
		json_decref(obj0);
	}
	{
		json_t *obj0 = json_object_get(in, "Field3");
		if (obj0 != NULL && json_typeof(obj0) == JSON_ARRAY) {
			unsigned int size = json_array_size(obj0);
			this->field3.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				json_t *obj1 = json_array_get(obj0, i);
				{
					if (json_is_integer(obj1)) {
						this->field3[i] = json_integer_value(obj1);
					}
				}
			}
		}
		json_decref(obj0);
	}
	{
		json_t *obj0 = json_object_get(in, "Field4");
		if (obj0 != NULL && json_typeof(obj0) == JSON_ARRAY) {
			unsigned int size = json_array_size(obj0);
			this->field4.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				json_t *obj1 = json_array_get(obj0, i);
				if (obj1 != NULL && json_typeof(obj1) == JSON_ARRAY) {
					unsigned int size = json_array_size(obj1);
					this->field4[i].resize(size);
					for (unsigned int ii = 0; ii < size; ii++) {
						json_t *obj2 = json_array_get(obj1, ii);
						{
							if (json_is_string(obj2)) {
								this->field4[i][ii] = json_string_value(obj2);
							}
						}
					}
				}
			}
		}
		json_decref(obj0);
	}
	{
		json_t *obj0 = json_object_get(in, "Field5");
		if (obj0 != NULL && json_typeof(obj0) == JSON_OBJECT) {
			const char *key;
			json_t *obj1;
			json_object_foreach(obj0, key, obj1) {
				string i;
				{
					std::stringstream s;
					s << key;
					s >> i;
				}
				string val;
				this->field5.insert(make_pair(i, val));
				{
					if (json_is_string(obj1)) {
						this->field5[i] = json_string_value(obj1);
					}
				}
			}
		}
		json_decref(obj0);
	}
	return true;
}

json_t* Model1::buildJsonObj() {
	json_t *obj = json_object();
	{
		json_t *out0 = json_string(this->field1.c_str());
		json_object_set(obj, "Field1", out0);
		json_decref(out0);
	}
	{
		json_t *out0 = this->field2.buildJsonObj();
		json_object_set(obj, "Field2", out0);
		json_decref(out0);
	}
	{
		json_t *out1 = json_array();
		for (unsigned int i = 0; i < this->field3.size(); i++) {
			json_t *out0 = json_integer(this->field3[i]);
			json_array_append(out1, out0);
			json_decref(out0);
		}
		json_object_set(obj, "Field3", out1);
		json_decref(out1);
	}
	{
		json_t *out2 = json_array();
		for (unsigned int i = 0; i < this->field4.size(); i++) {
			json_t *out1 = json_array();
			for (unsigned int ii = 0; ii < this->field4[i].size(); ii++) {
				json_t *out0 = json_string(this->field4[i][ii].c_str());
				json_array_append(out1, out0);
				json_decref(out0);
			}
			json_array_append(out2, out1);
			json_decref(out1);
		}
		json_object_set(obj, "Field4", out2);
		json_decref(out2);
	}
	{
		json_t *out1 = json_object();
		for (map<string, string >::iterator n = this->field5.begin(); n != this->field5.end(); n++) {
			std::stringstream s;
			string key;
			s << n->first;
			s >> key;
			json_t *out0 = json_string(this->field5[n->first].c_str());
			json_object_set(out1, key.c_str(), out0);
			json_decref(out0);
		}
		json_object_set(obj, "Field5", out1);
		json_decref(out1);
	}
	return obj;
}
