/*
   Copyright 2013 - 2014 gtalent2@gmail.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"github.com/gtalent/cyborgbear/parser"
	"io/ioutil"
	"strconv"
	"strings"
)

type Cpp struct {
	hppPrefix    string
	hpp          string
	constructor  string
	reader       string
	writer       string
	equals       string
	notEquals    string
	namespace    string
	lowerCase    bool
	lib          int
}

func NewCOut(namespace string, lib int, lowerCase bool) *Cpp {
	out := new(Cpp)
	out.lowerCase = lowerCase
	out.namespace = namespace
	out.lib = lib
	out.hppPrefix = `#include <string>
#include <sstream>

` + out.buildModelmakerDefsHeader() + `


`
	return out
}

func (me *Cpp) write(outFile string) string {
	return me.header("") + "\n" + me.body("")
}

func (me *Cpp) writeFile(outFile string) error {
	var err error
	err = ioutil.WriteFile(outFile+".hpp", []byte(me.header(outFile+".hpp")), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outFile+".cpp", []byte(me.body(outFile+".hpp")), 0644)
	return err
}

func (me *Cpp) typeMap(t string) string {
	switch t {
	case "float", "float32", "float64", "double":
		return "double"
	case "int", "bool", "string":
		return t
	case "unknown":
		return "cyborgbear::unknown"
	default:
		return me.namespace + "::" + t
	}
	return t
}

func (me *Cpp) buildTypeDec(t string, index []parser.VarType) string {
	array := ""
	out := ""
	for i := 0; i < len(index); i++ {
		if index[i].Type == "array" {
			out += "cyborgbear::Array<"
			array += ", " + index[i].Index + ">"
		} else if index[i].Type == "slice" {
			vector := ""
			if me.lib == USING_QT {
				vector = "QVector"
			} else {
				vector = "std::vector"
			}
			out += vector + "< "
			array += " >"
		} else if index[i].Type == "map" {
			cmap := ""
			if me.lib == USING_QT {
				cmap = "QMap"
			} else {
				cmap = "std::map"
			}
			out += cmap + "< " + index[i].Index + ", "
			array += " >"
		}
	}
	out += t + array
	return out
}

func (me *Cpp) buildBoostSerialize(name string) string {
	return "\tar & model." + name + ";\n"
}

func (me *Cpp) buildVar(v, t string, index []parser.VarType) string {
	return me.buildTypeDec(t, index) + " " + v + ";"
}

func (me *Cpp) addVar(v string, index []parser.VarType) {
	jsonV := v
	if me.lowerCase && len(v) > 0 && v[0] < 91 {
		v = string(v[0]+32) + v[1:]
	}
	t := me.typeMap(index[len(index)-1].Type)
	index = index[:len(index)-1]
	me.hpp += "\t\t" + me.buildVar(v, t, index) + "\n"
	var reader CppCode
	reader.tabs += "\t"
	me.constructor += me.buildConstructor(v, t, index)
	me.reader += me.buildReader(&reader, v, jsonV, t, "", index, 0)
	me.writer += me.buildWriter(v, jsonV, t, index)
	me.equals += me.buildEquals(v)
	me.notEquals += me.buildNotEquals(v)
}

func (me *Cpp) addClass(v string) {
	me.hpp += "\nnamespace " + me.namespace + " {\n"
	me.hpp += "\nusing cyborgbear::string;\n"
	me.hpp += "\nclass " + v + ": public cyborgbear::Model {\n"
	me.hpp += "\n\tpublic:\n"
	me.hpp += "\n\t\t" + v + "();\n"
	me.hpp += "\n\t\tcyborgbear::Error loadJsonObj(cyborgbear::JsonVal obj);\n"
	me.hpp += "\n\t\tcyborgbear::JsonValOut buildJsonObj();\n"
	me.hpp += "\n\t\tbool operator==(const " + v + "&) const;\n"
	me.hpp += "\n\t\tbool operator!=(const " + v + "&) const;\n"

	me.constructor += v + "::" + v + "() {\n"
	me.reader += "cyborgbear::Error " + v + `::loadJsonObj(cyborgbear::JsonVal in) {
	cyborgbear::Error retval = cyborgbear::Error_Ok;
	cyborgbear::JsonObjOut inObj = cyborgbear::toObj(in);
`
	me.writer += "cyborgbear::JsonValOut " + v + `::buildJsonObj() {
	cyborgbear::JsonObjOut obj = cyborgbear::newJsonObj();`
	me.equals += "bool " + v + "::operator==(const " + v + " &o) const {\n"
	me.notEquals += "bool " + v + "::operator!=(const " + v + " &o) const {\n"
}

func (me *Cpp) closeClass(v string) {
	me.hpp += "};\n\n"
	me.hpp += "}\n\n"
	me.constructor += "}\n\n"
	me.reader += "\n\treturn retval;\n}\n\n"
	me.writer += "\n\treturn obj;\n}\n\n"
	me.equals += "\n\treturn true;\n}\n\n"
	me.notEquals += "\n\treturn false;\n}\n\n"
}

func (me *Cpp) header(fileName string) string {
	n := strings.ToUpper(fileName)
	n = strings.Replace(n, ".", "_", -1)
	out := `//Generated Code

#ifndef ` + n + `
#define ` + n + `
` + me.hppPrefix + me.hpp + `
#endif`
	return out
}

func (me *Cpp) body(headername string) string {
	include := ""
	if headername != "" {
		include += `//Generated Code
` + me.buildModelmakerDefsBody(headername) + `

#include "string.h"
#include "` + headername + `"

using namespace ` + me.namespace + `;
using std::stringstream;

`
	}
	writer := me.writer
	if len(me.writer) > 1 {
		writer = me.writer[:len(me.writer)-1]
	}
	return include + me.constructor + me.reader + writer + me.equals + me.notEquals
}

func (me *Cpp) endsWithClose() bool {
	return me.hpp[len(me.hpp)-3:] == "};\n"
}

func (me *Cpp) buildConstructor(v, t string, index []parser.VarType) string {
	if len(index) < 1 {
		switch t {
		case "bool", "int", "double":
			return "\tthis->" + v + " = 0;\n"
		case "string":
			return "\tthis->" + v + " = \"\";\n"
		}
	} else if index[0].Type == "array" {
		switch t {
		case "bool", "int", "double":
			return "\tfor (int i = 0; i < " + index[0].Index + "; this->" + v + "[i++] = 0);\n"
		}
	}
	return ""
}

func (me *Cpp) buildReader(code *CppCode, v, jsonV, t, sub string, index []parser.VarType, depth int) string {
	if depth == 0 {
		code.PushBlock()
		code.Insert("cyborgbear::JsonValOut obj0 = cyborgbear::objRead(inObj, \"" + jsonV + "\");")
	}
	if len(index) > 0 {
		is := "i"
		for i := 0; i < depth; i++ {
			is += "i"
		}
		if index[0].Type == "slice" || index[0].Type == "array" {
			code.Insert("retval |= cyborgbear::readVal(obj" + strconv.Itoa(depth) + ", this->" + v + sub + ");")
		} else if index[0].Type == "map" {
			code.PushIfBlock("!cyborgbear::isNull(obj" + strconv.Itoa(depth) + ")")
			code.PushIfBlock("cyborgbear::isObj(obj" + strconv.Itoa(depth) + ")")
			code.Insert("cyborgbear::JsonObjOut map" + strconv.Itoa(depth) + " = cyborgbear::toObj(obj" + strconv.Itoa(depth) + ");")
			code.PushForBlock("cyborgbear::JsonObjIterator it" + strconv.Itoa(depth+1) + " = cyborgbear::jsonObjIterator(map" + strconv.Itoa(depth) + "); !cyborgbear::iteratorAtEnd(it" + strconv.Itoa(depth+1) + ", map" + strconv.Itoa(depth) + "); " + "it" + strconv.Itoa(depth+1) + " = cyborgbear::jsonObjIteratorNext(map" + strconv.Itoa(depth) + ",  it" + strconv.Itoa(depth+1) + ")")
			code.Insert(index[0].Index + " " + is + ";")
			code.Insert("cyborgbear::JsonValOut obj" + strconv.Itoa(depth+1) + " = cyborgbear::iteratorValue(it" + strconv.Itoa(depth+1) + ");")
			code.PushBlock()
			code.Insert("std::string key = cyborgbear::toStdString(cyborgbear::jsonObjIteratorKey(it" + strconv.Itoa(depth+1) + "));")
			switch index[0].Index {
			case "bool":
				code.Insert(is + " = key == \"true\";")
			case "string":
				code.Insert("std::string o;")
				code.Insert("std::stringstream s;")
				code.Insert("s << key;")
				code.Insert("s >> o;")
				code.Insert(is + " = o.c_str();")
			case "double", "int":
				code.Insert("std::stringstream s;")
				code.Insert("s << key;")
				code.Insert("s >> " + is + ";")
			}
			code.PopBlock()

			me.buildReader(code, v, jsonV, t, sub+"["+is+"]", index[1:], depth+1)
			code.PopBlock()
			code.PopBlock()
			code.PopBlock()
		}
	} else {
		code.Insert("retval |= cyborgbear::readVal(obj" + strconv.Itoa(depth) + ", this->" + v + sub + ");")
	}

	if depth == 0 {
		code.PopBlock()
	}

	return code.String()
}

func (me *Cpp) buildArrayWriter(code *CppCode, t, v, sub string, depth int, index []parser.VarType) {
	is := "i"
	ns := "n"
	for i := 0; i < depth; i++ {
		is += "i"
		ns += "n"
	}

	itKey := "->first"
	if me.lib == USING_QT {
		itKey = ".key()"
	}

	if len(index) > depth {
		if index[depth].Type == "array" || index[depth].Type == "slice" {
			code.Insert("cyborgbear::JsonArrayOut out" + strconv.Itoa(len(index[depth:])) + " = cyborgbear::newJsonArray();")
			if index[depth].Type == "slice" {
				code.PushForBlock("cyborgbear::VectorIterator " + is + " = 0; " + is + " < this->" + v + sub + ".size(); " + is + "++")
			} else { // array
				code.PushForBlock("cyborgbear::VectorIterator " + is + " = 0; " + is + " < " + index[depth].Index + "; " + is + "++")
			}
			me.buildArrayWriter(code, t, v, sub+"["+is+"]", depth+1, index)
			code.Insert("cyborgbear::arrayAdd(out" + strconv.Itoa(len(index[depth:])) + ", out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.Insert("cyborgbear::decref(out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.PopBlock()
		} else if index[depth].Type == "map" {
			code.Insert("cyborgbear::JsonObjOut out" + strconv.Itoa(len(index[depth:])) + " = cyborgbear::newJsonObj();")
			code.PushForBlock(me.buildTypeDec(t, index[depth:]) + "::iterator " + ns + " = this->" + v + sub + ".begin(); " + ns + " != this->" + v + sub + ".end(); ++" + ns)
			switch index[depth].Index {
			case "bool":
				code.Insert("string key = " + ns + itKey + " ? \"true\" : \"false\";")
			case "string":
				code.Insert("std::stringstream s;")
				code.Insert("string key;")
				code.Insert("std::string tmp;")
				code.Insert("s << cyborgbear::toStdString(cyborgbear::toString(" + ns + itKey + "));")
				code.Insert("s >> tmp;")
				code.Insert("key = cyborgbear::toString(tmp);")
			case "int", "double":
				code.Insert("std::stringstream s;")
				code.Insert("string key;")
				code.Insert("std::string tmp;")
				code.Insert("s << " + ns + itKey + ";")
				code.Insert("s >> tmp;")
				code.Insert("key = cyborgbear::toString(tmp);")
			}
			me.buildArrayWriter(code, t, v, sub+"["+ns+itKey+"]", depth+1, index)
			code.Insert("cyborgbear::objSet(out" + strconv.Itoa(len(index[depth:])) + ", key, out" + strconv.Itoa(len(index[depth:])-1) + ");")
			code.Insert("cyborgbear::decref(out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.PopBlock()
		}
	} else {
		switch t {
		case "int":
			code.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + sub + ");")
		case "double":
			code.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + sub + ");")
		case "bool":
			code.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + sub + ");")
		case "string":
			code.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + sub + ");")
		default:
			code.Insert("cyborgbear::JsonValOut obj0 = this->" + v + sub + ".buildJsonObj();")
			code.Insert("cyborgbear::JsonValOut out0 = obj0;")
		}
	}
}

func (me *Cpp) buildWriter(v, jsonV, t string, index []parser.VarType) string {
	var out CppCode
	out.tabs = "\t"
	out.PushBlock()
	if len(index) > 0 {
		me.buildArrayWriter(&out, t, v, "", 0, index)
		out.Insert("cyborgbear::objSet(obj, \"" + jsonV + "\", out" + strconv.Itoa(len(index)) + ");")
		out.Insert("cyborgbear::decref(out" + strconv.Itoa(len(index)) + ");")
	} else {
		switch t {
		case "int":
			out.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + ");")
		case "double":
			out.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + ");")
		case "bool":
			out.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + ");")
		case "string":
			out.Insert("cyborgbear::JsonValOut out0 = cyborgbear::toJsonVal(this->" + v + ");")
		default:
			out.Insert("cyborgbear::JsonValOut obj0 = this->" + v + ".buildJsonObj();")
			out.Insert("cyborgbear::JsonValOut out0 = obj0;")
		}
		out.Insert("cyborgbear::objSet(obj, \"" + jsonV + "\", out0);")
		out.Insert("cyborgbear::decref(out0);")
	}
	out.PopBlock()
	return out.String()
}

func (me *Cpp) buildEquals(v string) string {
	return "\tif (" + v + " != o." + v + ") return false;\n"
}

func (me *Cpp) buildNotEquals(v string) string {
	return "\tif (" + v + " != o." + v + ") return true;\n"
}

func (me *Cpp) buildModelmakerDefsHeader() string {
	using := ""
	if me.lib == USING_QT {
		using += "#define CYBORGBEAR_USING_QT\n"
	} else {
		using += "#define CYBORGBEAR_USING_JANSSON\n"
	}
	out := using + `

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

namespace ` + me.namespace + ` {

namespace cyborgbear {

typedef unsigned long int Error;
const Error Error_Ok = 0;
const Error Error_TypeMismatch = 1;
const Error Error_MissingField = 2;
const Error Error_CouldNotAccessFile = 4;
const Error Error_GenericParsingError = 8;
const Error Error_ArrayOverflow = 16;

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

template<typename T, long long len>
class Array {
	protected:
		T data[len];

	public:
		T &operator[](long long idx) {
			return data[idx];
		}

		long long size() const {
			return len;
		}

		bool operator==(const Array<T, len> &other) const {
			for (int i = 0; i < size(); i++) {
				if (data[i] != other.data[i]) {
					return false;
				}
			}
			return true;
		}

		bool operator!=(const Array<T, len> &other) const {
			return !(*this == other);
		}

		template<class Archive>
		void serialize(Archive &ar, const unsigned int) {
		}
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

class unknown;
cyborgbear::Error readVal(JsonObjOut v, class Model &val);

class Model {
	friend class unknown;
	friend cyborgbear::Error readVal(JsonObjOut v, class Model &val);
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
};

/**
 * Version of cyborgbear.
 */
extern string version;

//string ops
std::string toStdString(string str);


JsonObjOut read(string json);


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


JsonArrayOut toArray(JsonVal);
JsonObjOut toObj(JsonVal);

JsonValOut toJsonVal(int);
JsonValOut toJsonVal(double);
JsonValOut toJsonVal(bool);
JsonValOut toJsonVal(string);
JsonValOut toJsonVal(JsonArray);
JsonValOut toJsonVal(JsonObj);


cyborgbear::Error readVal(JsonVal, int&);
cyborgbear::Error readVal(JsonVal, double&);
cyborgbear::Error readVal(JsonVal, bool&);
cyborgbear::Error readVal(JsonVal, string&);

template<typename T, long long len>
inline cyborgbear::Error readVal(JsonValOut v, Array<T, len> &val) {
	cyborgbear::Error retval = cyborgbear::Error_Ok;
	if (!cyborgbear::isNull(v)) {
		if (cyborgbear::isArray(v)) {
			cyborgbear::JsonArrayOut array = cyborgbear::toArray(v);
			unsigned int size = cyborgbear::arraySize(array);
			if (size > val.size()) {
				size = val.size();
				retval |= cyborgbear::Error_ArrayOverflow;
			}
			for (unsigned int i = 0; i < size; i++) {
				retval |= cyborgbear::readVal(cyborgbear::arrayRead(array, i), val[i]);
			}
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	} else {
		retval |= cyborgbear::Error_MissingField;
	}
	return retval;
}

template<typename T>
#ifdef CYBORGBEAR_USING_QT
inline cyborgbear::Error readVal(JsonValOut v, QVector<T> &val) {
#else
inline cyborgbear::Error readVal(JsonValOut v, std::vector<T> &val) {
#endif
	cyborgbear::Error retval = cyborgbear::Error_Ok;
	if (!cyborgbear::isNull(v)) {
		if (cyborgbear::isArray(v)) {
			cyborgbear::JsonArrayOut array = cyborgbear::toArray(v);
			unsigned int size = cyborgbear::arraySize(array);
			val.resize(size);
			for (unsigned int i = 0; i < size; i++) {
				retval |= cyborgbear::readVal(cyborgbear::arrayRead(array, i), val[i]);
			}
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	} else {
		retval |= cyborgbear::Error_MissingField;
	}
	return retval;
}

inline cyborgbear::Error readVal(JsonObjOut v, class Model &val) {
	cyborgbear::Error retval = cyborgbear::Error_Ok;
	if (cyborgbear::isObj(v)) {
		val.loadJsonObj(v);
	} else {
		if (cyborgbear::isNull(v)) {
			retval |= cyborgbear::Error_MissingField;
		} else {
			retval |= cyborgbear::Error_TypeMismatch;
		}
	}
	return retval;
}

inline cyborgbear::Error readVal(JsonObjOut v, class unknown &val) {
	cyborgbear::Error retval = cyborgbear::Error_Ok;
	if (!cyborgbear::isNull(v)) {
		retval |= val.loadJsonObj(v);
	} else {
		retval |= cyborgbear::Error_MissingField;
	}
	return retval;
}

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

};

};
`
	return out
}

func (me *Cpp) buildModelmakerDefsBody(headername string) string {
	out := `
#include <fstream>
#include "` + headername + `"

using namespace ` + me.namespace + `;
using namespace ` + me.namespace + `::cyborgbear;

string ` + me.namespace + `::cyborgbear::version = "` + cyborgbear_version + `";

int Model::readJsonFile(string path) {
	try {
		std::ifstream in;
		in.open(cyborgbear::toStdString(path).c_str());
		if (in.is_open()) {
			std::string json((std::istreambuf_iterator<char>(in)), std::istreambuf_iterator<char>());
			in.close();
			return fromJson(cyborgbear::toString(json));
		}
	} catch (...) {
	}
	return cyborgbear::Error_CouldNotAccessFile;
}

void Model::writeJsonFile(string path, cyborgbear::JsonSerializationSettings sttngs) {
	std::ofstream out;
	out.open(cyborgbear::toStdString(path).c_str());
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
	cyborgbear::JsonValOut obj = buildJsonObj();
	bool out;
	cyborgbear::readVal(obj, out);
	return out;
}

int unknown::toInt() {
	cyborgbear::JsonValOut obj = buildJsonObj();
	int out;
	cyborgbear::readVal(obj, out);
	return out;
}

double unknown::toDouble() {
	cyborgbear::JsonValOut obj = buildJsonObj();
	double out;
	cyborgbear::readVal(obj, out);
	return out;
}

string unknown::toString() {
	cyborgbear::JsonValOut obj = buildJsonObj();
	string out;
	cyborgbear::readVal(obj, out);
	return out;
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

bool unknown::operator==(const unknown &o) const {
	return m_type == o.m_type && m_data == o.m_data;
}

bool unknown::operator!=(const unknown &o) const {
	return m_type != o.m_type || m_data != o.m_data;
}
`
	return out
}
