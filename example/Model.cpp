//Generated Code

#include <fstream>
#include "Model.hpp"

using namespace models;
using namespace models::cyborgbear;

string models::cyborgbear::version = "1.0.0-beta8";

int Model::readJsonFile(string path) {
	std::ifstream in;
	in.open(cyborgbear::toCString(path));
	if (in.is_open()) {
		std::string json((std::istreambuf_iterator<char>(in)), std::istreambuf_iterator<char>());
		in.close();
		return fromJson(cyborgbear::toString(json));
	}
	return cyborgbear::Error_CouldNotAccessFile;
}

void Model::writeJsonFile(string path, cyborgbear::JsonSerializationSettings sttngs) {
	std::ofstream out;
	out.open(cyborgbear::toCString(path));
	std::string json = cyborgbear::toStdString(toJson(sttngs));
	out << json << "\0";
	out.close();
}

int Model::fromJson(string json) {
	cyborgbear::JsonValOut obj = cyborgbear::read(json);
	cyborgbear::Error retval = loadJsonObj(obj);
	cyborgbear::decref(obj);
	return retval;
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

cyborgbear::Error unknown::loadJsonObj(cyborgbear::JsonVal obj) {
	cyborgbear::JsonObjOut wrapper = cyborgbear::newJsonObj();
	cyborgbear::objSet(wrapper, "Value", obj);
	m_data = cyborgbear::write(wrapper, cyborgbear::Compact);
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

	if (cyborgbear::isNull(obj)) {
		return cyborgbear::Error_GenericParsingError;
	} else {
		return cyborgbear::Error_Ok;
	}
}

cyborgbear::JsonValOut unknown::buildJsonObj() {
	cyborgbear::JsonObjOut obj = cyborgbear::read(m_data);
#ifdef CYBORGBEAR_USING_QT
	cyborgbear::JsonValOut val = cyborgbear::objRead(obj, "Value");
#else
	cyborgbear::JsonValOut val = cyborgbear::incref(cyborgbear::objRead(obj, "Value"));
#endif
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
	cyborgbear::JsonValOut out = buildJsonObj();
	return cyborgbear::toBool(out);
}

int unknown::toInt() {
	cyborgbear::JsonValOut out = buildJsonObj();
	return cyborgbear::toInt(out);
}

double unknown::toDouble() {
	cyborgbear::JsonValOut out = buildJsonObj();
	return cyborgbear::toDouble(out);
}

string unknown::toString() {
	cyborgbear::JsonValOut out = buildJsonObj();
	return cyborgbear::toString(out);
}

void unknown::set(Model *v) {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();
	cyborgbear::JsonValOut val = v->buildJsonObj();
	cyborgbear::objSet(obj, "Value", val);
	m_type = cyborgbear::Object;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);

	unknown *unk = dynamic_cast<unknown*>(v);
	if (unk)
		m_type = unk->m_type;
}

void unknown::set(bool v) {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();
	cyborgbear::JsonValOut val = cyborgbear::toJsonVal(v);
	cyborgbear::objSet(obj, "Value", val);
	m_type = cyborgbear::Bool;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

void unknown::set(int v) {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();
	cyborgbear::JsonValOut val = cyborgbear::toJsonVal(v);
	cyborgbear::objSet(obj, "Value", val);
	m_type = cyborgbear::Integer;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

void unknown::set(double v) {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();
	cyborgbear::JsonValOut val = cyborgbear::toJsonVal(v);
	cyborgbear::objSet(obj, "Value", val);
	m_type = cyborgbear::Double;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

void unknown::set(string v) {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();
	cyborgbear::JsonValOut val = cyborgbear::toJsonVal(v);
	cyborgbear::objSet(obj, "Value", val);
	m_type = cyborgbear::String;
	m_data = cyborgbear::write(obj, cyborgbear::Compact);
	cyborgbear::decref(obj);
}

#ifdef CYBORGBEAR_BOOST_ENABLED

void unknown::fromBoostBinary(string dat) {
	std::stringstream in(dat);
	boost::archive::binary_iarchive ia(in);
	ia >> *this;
}

string unknown::toBoostBinary() {
	std::stringstream out;
	{
		boost::archive::binary_oarchive oa(out);
		oa << *this;
	}
	string str;
	while (out.good())
		str += out.get();
	return str;
}
#endif


#include "string.h"
#include "Model.hpp"

using namespace models;
using std::stringstream;

Model1::Model1() {
	this->Field1 = "";
	for (int i = 0; i < 4; this->Field3[i++] = 0);
}

cyborgbear::Error Model1::loadJsonObj(cyborgbear::JsonVal in) {
	cyborgbear::Error retval = cyborgbear::Error_Ok;
	cyborgbear::JsonObjOut inObj = cyborgbear::toObj(in);

	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field1");
		{
			if (cyborgbear::isString(obj0)) {
				this->Field1 = cyborgbear::toString(obj0);
			} else {
				if (cyborgbear::isNull(obj0)) {
					retval |= cyborgbear::Error_MissingField;
				} else {
					retval |= cyborgbear::Error_TypeMismatch;
				}
			}
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field2");
		{
			retval |= this->Field2.loadJsonObj(obj0);
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field3");
		if (!cyborgbear::isNull(obj0)) {
			if (cyborgbear::isArray(obj0)) {
				cyborgbear::JsonArrayOut array0 = cyborgbear::toArray(obj0);
				unsigned int size = cyborgbear::arraySize(array0);
				for (unsigned int i = 0; i < size; i++) {
					cyborgbear::JsonValOut obj1 = cyborgbear::arrayRead(array0, i);
					{
						if (cyborgbear::isInt(obj1)) {
							this->Field3[i] = cyborgbear::toInt(obj1);
						} else {
							if (cyborgbear::isNull(obj1)) {
								retval |= cyborgbear::Error_MissingField;
							} else {
								retval |= cyborgbear::Error_TypeMismatch;
							}
						}
					}
				}
			} else {
				retval |= cyborgbear::Error_TypeMismatch;
			}
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field4");
		if (!cyborgbear::isNull(obj0)) {
			if (cyborgbear::isArray(obj0)) {
				cyborgbear::JsonArrayOut array0 = cyborgbear::toArray(obj0);
				unsigned int size = cyborgbear::arraySize(array0);
				this->Field4.resize(size);
				for (unsigned int i = 0; i < size; i++) {
					cyborgbear::JsonValOut obj1 = cyborgbear::arrayRead(array0, i);
					if (!cyborgbear::isNull(obj1)) {
						if (cyborgbear::isArray(obj1)) {
							cyborgbear::JsonArrayOut array1 = cyborgbear::toArray(obj1);
							unsigned int size = cyborgbear::arraySize(array1);
							this->Field4[i].resize(size);
							for (unsigned int ii = 0; ii < size; ii++) {
								cyborgbear::JsonValOut obj2 = cyborgbear::arrayRead(array1, ii);
								{
									if (cyborgbear::isString(obj2)) {
										this->Field4[i][ii] = cyborgbear::toString(obj2);
									} else {
										if (cyborgbear::isNull(obj2)) {
											retval |= cyborgbear::Error_MissingField;
										} else {
											retval |= cyborgbear::Error_TypeMismatch;
										}
									}
								}
							}
						} else {
							retval |= cyborgbear::Error_TypeMismatch;
						}
					}
				}
			} else {
				retval |= cyborgbear::Error_TypeMismatch;
			}
		}
	}
	{
		cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, "Field5");
		if (!cyborgbear::isNull(obj0)) {
			if (cyborgbear::isObj(obj0)) {
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
						} else {
							if (cyborgbear::isNull(obj1)) {
								retval |= cyborgbear::Error_MissingField;
							} else {
								retval |= cyborgbear::Error_TypeMismatch;
							}
						}
					}
				}
			}
		}
	}
	return retval;
}

cyborgbear::JsonValOut Model1::buildJsonObj() {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();
	{
		cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->Field1);
		cyborgbear::objSet(obj, "Field1", out0);
		cyborgbear::decref(out0);
	}
	{
		cyborgbear::JsonValOut obj0 = this->Field2.buildJsonObj();
		cyborgbear::JsonValOut out0 = obj0;
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

namespace models {

#ifdef CYBORGBEAR_BOOST_ENABLED
void Model1::fromBoostBinary(string dat) {
	std::stringstream in(dat);
	boost::archive::binary_iarchive ia(in);
	ia >> *this;
}

string Model1::toBoostBinary() {
	std::stringstream out;
	{
		boost::archive::binary_oarchive oa(out);
		oa << *this;
	}

	string str;
	while (out.good())
		str += out.get();

	return str;
}
#endif
}
