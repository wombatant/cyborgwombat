//Generated Code

#ifndef MODEL_HPP
#define MODEL_HPP
#include <string>
#include <sstream>

#define CYBORGBEAR_USING_JANSSON
#define CYBORGBEAR_BOOST_ENABLED


#ifdef CYBORGBEAR_USING_QT
#include <QString>
#include <QJsonDocument>
#include <QJsonArray>
#include <QJsonObject>
#include <QJsonValue>
#include <QMap>
#include <QVector>
#else
#include <vector>
#include <map>
#include <string>
#include <jansson.h>
#endif

#ifdef CYBORGBEAR_BOOST_ENABLED
#include <boost/serialization/vector.hpp>
#include <boost/serialization/map.hpp>
#include <boost/serialization/string.hpp>
#include <boost/archive/binary_oarchive.hpp>
#include <boost/archive/binary_iarchive.hpp>
#endif

namespace models {

namespace cyborgbear {

typedef unsigned long int Error;
const Error Error_Ok = 0;
const Error Error_TypeMismatch = 1;
const Error Error_MissingField = 2;
const Error Error_CouldNotAccessFile = 4;
const Error Error_GenericParsingError = 8;

enum JsonSerializationSettings {
	Compact = 0,
	Readable = 1
};

enum Type {
	Bool,
	Integer,
	Double,
	String,
	Object
};

#ifdef CYBORGBEAR_USING_QT
typedef QJsonObject& JsonObj;
typedef QJsonValue&  JsonVal;
typedef QJsonArray&  JsonArray;

typedef QJsonObject  JsonObjOut;
typedef QJsonValue   JsonValOut;
typedef QJsonArray   JsonArrayOut;

typedef QJsonObject::iterator JsonObjIterator;
typedef QString               JsonObjIteratorKey;
typedef QJsonValueRef         JsonObjIteratorVal;

typedef QString string;

typedef int VectorIterator;

#else

typedef json_t* JsonObj;
typedef json_t* JsonVal;
typedef json_t* JsonArray;

typedef json_t* JsonObjOut;
typedef json_t* JsonValOut;
typedef json_t* JsonArrayOut;

typedef const char* JsonObjIterator;
typedef const char* JsonObjIteratorKey;
typedef json_t*     JsonObjIteratorVal;

typedef std::string string;

typedef unsigned VectorIterator;
#endif

/**
 * Version of cyborgbear.
 */
extern string version;

//string ops
std::string toStdString(string str);


JsonObjOut read(string json);

cyborgbear::Error readVal(JsonVal, int&);
cyborgbear::Error readVal(JsonVal, double&);
cyborgbear::Error readVal(JsonVal, bool&);
cyborgbear::Error readVal(JsonVal, string&);
JsonArrayOut toArray(JsonVal);
JsonObjOut toObj(JsonVal);

JsonValOut toJsonVal(int);
JsonValOut toJsonVal(double);
JsonValOut toJsonVal(bool);
JsonValOut toJsonVal(string);
JsonValOut toJsonVal(JsonArray);
JsonValOut toJsonVal(JsonObj);


//value methods

template<typename T>
bool isBool(T);

template<typename T>
bool isInt(T);

template<typename T>
bool isDouble(T);

template<typename T>
bool isString(T);

template<typename T>
bool isArray(T);

template<typename T>
bool isObj(T);

template<typename T>
bool isNull(T v);


JsonObj incref(JsonObj);
JsonVal incref(JsonVal);
JsonArray incref(JsonArray);

void decref(JsonObj);
void decref(JsonVal);
void decref(JsonArray);


JsonArrayOut newJsonArray();

void arrayAdd(JsonArray, JsonObj);
void arrayAdd(JsonArray, JsonVal);
void arrayAdd(JsonArray, JsonArray);

int arraySize(JsonArray);

JsonValOut arrayRead(JsonArray, int);


JsonObjOut newJsonObj();

void objSet(JsonObj, string, JsonObj);
void objSet(JsonObj, string, JsonVal);
void objSet(JsonObj, string, JsonArray);

JsonValOut objRead(JsonObj, string);


JsonObjIterator jsonObjIterator(JsonObj);
JsonObjIterator jsonObjIteratorNext(JsonObj, JsonObjIterator);
JsonObjIteratorKey jsonObjIteratorKey(JsonObjIterator);
JsonObjIteratorVal iteratorValue(JsonObjIterator);
bool iteratorAtEnd(JsonObjIterator, JsonObj);



inline string toString(string str) {
	return str;
}


#ifdef CYBORGBEAR_USING_QT

//string conversions
inline std::string toStdString(string str) {
	return str.toStdString();
}

inline string toString(std::string str) {
	return QString::fromStdString(str);
}


inline JsonObjOut read(string json) {
	return QJsonDocument::fromJson(json.toUtf8()).object();
}


//from JsonObj or JsonObjIteratorVal
template<typename T>
inline cyborgbear::Error readVal(T v, int &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isInt(v)) {
		val = v.toInt();
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

template<typename T>
inline cyborgbear::Error readVal(T v, double &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isDouble(v)) {
		val = v.toDouble();
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

template<typename T>
inline cyborgbear::Error readVal(T v, bool &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isBool(v)) {
		val = v.toBool();
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

template<typename T>
inline cyborgbear::Error readVal(T v, string &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isString(v)) {
		val = v.toString();
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

inline int toArray(T v, JsonArrayOut &val) {
	val = v.toArray();
	return 0;
}

inline int toObj(T v, JsonArrayOut &val) {
	val = v.toObject();
	return 0;
}


inline JsonValOut toJsonVal(int v) {
	return QJsonValue(v);
}

inline JsonValOut toJsonVal(double v) {
	return QJsonValue(v);
}

inline JsonValOut toJsonVal(bool v) {
	return QJsonValue(v);
}

inline JsonValOut toJsonVal(string v) {
	return QJsonValue(v);
}

inline JsonValOut toJsonVal(JsonArray v) {
	return QJsonValue(v);
}

inline JsonValOut toJsonVal(JsonObj v) {
	return v;
}


template<typename T>
inline bool isNull(T v) {
	return v.isNull();
}

template<typename T>
inline bool isBool(T v) {
	return v.isBool();
}

template<typename T>
inline bool isInt(T v) {
	return v.isDouble();
}

template<typename T>
inline bool isDouble(T v) {
	return v.isDouble();
}

template<typename T>
inline bool isString(T v) {
	return v.isString();
}

template<typename T>
inline bool isArray(T v) {
	return v.isArray();
}

template<typename T>
inline bool isObj(T v) {
	return v.isObject();
}


inline JsonVal incref(JsonVal v) {
	return v;
}

inline void decref(JsonVal) {}

inline JsonObj incref(JsonObj v) {
	return v;
}

inline void decref(JsonObj) {}

inline JsonArray incref(JsonArray v) {
	return v;
}

inline void decref(JsonArray) {}


inline JsonArrayOut newJsonArray() {
	return QJsonArray();
}

inline void arrayAdd(JsonArray a, JsonArray v) {
	JsonValOut tmp = cyborgbear::toJsonVal(v);
	a.append(tmp);
}

inline void arrayAdd(JsonArray a, JsonObj v) {
	JsonValOut tmp = cyborgbear::toJsonVal(v);
	a.append(tmp);
}

inline void arrayAdd(JsonArray a, JsonVal v) {
	a.append(v);
}


inline JsonValOut arrayRead(JsonArray a, int i) {
	return a[i];
}

inline int arraySize(JsonArray a) {
	return a.size();
}


inline JsonObjOut newJsonObj() {
	return QJsonObject();
}

inline void objSet(JsonObj o, string k, JsonArray v) {
	JsonValOut tmp = cyborgbear::toJsonVal(v);
	o.insert(k, tmp);
}

inline void objSet(JsonObj o, string k, JsonObj v) {
	JsonValOut tmp = cyborgbear::toJsonVal(v);
	o.insert(k, tmp);
}

inline void objSet(JsonObj o, string k, JsonVal v) {
	o.insert(k, v);
}


inline JsonValOut objRead(JsonObj o, string k) {
	return o[k];
}

inline JsonObjIterator jsonObjIterator(JsonObj o) {
	return o.begin();
}

inline JsonObjIterator jsonObjIteratorNext(JsonObj, JsonObjIterator i) {
	return i + 1;
}

inline JsonObjIteratorKey jsonObjIteratorKey(JsonObjIterator i) {
	return i.key();
}

inline bool iteratorAtEnd(JsonObjIterator i, JsonObj o) {
	return i == o.end();
}

inline JsonObjIteratorVal iteratorValue(JsonObjIterator i) {
	return i.value();
}

inline string write(JsonObj obj, JsonSerializationSettings sttngs) {
	QJsonDocument doc(obj);
	return doc.toJson(sttngs == Compact ? QJsonDocument::Compact : QJsonDocument::Indented);
}

#else

inline std::string toStdString(string str) {
	return str;
}


inline JsonObjOut read(string json) {
	return json_loads(json.c_str(), 0, NULL);
}

inline string write(JsonObj obj, JsonSerializationSettings sttngs) {
	char *tmp = json_dumps(obj, sttngs == Compact ? JSON_COMPACT : JSON_INDENT(3));
	if (!tmp)
		return "{}";
	string out = tmp;
	free(tmp);
	cyborgbear::decref(obj);
	return out;
}

//value methods

inline cyborgbear::Error readVal(JsonVal v, int &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isInt(v)) {
		val = (int) json_integer_value(v);
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

inline cyborgbear::Error readVal(JsonVal v, double &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isDouble(v)) {
		val = (double) json_real_value(v);
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

inline cyborgbear::Error readVal(JsonVal v, bool &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isBool(v)) {
		val = json_is_true(v);
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

inline cyborgbear::Error readVal(JsonVal v, string &val) {
	int retval = cyborgbear::Error_Ok;
	if (cyborgbear::isString(v)) {
		val = json_string_value(v);
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

inline JsonArray toArray(JsonVal v) {
	return v;
}

inline JsonObj toObj(JsonVal v) {
	return v;
}


inline JsonVal toJsonVal(int v) {
	return json_integer(v);
}

inline JsonVal toJsonVal(double v) {
	return json_real(v);
}

inline JsonVal toJsonVal(bool v) {
	return json_boolean(v);
}

inline JsonVal toJsonVal(string v) {
	return json_string(v.c_str());
}

inline JsonVal toJsonVal(JsonArray v) {
	return v;
}


template<typename T>
inline bool isNull(T v) {
	return !v;
}

template<typename T>
inline bool isBool(T v) {
	return json_is_boolean(v);
}

template<typename T>
inline bool isInt(T v) {
	return json_is_integer(v);
}

template<typename T>
inline bool isDouble(T v) {
	return json_is_real(v);
}

template<typename T>
inline bool isString(T v) {
	return json_is_string(v);
}

template<typename T>
inline bool isArray(T v) {
	return json_is_array(v);
}

template<typename T>
inline bool isObj(T v) {
	return json_is_object(v);
}

inline JsonVal incref(JsonVal v) {
	return json_incref(v);
}

inline void decref(JsonVal v) {
	json_decref(v);
}

//array methods

inline JsonArrayOut newJsonArray() {
	return json_array();
}

inline void arrayAdd(JsonArray a, JsonVal v) {
	json_array_append(a, v);
}

inline JsonVal arrayRead(JsonArray a, int i) {
	return json_array_get(a, i);
}

inline int arraySize(JsonArray a) {
	return json_array_size(a);
}

//object methods

inline JsonObjOut newJsonObj() {
	return json_object();
}

inline void objSet(JsonObj o, string k, JsonVal v) {
	json_object_set(o, k.c_str(), v);
}

inline JsonVal objRead(JsonObj o, string k) {
	return json_object_get(o, k.c_str());
}


inline JsonObjIterator jsonObjIterator(JsonObj o) {
	return json_object_iter_key(json_object_iter(o));
}

inline JsonObjIterator jsonObjIteratorNext(JsonObj o, JsonObjIterator i) {
	return json_object_iter_key(json_object_iter_next(o, json_object_key_to_iter(i)));
}

inline JsonObjIteratorKey jsonObjIteratorKey(JsonObjIterator i) {
	return i;
}

inline JsonObjIteratorVal iteratorValue(JsonObjIterator i) {
	return json_object_iter_value(json_object_key_to_iter(i));
}

inline bool iteratorAtEnd(JsonObjIterator i, JsonObj) {
	return !i;
}

#endif

class unknown;

class Model {
	friend class unknown;
	public:
		/**
		 * Reads fields of this Model from file of the given path.
		 */
		int readJsonFile(string path);

		/**
		 * Writes JSON representation of this Model to JSON file of the given path.
		 */
		void writeJsonFile(string path, cyborgbear::JsonSerializationSettings sttngs = Compact);

		/**
		 * Loads fields of this Model from the given JSON text.
		 */
		int fromJson(string json);

		/**
		 * Returns JSON representation of this Model.
		 */
		string toJson(cyborgbear::JsonSerializationSettings sttngs = Compact);

#ifdef CYBORGBEAR_BOOST_ENABLED
		/**
		 * Returns Boost serialization version of this object.
		 */
		virtual string toBoostBinary() = 0;

		/**
		 * Loads fields of this Model from the given Boost serialization text.
		 */
		virtual void fromBoostBinary(string dat) = 0;
#endif

#ifdef CYBORGBEAR_USING_QT
		cyborgbear::Error loadJsonObj(cyborgbear::JsonObjIteratorVal &obj) { return loadJsonObj(obj); };
#endif
	protected:
		virtual cyborgbear::JsonValOut buildJsonObj() = 0;
		virtual cyborgbear::Error loadJsonObj(cyborgbear::JsonVal obj) = 0;
};

class unknown: public Model {
	cyborgbear::string m_data;
	cyborgbear::Type m_type;

#ifdef CYBORGBEAR_BOOST_ENABLED
	friend class boost::serialization::access;

	template<class Archive>
	void serialize(Archive &ar, const unsigned int) {
		ar & m_type;
		ar & m_data;
	}
#endif

	public:
		unknown();
		unknown(Model *v);
		unknown(bool v);
		unknown(int v);
		unknown(double v);
		unknown(string v);
		virtual ~unknown();

		bool loaded();
		cyborgbear::Error loadJsonObj(cyborgbear::JsonVal obj);
		cyborgbear::JsonValOut buildJsonObj();

		bool toBool();
		int toInt();
		double toDouble();
		string toString();

		bool isBool();
		bool isInt();
		bool isDouble();
		bool isString();
		bool isObject();

		void set(Model* v);
		void set(bool v);
		void set(int v);
		void set(double v);
		void set(string v);

		bool operator==(const unknown&) const;
		bool operator!=(const unknown&) const;

#ifdef CYBORGBEAR_BOOST_ENABLED
		/**
		 * Returns Boost serialization version of this object.
		 */
		string toBoostBinary();

		/**
		 * Loads fields of this Model from the given Boost serialization text.
		 */
		void fromBoostBinary(string dat);
#endif
};

};

};




namespace models {

using cyborgbear::string;

class Model1: public cyborgbear::Model {

	public:

		Model1();

		cyborgbear::Error loadJsonObj(cyborgbear::JsonVal obj);

		cyborgbear::JsonValOut buildJsonObj();

		bool operator==(const Model1&) const;

		bool operator!=(const Model1&) const;
#ifdef CYBORGBEAR_BOOST_ENABLED

		virtual string toBoostBinary();

		virtual void fromBoostBinary(string dat);
#endif
		string Field1;
		cyborgbear::unknown Field2;
		int Field3[4];
		std::vector< std::vector< string > > Field4;
		std::map< string, string > Field5;
};

}


#ifdef CYBORGBEAR_BOOST_ENABLED
#include <boost/serialization/list.hpp>
#include <boost/serialization/string.hpp>
#include <boost/archive/binary_oarchive.hpp>
#include <boost/archive/binary_iarchive.hpp>

namespace boost {
namespace serialization {

template<class Archive>
void serialize(Archive &ar, models::Model1 &model, const unsigned int) {
	ar & model.Field1;
	ar & model.Field2;
	ar & model.Field3;
	ar & model.Field4;
	ar & model.Field5;
}

}
}
#endif
#endif