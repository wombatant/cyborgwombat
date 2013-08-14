//Generated Code

#include "Model.hpp"

using namespace models;
using std::stringstream;

Model1::Model1() {
	this->field1 = "";
}

bool Model1::loadJsonObj(modelmaker::JsonVal in) {
	modelmaker::JsonObjOut inObj = modelmaker::toObj(in);
	{
		modelmaker::JsonValOut obj0 = modelmaker::objRead(inObj, "Field1");
		{
			if (modelmaker::isString(obj0)) {
				this->field1 = modelmaker::toString(obj0);
			}
		}
	}
	{
		modelmaker::JsonValOut obj0 = modelmaker::objRead(inObj, "Field2");
		{
			this->field2.loadJsonObj(obj0);
		}
	}
	{
		modelmaker::JsonValOut obj0 = modelmaker::objRead(inObj, "Field3");
		if (!modelmaker::isNull(obj0) && modelmaker::isArray(obj0)) {
			modelmaker::JsonArrayOut array0 = modelmaker::toArray(obj0);
			unsigned int size = modelmaker::arraySize(array0);
			this->field3.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				modelmaker::JsonValOut obj1 = modelmaker::arrayRead(array0, i);
				{
					if (modelmaker::isInt(obj1)) {
						this->field3[i] = modelmaker::toInt(obj1);
					}
				}
			}
		}
	}
	{
		modelmaker::JsonValOut obj0 = modelmaker::objRead(inObj, "Field4");
		if (!modelmaker::isNull(obj0) && modelmaker::isArray(obj0)) {
			modelmaker::JsonArrayOut array0 = modelmaker::toArray(obj0);
			unsigned int size = modelmaker::arraySize(array0);
			this->field4.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				modelmaker::JsonValOut obj1 = modelmaker::arrayRead(array0, i);
				if (!modelmaker::isNull(obj1) && modelmaker::isArray(obj1)) {
					modelmaker::JsonArrayOut array1 = modelmaker::toArray(obj1);
					unsigned int size = modelmaker::arraySize(array1);
					this->field4[i].resize(size);
					for (unsigned int ii = 0; ii < size; ii++) {
						modelmaker::JsonValOut obj2 = modelmaker::arrayRead(array1, ii);
						{
							if (modelmaker::isString(obj2)) {
								this->field4[i][ii] = modelmaker::toString(obj2);
							}
						}
					}
				}
			}
		}
	}
	{
		modelmaker::JsonValOut obj0 = modelmaker::objRead(inObj, "Field5");
		if (!modelmaker::isNull(obj0) && modelmaker::isObj(obj0)) {
			modelmaker::JsonObjOut map0 = modelmaker::toObj(obj0);
			for (modelmaker::JsonObjIterator it1 = modelmaker::iterator(map0); !modelmaker::iteratorAtEnd(it1, map0); it1 = modelmaker::iteratorNext(map0,  it1)) {
				string i;
				modelmaker::JsonObjIteratorVal obj1 = modelmaker::iteratorValue(it1);
				{
					std::string key = modelmaker::toStdString(modelmaker::iteratorKey(it1));
					std::string o;
					std::stringstream s;
					s << key;
					s >> o;
					i = o.c_str();
				}
				{
					if (modelmaker::isString(obj1)) {
						this->field5[i] = modelmaker::toString(obj1);
					}
				}
			}
		}
	}
	return true;
}

modelmaker::JsonValOut Model1::buildJsonObj() {
	modelmaker::JsonObjOut obj = modelmaker::newJsonObj();
	{
		modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->field1);
		modelmaker::objSet(obj, "Field1", out0);
		modelmaker::decref(out0);
	}
	{
		modelmaker::JsonValOut obj0 = this->field2.buildJsonObj();
		modelmaker::JsonValOut out0 = obj0;
		modelmaker::objSet(obj, "Field2", out0);
		modelmaker::decref(out0);
	}
	{
		modelmaker::JsonArrayOut out1 = modelmaker::newJsonArray();
		for (unsigned int i = 0; i < this->field3.size(); i++) {
			modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->field3[i]);
			modelmaker::arrayAdd(out1, out0);
			modelmaker::decref(out0);
		}
		modelmaker::objSet(obj, "Field3", out1);
		modelmaker::decref(out1);
	}
	{
		modelmaker::JsonArrayOut out2 = modelmaker::newJsonArray();
		for (unsigned int i = 0; i < this->field4.size(); i++) {
			modelmaker::JsonArrayOut out1 = modelmaker::newJsonArray();
			for (unsigned int ii = 0; ii < this->field4[i].size(); ii++) {
				modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->field4[i][ii]);
				modelmaker::arrayAdd(out1, out0);
				modelmaker::decref(out0);
			}
			modelmaker::arrayAdd(out2, out1);
			modelmaker::decref(out1);
		}
		modelmaker::objSet(obj, "Field4", out2);
		modelmaker::decref(out2);
	}
	{
		modelmaker::JsonObjOut out1 = modelmaker::newJsonObj();
		for (map< string, string >::iterator n = this->field5.begin(); n != this->field5.end(); ++n) {
			std::stringstream s;
			string key;
			std::string tmp;
			s << modelmaker::toStdString(n->first);
			s >> tmp;
			key = modelmaker::toString(tmp);
			modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->field5[n->first]);
			modelmaker::objSet(out1, key, out0);
			modelmaker::decref(out0);
		}
		modelmaker::objSet(obj, "Field5", out1);
		modelmaker::decref(out1);
	}
	return obj;
}
