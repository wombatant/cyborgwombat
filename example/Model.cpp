//Generated Code

#include <fstream>
#include "Model.hpp"

using namespace models;
using namespace models::cyborgbear;

bool Model::readJsonFile(string path) {
	std::ifstream in;
	in.open(cyborgbear::toCString(path));
	std::string json;
	if (in.is_open()) {
		while (in.good()) {
			std::string s;
			in >> s;
			json += s;
		}
		in.close();
		fromJson(cyborgbear::toString(json));
		return true;
	}
	return false;
}

void Model::writeJsonFile(string path, cyborgbear::JsonSerializationSettings sttngs) {
	std::ofstream out;
	out.open(cyborgbear::toCString(path));
	std::string json = cyborgbear::toStdString(toJson(sttngs));
	out << json << "\n";
	out.close();
}

void Model::fromJson(string json) {
	cyborgbear::JsonValOut obj = cyborgbear::read(json);
	loadJsonObj(obj);
	cyborgbear::decref(obj);
}

string Model::toJson(cyborgbear::JsonSerializationSettings sttngs) {
	cyborgbear::JsonValOut val = buildJsonObj();
	cyborgbear::JsonObjOut obj = cyborgbear::toObj(val);
	return cyborgbear::write(obj, sttngs);
}

unknown::unknown() {
}

unknown::unknown(Model *v) {
	set(v);
}

unknown::unknown(bool v) {
	set(v);
}

unknown::unknown(int v) {
	set(v);
}

unknown::unknown(double v) {
	set(v);
}

unknown::unknown(string v) {
	set(v);
}

unknown::~unknown() {
}

bool unknown::loadJsonObj(cyborgbear::JsonVal obj) {
	cyborgbear::JsonValOut wrapper = cyborgbear::newJsonObj();
	cyborgbear::objSet(wrapper, "Value", obj);
	m_data = cyborgbear::write(wrapper, cyborgbear::Compact);
	printf("loadJsonObj: %s\n", m_data.c_str());
	if (cyborgbear::isBool(obj)) {
		m_type = cyborgbear::Bool;
	} else if (cyborgbear::isInt(obj)) {
		m_type = cyborgbear::Integer;
	} else if (cyborgbear::isDouble(obj)) {
		m_type = cyborgbear::Double;
	} else if (cyborgbear::isString(obj)) {
		m_type = cyborgbear::String;
	} else if (cyborgbear::isObj(obj)) {
		m_type = cyborgbear::Object;
	}

	return !cyborgbear::isNull(obj);
}

cyborgbear::JsonValOut unknown::buildJsonObj() {
	cyborgbear::JsonValOut obj = cyborgbear::read(m_data);
	cyborgbear::JsonValOut val = cyborgbear::incref(cyborgbear::objRead(obj, "Value"));
	cyborgbear::decref(obj);
	return val;
}

bool unknown::loaded() {
	return m_data != "";
}

bool unknown::isBool() {
	return m_type == cyborgbear::Bool;
}

bool unknown::isInt() {
	return m_type == cyborgbear::Integer;
}

bool unknown::isDouble() {
	return m_type == cyborgbear::Double;
}

bool unknown::isString() {
	return m_type == cyborgbear::String;
}

bool unknown::isObject() {
	return m_type == cyborgbear::Object;
}

bool unknown::toBool() {
	return cyborgbear::toBool(buildJsonObj());
}

int unknown::toInt() {
	return cyborgbear::toInt(buildJsonObj());
}

double unknown::toDouble() {
	return cyborgbear::toDouble(buildJsonObj());
}

string unknown::toString() {
	return cyborgbear::toString(buildJsonObj());
}

void unknown::set(Model *v) {
	cyborgbear::JsonValOut obj = cyborgbear::newJsonObj();
	cyborgbear::objSet(obj, "Value", v->buildJsonObj());
	m_type = cyborgbear::Object;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

void unknown::set(bool v) {
	cyborgbear::JsonValOut obj = cyborgbear::newJsonObj();
	cyborgbear::objSet(obj, "Value", cyborgbear::toJsonVal(v));
	m_type = cyborgbear::Bool;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

void unknown::set(int v) {
	cyborgbear::JsonValOut obj = cyborgbear::newJsonObj();
	cyborgbear::objSet(obj, "Value", cyborgbear::toJsonVal(v));
	m_type = cyborgbear::Integer;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

void unknown::set(double v) {
	cyborgbear::JsonValOut obj = cyborgbear::newJsonObj();
	cyborgbear::objSet(obj, "Value", cyborgbear::toJsonVal(v));
	m_type = cyborgbear::Double;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

void unknown::set(string v) {
	cyborgbear::JsonValOut obj = cyborgbear::newJsonObj();
	cyborgbear::objSet(obj, "Value", cyborgbear::toJsonVal(v));
	m_type = cyborgbear::String;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}


#include "string.h"
#include "Model.hpp"

using namespace models;
using std::stringstream;

Model1::Model1() {
	this->field1 = "";
	for (int i = 0; i < 4; this->field3[i++] = 0);
}

bool Model1::loadJsonObj(cyborgbear::JsonVal in) {
	cyborgbear::JsonObjOut inObj = cyborgbear::toObj(in);
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field1");
		{
			if (cyborgbear::isString(obj0)) {
				this->field1 = cyborgbear::toString(obj0);
			}
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field2");
		{
			this->field2.loadJsonObj(obj0);
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field3");
		if (!cyborgbear::isNull(obj0) && cyborgbear::isArray(obj0)) {
			cyborgbear::JsonArrayOut array0 = cyborgbear::toArray(obj0);
			unsigned int size = cyborgbear::arraySize(array0);
			for (unsigned int i = 0; i < size; i++) {
				cyborgbear::JsonValOut obj1 = cyborgbear::arrayRead(array0, i);
				{
					if (cyborgbear::isInt(obj1)) {
						this->field3[i] = cyborgbear::toInt(obj1);
					}
				}
			}
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field4");
		if (!cyborgbear::isNull(obj0) && cyborgbear::isArray(obj0)) {
			cyborgbear::JsonArrayOut array0 = cyborgbear::toArray(obj0);
			unsigned int size = cyborgbear::arraySize(array0);
			this->field4.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				cyborgbear::JsonValOut obj1 = cyborgbear::arrayRead(array0, i);
				if (!cyborgbear::isNull(obj1) && cyborgbear::isArray(obj1)) {
					cyborgbear::JsonArrayOut array1 = cyborgbear::toArray(obj1);
					unsigned int size = cyborgbear::arraySize(array1);
					this->field4[i].resize(size);
					for (unsigned int ii = 0; ii < size; ii++) {
						cyborgbear::JsonValOut obj2 = cyborgbear::arrayRead(array1, ii);
						{
							if (cyborgbear::isString(obj2)) {
								this->field4[i][ii] = cyborgbear::toString(obj2);
							}
						}
					}
				}
			}
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field5");
		if (!cyborgbear::isNull(obj0) && cyborgbear::isObj(obj0)) {
			cyborgbear::JsonObjOut map0 = cyborgbear::toObj(obj0);
			for (cyborgbear::JsonObjIterator it1 = cyborgbear::jsonObjIterator(map0); !cyborgbear::iteratorAtEnd(it1, map0); it1 = cyborgbear::jsonObjIteratorNext(map0,  it1)) {
				string i;
				cyborgbear::JsonValOut obj1 = cyborgbear::iteratorValue(it1);
				{
					std::string key = cyborgbear::toStdString(cyborgbear::jsonObjIteratorKey(it1));
					std::string o;
					std::stringstream s;
					s << key;
					s >> o;
					i = o.c_str();
				}
				{
					if (cyborgbear::isString(obj1)) {
						this->field5[i] = cyborgbear::toString(obj1);
					}
				}
			}
		}
	}
	return true;
}

cyborgbear::JsonValOut Model1::buildJsonObj() {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();
	{
		cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->field1);
		cyborgbear::objSet(obj, "Field1", out0);
		cyborgbear::decref(out0);
	}
	{
		cyborgbear::JsonValOut obj0 = this->field2.buildJsonObj();
		cyborgbear::JsonValOut out0 = obj0;
		cyborgbear::objSet(obj, "Field2", out0);
		cyborgbear::decref(out0);
	}
	{
		cyborgbear::JsonArrayOut out1 = cyborgbear::newJsonArray();
		for (cyborgbear::VectorIterator i = 0; i < 4; i++) {
			cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->field3[i]);
			cyborgbear::arrayAdd(out1, out0);
			cyborgbear::decref(out0);
		}
		cyborgbear::objSet(obj, "Field3", out1);
		cyborgbear::decref(out1);
	}
	{
		cyborgbear::JsonArrayOut out2 = cyborgbear::newJsonArray();
		for (cyborgbear::VectorIterator i = 0; i < this->field4.size(); i++) {
			cyborgbear::JsonArrayOut out1 = cyborgbear::newJsonArray();
			for (cyborgbear::VectorIterator ii = 0; ii < this->field4[i].size(); ii++) {
				cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->field4[i][ii]);
				cyborgbear::arrayAdd(out1, out0);
				cyborgbear::decref(out0);
			}
			cyborgbear::arrayAdd(out2, out1);
			cyborgbear::decref(out1);
		}
		cyborgbear::objSet(obj, "Field4", out2);
		cyborgbear::decref(out2);
	}
	{
		cyborgbear::JsonObjOut out1 = cyborgbear::newJsonObj();
		for (std::map< string, string >::iterator n = this->field5.begin(); n != this->field5.end(); ++n) {
			std::stringstream s;
			string key;
			std::string tmp;
			s << cyborgbear::toStdString(cyborgbear::toString(n->first));
			s >> tmp;
			key = cyborgbear::toString(tmp);
			cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->field5[n->first]);
			cyborgbear::objSet(out1, key, out0);
			cyborgbear::decref(out0);
		}
		cyborgbear::objSet(obj, "Field5", out1);
		cyborgbear::decref(out1);
	}
	return obj;
}
