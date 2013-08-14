//Generated Code

#include <fstream>
#include "modelmakerdefs.hpp"

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

void Model::writeFile(string path) {
	std::ofstream out;
	out.open(modelmaker::toCString(path));
	std::string json = modelmaker::toStdString(write());
	out << json << "\n";
	out.close();
}

void Model::load(string json) {
	modelmaker::JsonValOut obj = modelmaker::read(modelmaker::toCString(json));
	loadJsonObj(obj);
	modelmaker::decref(obj);
}

string Model::write() {
	modelmaker::JsonValOut val = buildJsonObj();
	modelmaker::JsonObjOut obj = modelmaker::toObj(val);
	return modelmaker::write(obj);
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
	return modelmaker::isNull(m_obj);
}

bool unknown::isBool() {
	return modelmaker::isNull(m_obj) && modelmaker::isBool(m_obj);
}

bool unknown::isInt() {
	return modelmaker::isNull(m_obj) && modelmaker::isInt(m_obj);
}

bool unknown::isDouble() {
	return modelmaker::isNull(m_obj) && modelmaker::isDouble(m_obj);
}

bool unknown::isString() {
	return modelmaker::isNull(m_obj) && modelmaker::isString(m_obj);
}

bool unknown::isObject() {
	return modelmaker::isNull(m_obj) && modelmaker::isObj(m_obj);
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
