//Generated Code


#include <fstream>
#include "Model.hpp"

using namespace models;
using namespace models::modelmaker;

bool Model::loadFile(string path) {
	std::ifstream in;
	in.open(modelmaker::toCString(path));
	std::string json;
	if (in.is_open()) {
		while (in.good()) {
			std::string s;
			in >> s;
			json += s;
		}
		in.close();
		load(modelmaker::toString(json));
		return true;
	}
	return false;
}

void Model::writeFile(string path, modelmaker::JsonSerializationSettings sttngs) {
	std::ofstream out;
	out.open(modelmaker::toCString(path));
	std::string json = modelmaker::toStdString(write(sttngs));
	out << json << "\n";
	out.close();
}

void Model::load(string json) {
	modelmaker::JsonValOut obj = modelmaker::read(modelmaker::toCString(json));
	loadJsonObj(obj);
	modelmaker::decref(obj);
}

string Model::write(modelmaker::JsonSerializationSettings sttngs) {
	modelmaker::JsonValOut val = buildJsonObj();
	modelmaker::JsonObjOut obj = modelmaker::toObj(val);
	return modelmaker::write(obj, sttngs);
}

unknown::unknown() {
#ifndef USING_QT
	m_obj = 0;
#endif
}

unknown::unknown(Model *v) {
#ifndef USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(bool v) {
#ifndef USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(int v) {
#ifndef USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(double v) {
#ifndef USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::unknown(string v) {
#ifndef USING_QT
	m_obj = 0;
#endif
	set(v);
}

unknown::~unknown() {
	modelmaker::decref(m_obj);
}

bool unknown::loadJsonObj(modelmaker::JsonVal obj) {
#ifdef USING_JANSSON
	m_obj = modelmaker::incref(obj);
#else
	m_obj = obj;
#endif
	return !modelmaker::isNull(obj);
}

modelmaker::JsonValOut unknown::buildJsonObj() {
#ifdef USING_JANSSON
	return modelmaker::incref(m_obj);
#else
	return m_obj;
#endif
}

bool unknown::loaded() {
	return !modelmaker::isNull(m_obj);
}

bool unknown::isBool() {
	return !modelmaker::isNull(m_obj) && modelmaker::isBool(m_obj);
}

bool unknown::isInt() {
	return !modelmaker::isNull(m_obj) && modelmaker::isInt(m_obj);
}

bool unknown::isDouble() {
	return !modelmaker::isNull(m_obj) && modelmaker::isDouble(m_obj);
}

bool unknown::isString() {
	return !modelmaker::isNull(m_obj) && modelmaker::isString(m_obj);
}

bool unknown::isObject() {
	return !modelmaker::isNull(m_obj) && modelmaker::isObj(m_obj);
}

bool unknown::toBool() {
	return modelmaker::toBool(m_obj);
}

int unknown::toInt() {
	return modelmaker::toInt(m_obj);
}

double unknown::toDouble() {
	return modelmaker::toDouble(m_obj);
}

string unknown::toString() {
	return modelmaker::toString(m_obj);
}

void unknown::set(Model *v) {
	modelmaker::JsonValOut obj = v->buildJsonObj();
	modelmaker::JsonVal old = m_obj;
	m_obj = obj;
	if (!modelmaker::isNull(old)) {
		modelmaker::decref(old);
	}
}

void unknown::set(bool v) {
	modelmaker::JsonValOut obj = modelmaker::toJsonVal(v);
	modelmaker::JsonVal old = m_obj;
	m_obj = obj;
	if (!modelmaker::isNull(old)) {
		modelmaker::decref(old);
	}
}

void unknown::set(int v) {
	modelmaker::JsonValOut obj = modelmaker::toJsonVal(v);
	modelmaker::JsonVal old = m_obj;
	m_obj = obj;
	if (!modelmaker::isNull(old)) {
		modelmaker::decref(old);
	}
}

void unknown::set(double v) {
	modelmaker::JsonValOut obj = modelmaker::toJsonVal(v);
	modelmaker::JsonVal old = m_obj;
	m_obj = obj;
	if (!modelmaker::isNull(old)) {
		modelmaker::decref(old);
	}
}

void unknown::set(string v) {
	modelmaker::JsonValOut obj = modelmaker::toJsonVal(v);
	modelmaker::JsonVal old = m_obj;
	m_obj = obj;
	if (!modelmaker::isNull(old)) {
		modelmaker::decref(old);
	}
}


#include "string.h"
#include "Model.hpp"

using namespace models;
using std::stringstream;

Model1::Model1() {
	this->field1 = "";
	for (int i = 0; i < 4; this->field3[i++] = 0);
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
				modelmaker::JsonValOut obj1 = modelmaker::iteratorValue(it1);
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
		for (unsigned int i = 0; i < 4; i++) {
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
		for (std::map< string, string >::iterator n = this->field5.begin(); n != this->field5.end(); ++n) {
			std::stringstream s;
			string key;
			std::string tmp;
			s << modelmaker::toStdString(modelmaker::toString(n->first));
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
