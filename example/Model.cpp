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
#ifndef CYBORGBEAR_USING_QT
	m_obj = 0;
#endif
}

unknown::unknown(Model *v) {
#ifndef CYBORGBEAR_USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(bool v) {
#ifndef CYBORGBEAR_USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(int v) {
#ifndef CYBORGBEAR_USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(double v) {
#ifndef CYBORGBEAR_USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(string v) {
#ifndef CYBORGBEAR_USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::~unknown() {
	cyborgbear::decref(m_obj);
}

bool unknown::loadJsonObj(cyborgbear::JsonVal obj) {
#ifdef CYBORGBEAR_USING_JANSSON
	m_obj = cyborgbear::incref(obj);
#else
	m_obj = obj;
#endif
	return !cyborgbear::isNull(obj);
}

cyborgbear::JsonValOut unknown::buildJsonObj() {
#ifdef CYBORGBEAR_USING_JANSSON
	return cyborgbear::incref(m_obj);
#else
	return m_obj;
#endif
}

bool unknown::loaded() {
	return !cyborgbear::isNull(m_obj);
}

bool unknown::isBool() {
	return !cyborgbear::isNull(m_obj) && cyborgbear::isBool(m_obj);
}

bool unknown::isInt() {
	return !cyborgbear::isNull(m_obj) && cyborgbear::isInt(m_obj);
}

bool unknown::isDouble() {
	return !cyborgbear::isNull(m_obj) && cyborgbear::isDouble(m_obj);
}

bool unknown::isString() {
	return !cyborgbear::isNull(m_obj) && cyborgbear::isString(m_obj);
}

bool unknown::isObject() {
	return !cyborgbear::isNull(m_obj) && cyborgbear::isObj(m_obj);
}

bool unknown::toBool() {
	return cyborgbear::toBool(m_obj);
}

int unknown::toInt() {
	return cyborgbear::toInt(m_obj);
}

double unknown::toDouble() {
	return cyborgbear::toDouble(m_obj);
}

string unknown::toString() {
	return cyborgbear::toString(m_obj);
}

void unknown::set(Model *v) {
	cyborgbear::JsonValOut obj = v->buildJsonObj();
	cyborgbear::JsonVal old = m_obj;
	m_obj = obj;
	if (!cyborgbear::isNull(old)) {
		cyborgbear::decref(old);
	}
}

void unknown::set(bool v) {
	cyborgbear::JsonValOut obj = cyborgbear::toJsonVal(v);
	cyborgbear::JsonVal old = m_obj;
	m_obj = obj;
	if (!cyborgbear::isNull(old)) {
		cyborgbear::decref(old);
	}
}

void unknown::set(int v) {
	cyborgbear::JsonValOut obj = cyborgbear::toJsonVal(v);
	cyborgbear::JsonVal old = m_obj;
	m_obj = obj;
	if (!cyborgbear::isNull(old)) {
		cyborgbear::decref(old);
	}
}

void unknown::set(double v) {
	cyborgbear::JsonValOut obj = cyborgbear::toJsonVal(v);
	cyborgbear::JsonVal old = m_obj;
	m_obj = obj;
	if (!cyborgbear::isNull(old)) {
		cyborgbear::decref(old);
	}
}

void unknown::set(string v) {
	cyborgbear::JsonValOut obj = cyborgbear::toJsonVal(v);
	cyborgbear::JsonVal old = m_obj;
	m_obj = obj;
	if (!cyborgbear::isNull(old)) {
		cyborgbear::decref(old);
	}
}


#include "string.h"
#include "Model.hpp"

using namespace models;
using std::stringstream;

Model1::Model1() {
	this->Field1 = "";
	for (int i = 0; i < 4; this->Field3[i++] = 0);
}

bool Model1::loadJsonObj(cyborgbear::JsonVal in) {
	cyborgbear::JsonObjOut inObj = cyborgbear::toObj(in);
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field1");
		{
			if (cyborgbear::isString(obj0)) {
				this->Field1 = cyborgbear::toString(obj0);
			}
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field2");
		{
			if (cyborgbear::isInt(obj0)) {
				this->Field2 = cyborgbear::toInt(obj0);
			}
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
						this->Field3[i] = cyborgbear::toInt(obj1);
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
			this->Field4.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				cyborgbear::JsonValOut obj1 = cyborgbear::arrayRead(array0, i);
				if (!cyborgbear::isNull(obj1) && cyborgbear::isArray(obj1)) {
					cyborgbear::JsonArrayOut array1 = cyborgbear::toArray(obj1);
					unsigned int size = cyborgbear::arraySize(array1);
					this->Field4[i].resize(size);
					for (unsigned int ii = 0; ii < size; ii++) {
						cyborgbear::JsonValOut obj2 = cyborgbear::arrayRead(array1, ii);
						{
							if (cyborgbear::isString(obj2)) {
								this->Field4[i][ii] = cyborgbear::toString(obj2);
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
						this->Field5[i] = cyborgbear::toString(obj1);
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
		cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->Field1);
		cyborgbear::objSet(obj, "Field1", out0);
		cyborgbear::decref(out0);
	}
	{
		cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->Field2);
		cyborgbear::objSet(obj, "Field2", out0);
		cyborgbear::decref(out0);
	}
	{
		cyborgbear::JsonArrayOut out1 = cyborgbear::newJsonArray();
		for (cyborgbear::VectorIterator i = 0; i < 4; i++) {
			cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->Field3[i]);
			cyborgbear::arrayAdd(out1, out0);
			cyborgbear::decref(out0);
		}
		cyborgbear::objSet(obj, "Field3", out1);
		cyborgbear::decref(out1);
	}
	{
		cyborgbear::JsonArrayOut out2 = cyborgbear::newJsonArray();
		for (cyborgbear::VectorIterator i = 0; i < this->Field4.size(); i++) {
			cyborgbear::JsonArrayOut out1 = cyborgbear::newJsonArray();
			for (cyborgbear::VectorIterator ii = 0; ii < this->Field4[i].size(); ii++) {
				cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->Field4[i][ii]);
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
		for (std::map< string, string >::iterator n = this->Field5.begin(); n != this->Field5.end(); ++n) {
			std::stringstream s;
			string key;
			std::string tmp;
			s << cyborgbear::toStdString(cyborgbear::toString(n->first));
			s >> tmp;
			key = cyborgbear::toString(tmp);
			cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->Field5[n->first]);
			cyborgbear::objSet(out1, key, out0);
			cyborgbear::decref(out0);
		}
		cyborgbear::objSet(obj, "Field5", out1);
		cyborgbear::decref(out1);
	}
	return obj;
}
