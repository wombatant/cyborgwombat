/*
   Copyright 2013 gtalent2@gmail.com

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
	"./parser"
	"strconv"
	"strings"
)

const (
	USING_JANSSON = iota
	USING_QT      = iota
)

type Cpp struct {
	hppPrefix   string
	hpp         string
	constructor string
	reader      string
	writer      string
	namespace   string
	lib         int
}

func NewCOut(namespace string, lib int) *Cpp {
	out := new(Cpp)
	out.namespace = namespace
	out.lib = lib
	out.hppPrefix = `#include <string>
#include <sstream>

` + out.buildModelmakerDefsHeader() + `


`
	return out
}

func (me *Cpp) typeMap(t string) string {
	switch t {
	case "float", "float32", "float64", "double":
		return "double"
	case "unknown":
		return "modelmaker::unknown"
	default:
		return t
	}
	return t
}

func (me *Cpp) buildTypeDec(t string, index []parser.VarType) string {
	array := ""
	out := ""
	for i := 0; i < len(index); i++ {
		if i != len(index)-1 && index[i].Type == "array" {
			array += "[" + index[i].Index + "]"
		} else if index[i].Type == "slice" {
			out += "std::vector< "
			array += " >"
		} else if index[i].Type == "map" {
			out += "std::map< " + index[i].Index + ", "
			array += " >"
		}
	}
	out += t + array
	return out
}

func (me *Cpp) buildVar(v, t string, index []parser.VarType) string {
	array := ""
	if len(index) > 0 && index[len(index)-1].Type == "array" {
		array = "[" + index[len(index)-1].Index + "]"
	}
	return me.buildTypeDec(t, index) + " " + v + array + ";"
}

func (me *Cpp) addVar(v string, index []parser.VarType) {
	jsonV := v
	if len(v) > 0 && v[0] < 91 {
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
}

func (me *Cpp) addClass(v string) {
	me.hpp += "\nnamespace " + me.namespace + " {\n"
	me.hpp += "\nusing modelmaker::string;\n"
	me.hpp += "\nclass " + v + ": public modelmaker::Model {\n"
	me.hpp += "\n\tpublic:\n"
	me.hpp += "\n\t\t" + v + "();\n"
	me.hpp += "\n\t\tbool loadJsonObj(modelmaker::JsonVal obj);\n"
	me.hpp += "\n\t\tmodelmaker::JsonValOut buildJsonObj();\n\n"
	me.constructor += v + "::" + v + "() {\n"
	me.reader += "bool " + v + "::loadJsonObj(modelmaker::JsonVal in) {\n"
	me.reader += "\tmodelmaker::JsonObjOut inObj = modelmaker::toObj(in);"
	me.writer += "modelmaker::JsonValOut " + v + `::buildJsonObj() {
	modelmaker::JsonObjOut obj = modelmaker::newJsonObj();`
}

func (me *Cpp) closeClass() {
	me.hpp += "};\n\n"
	me.hpp += "}\n\n"
	me.constructor += "}\n\n"
	me.reader += "\n\treturn true;\n}\n\n"
	me.writer += "\n\treturn obj;\n}\n\n"
}

func (me *Cpp) header(fileName string) string {
	n := strings.ToUpper(fileName)
	n = strings.Replace(n, ".", "_", -1)
	return `//Generated Code
#ifndef ` + n + `
#define ` + n + `
` + me.hppPrefix + me.hpp + `
#endif`
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
	return include + me.constructor + me.reader + me.writer[:len(me.writer)-1]
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
		code.Insert("modelmaker::JsonValOut obj0 = modelmaker::objRead(inObj, \"" + jsonV + "\");")
	}
	if len(index) > 0 {
		is := "i"
		for i := 0; i < depth; i++ {
			is += "i"
		}
		if index[0].Type == "array" || index[0].Type == "slice" {
			code.PushIfBlock("!modelmaker::isNull(obj" + strconv.Itoa(depth) + ") && modelmaker::isArray(obj" + strconv.Itoa(depth) + ")")
			code.Insert("modelmaker::JsonArrayOut array" + strconv.Itoa(depth) + " = modelmaker::toArray(obj" + strconv.Itoa(depth) + ");")
			code.Insert("unsigned int size = modelmaker::arraySize(array" + strconv.Itoa(depth) + ");")
			if index[0].Type == "slice" {
				code.Insert("this->" + v + sub + ".resize(size);")
			}
			code.PushForBlock("unsigned int " + is + " = 0; " + is + " < size; " + is + "++")
			code.Insert("modelmaker::JsonValOut obj" + strconv.Itoa(depth+1) + " = modelmaker::arrayRead(array" + strconv.Itoa(depth) + ", " + is + ");")
			me.buildReader(code, v, jsonV, t, sub+"["+is+"]", index[1:], depth+1)
			code.PopBlock()
			code.PopBlock()
		} else if index[0].Type == "map" {
			code.PushIfBlock("!modelmaker::isNull(obj" + strconv.Itoa(depth) + ") && modelmaker::isObj(obj" + strconv.Itoa(depth) + ")")
			code.Insert("modelmaker::JsonObjOut map" + strconv.Itoa(depth) + " = modelmaker::toObj(obj" + strconv.Itoa(depth) + ");")
			code.PushForBlock("modelmaker::JsonObjIterator it" + strconv.Itoa(depth+1) + " = modelmaker::iterator(map" + strconv.Itoa(depth) + "); !modelmaker::iteratorAtEnd(it" + strconv.Itoa(depth+1) + ", map" + strconv.Itoa(depth) + "); " + "it" + strconv.Itoa(depth+1) + " = modelmaker::iteratorNext(map" + strconv.Itoa(depth) + ",  it" + strconv.Itoa(depth+1) + ")")
			code.Insert(index[0].Index + " " + is + ";")
			code.Insert("modelmaker::JsonValOut obj" + strconv.Itoa(depth+1) + " = modelmaker::iteratorValue(it" + strconv.Itoa(depth+1) + ");")
			code.PushBlock()
			code.Insert("std::string key = modelmaker::toStdString(modelmaker::iteratorKey(it" + strconv.Itoa(depth+1) + "));")
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
		}
	} else {
		code.PushBlock()
		switch t { //type
		case "int":
			code.PushIfBlock("modelmaker::isInt(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = modelmaker::toInt(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "double":
			code.PushIfBlock("modelmaker::isDouble(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = modelmaker::toDouble(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "bool":
			code.PushIfBlock("modelmaker::isBool(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = modelmaker::toBool(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "string":
			code.PushIfBlock("modelmaker::isString(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = modelmaker::toString(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "modelmaker::unknown":
			code.Insert("this->" + v + sub + ".loadJsonObj(obj" + strconv.Itoa(depth) + ");")
		default:
			code.Insert("modelmaker::JsonValOut finalObj = modelmaker::toObj(obj" + strconv.Itoa(depth) + ");")
			code.PushIfBlock("modelmaker::isObj(finalObj)")
			code.Insert("this->" + v + sub + ".loadJsonObj(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		}
		code.PopBlock()
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

	if len(index) > depth {
		if index[depth].Type == "array" || index[depth].Type == "slice" {
			code.Insert("modelmaker::JsonArrayOut out" + strconv.Itoa(len(index[depth:])) + " = modelmaker::newJsonArray();")
			if index[depth].Type == "slice" {
				code.PushForBlock("unsigned int " + is + " = 0; " + is + " < this->" + v + sub + ".size(); " + is + "++")
			} else { // array
				code.PushForBlock("unsigned int " + is + " = 0; " + is + " < " + index[depth].Index + "; " + is + "++")
			}
			me.buildArrayWriter(code, t, v, sub+"["+is+"]", depth+1, index)
			code.Insert("modelmaker::arrayAdd(out" + strconv.Itoa(len(index[depth:])) + ", out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.Insert("modelmaker::decref(out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.PopBlock()
		} else if index[depth].Type == "map" {
			code.Insert("modelmaker::JsonObjOut out" + strconv.Itoa(len(index[depth:])) + " = modelmaker::newJsonObj();")
			code.PushForBlock(me.buildTypeDec(t, index[depth:]) + "::iterator " + ns + " = this->" + v + sub + ".begin(); " + ns + " != this->" + v + sub + ".end(); ++" + ns)
			switch index[depth].Index {
			case "bool":
				code.Insert("string key = " + ns + "->first ? \"true\" : \"false\";")
			case "string":
				code.Insert("std::stringstream s;")
				code.Insert("string key;")
				code.Insert("std::string tmp;")
				code.Insert("s << modelmaker::toStdString(modelmaker::toString(" + ns + "->first));")
				code.Insert("s >> tmp;")
				code.Insert("key = modelmaker::toString(tmp);")
			case "int", "double":
				code.Insert("std::stringstream s;")
				code.Insert("string key;")
				code.Insert("std::string tmp;")
				code.Insert("s << " + ns + "->first;")
				code.Insert("s >> tmp;")
				code.Insert("key = modelmaker::toString(tmp);")
			}
			me.buildArrayWriter(code, t, v, sub+"["+ns+"->first]", depth+1, index)
			code.Insert("modelmaker::objSet(out" + strconv.Itoa(len(index[depth:])) + ", key, out" + strconv.Itoa(len(index[depth:])-1) + ");")
			code.Insert("modelmaker::decref(out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.PopBlock()
		}
	} else {
		switch t {
		case "int":
			code.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + sub + ");")
		case "double":
			code.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + sub + ");")
		case "bool":
			code.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + sub + ");")
		case "string":
			code.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + sub + ");")
		default:
			code.Insert("modelmaker::JsonValOut obj0 = this->" + v + sub + ".buildJsonObj();")
			code.Insert("modelmaker::JsonValOut out0 = obj0;")
		}
	}
}

func (me *Cpp) buildWriter(v, jsonV, t string, index []parser.VarType) string {
	var out CppCode
	out.tabs = "\t"
	out.PushBlock()
	if len(index) > 0 {
		me.buildArrayWriter(&out, t, v, "", 0, index)
		out.Insert("modelmaker::objSet(obj, \"" + jsonV + "\", out" + strconv.Itoa(len(index)) + ");")
		out.Insert("modelmaker::decref(out" + strconv.Itoa(len(index)) + ");")
	} else {
		switch t {
		case "int":
			out.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + ");")
		case "double":
			out.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + ");")
		case "bool":
			out.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + ");")
		case "string":
			out.Insert("modelmaker::JsonValOut out0 = modelmaker::toJsonVal(this->" + v + ");")
		default:
			out.Insert("modelmaker::JsonValOut obj0 = this->" + v + ".buildJsonObj();")
			out.Insert("modelmaker::JsonValOut out0 = obj0;")
		}
		out.Insert("modelmaker::objSet(obj, \"" + jsonV + "\", out0);")
		out.Insert("modelmaker::decref(out0);")
	}
	out.PopBlock()
	return out.String()
}

func (me *Cpp) buildModelmakerDefsHeader() string {
	using := ""
	if me.lib == USING_QT {
		using = "USING_QT"
	} else {
		using = "USING_JANSSON"
	}
	out := `//Generated Code

#define ` + using + `

#include <string>

#include <vector>
#include <map>

#ifdef USING_QT
#include <QString>
#include <QJsonDocument>
#include <QJsonArray>
#include <QJsonObject>
#include <QJsonValue>
#include <QMap>
#include <QVector>
#else
#include <string>
#include <jansson.h>
#endif

namespace ` + me.namespace + ` {

namespace modelmaker {

#ifdef USING_QT
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
#endif

//string ops
std::string toStdString(string str);
const char* toCString(string str);


JsonObjOut read(const char *json);

int toInt(JsonVal);
double toDouble(JsonVal);
bool toBool(JsonVal);
string toString(JsonVal);
JsonArrayOut toArray(JsonVal);
JsonObjOut toObj(JsonVal);

JsonValOut toJsonVal(int);
JsonValOut toJsonVal(double);
JsonValOut toJsonVal(bool);
JsonValOut toJsonVal(string);
JsonValOut toJsonVal(JsonArray);
JsonValOut toJsonVal(JsonObj);


//value methods

bool isBool(JsonVal);
bool isInt(JsonVal);
bool isDouble(JsonVal);
bool isString(JsonVal);
bool isArray(JsonVal);
bool isObj(JsonVal);

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


JsonObjIterator iterator(JsonObj);
JsonObjIterator iteratorNext(JsonObj, JsonObjIterator);
JsonObjIteratorKey iteratorKey(JsonObjIterator);
JsonObjIteratorVal iteratorValue(JsonObjIterator);
bool iteratorAtEnd(JsonObjIterator, JsonObj);



inline string toString(string str) {
	return str;
}


#ifdef USING_QT

//string conversions
inline std::string toStdString(string str) {
	return str.toStdString();
}

inline const char* toCString(std::string str) {
	return str.c_str();
}

inline const char* toCString(string str) {
	return toStdString(str).c_str();
}

inline string toString(std::string str) {
	return QString::fromStdString(str);
}


inline JsonObjOut read(const char *json) {
	return QJsonDocument::fromJson(QByteArray(json)).object();
}


//from JsonObjIteratorVal
inline int toInt(JsonObjIteratorVal v) {
	return (int) v.toDouble();
}

inline double toDouble(JsonObjIteratorVal v) {
	return v.toDouble();
}

inline bool toBool(JsonObjIteratorVal v) {
	return v.toBool();
}

inline string toString(JsonObjIteratorVal v) {
	return v.toString();
}

inline JsonArrayOut toArray(JsonObjIteratorVal v) {
	return v.toArray();
}

inline JsonObjOut toObj(JsonObjIteratorVal v) {
	return v.toObject();
}

//from JsonVal
inline int toInt(JsonVal v) {
	return (int) v.toDouble();
}

inline double toDouble(JsonVal v) {
	return v.toDouble();
}

inline bool toBool(JsonVal v) {
	return v.toBool();
}

inline string toString(JsonVal v) {
	return v.toString();
}

inline JsonArrayOut toArray(JsonVal v) {
	return v.toArray();
}

inline JsonObjOut toObj(JsonVal v) {
	return v.toObject();
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


inline bool isNull(JsonObjIteratorVal v) {
	return v.isNull();
}

inline bool isBool(JsonObjIteratorVal v) {
	return v.isBool();
}

inline bool isInt(JsonObjIteratorVal v) {
	return v.isDouble();
}

inline bool isDouble(JsonObjIteratorVal v) {
	return v.isDouble();
}

inline bool isString(JsonObjIteratorVal v) {
	return v.isString();
}

inline bool isArray(JsonObjIteratorVal v) {
	return v.isArray();
}

inline bool isObj(JsonObjIteratorVal v) {
	return v.isObject();
}

inline bool isBool(JsonVal v) {
	return v.isBool();
}

inline bool isInt(JsonVal v) {
	return v.isDouble();
}

inline bool isDouble(JsonVal v) {
	return v.isDouble();
}

inline bool isString(JsonVal v) {
	return v.isString();
}

inline bool isArray(JsonVal v) {
	return v.isArray();
}

inline bool isObj(JsonVal v) {
	return v.isObject();
}

inline bool isNull(JsonVal v) {
	return v.isNull();
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
	JsonValOut tmp = modelmaker::toJsonVal(v);
	a.append(tmp);
}

inline void arrayAdd(JsonArray a, JsonObj v) {
	JsonValOut tmp = modelmaker::toJsonVal(v);
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
	JsonValOut tmp = modelmaker::toJsonVal(v);
	o.insert(k, tmp);
}

inline void objSet(JsonObj o, string k, JsonObj v) {
	JsonValOut tmp = modelmaker::toJsonVal(v);
	o.insert(k, tmp);
}

inline void objSet(JsonObj o, string k, JsonVal v) {
	o.insert(k, v);
}


inline JsonValOut objRead(JsonObj o, string k) {
	return o[k];
}

inline JsonObjIterator iterator(JsonObj o) {
	return o.begin();
}

inline JsonObjIterator iteratorNext(JsonObj, JsonObjIterator i) {
	return i + 1;
}

inline JsonObjIteratorKey iteratorKey(JsonObjIterator i) {
	return i.key();
}

inline bool iteratorAtEnd(JsonObjIterator i, JsonObj o) {
	return i == o.end();
}

inline JsonObjIteratorVal iteratorValue(JsonObjIterator i) {
	return i.value();
}

inline string write(JsonObj obj) {
	QJsonDocument doc(obj);
	return doc.toJson();
}

#else

inline std::string toStdString(string str) {
	return str;
}

inline const char* toCString(string str) {
	return str.c_str();
}


inline JsonObjOut read(const char *json) {
	return json_loads(json, 0, NULL);
}

inline string write(JsonObj obj) {
	char *tmp = json_dumps(obj, JSON_COMPACT);
	if (!tmp)
		return "{}";
	string out = tmp;
	free(tmp);
	modelmaker::decref(obj);
	return out;
}

//value methods

inline int toInt(JsonVal v) {
	return (int) json_integer_value(v);
}

inline double toDouble(JsonVal v) {
	return (double) json_real_value(v);
}

inline bool toBool(JsonVal v) {
	return json_is_true(v);
}

inline string toString(JsonVal v) {
	return json_string_value(v);
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


inline bool isNull(JsonVal v) {
	return !v;
}

inline bool isBool(JsonVal v) {
	return json_is_boolean(v);
}

inline bool isInt(JsonVal v) {
	return json_is_integer(v);
}

inline bool isDouble(JsonVal v) {
	return json_is_real(v);
}

inline bool isString(JsonVal v) {
	return json_is_string(v);
}

inline bool isArray(JsonVal v) {
	return json_is_array(v);
}

inline bool isObj(JsonVal v) {
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


inline JsonObjIterator iterator(JsonObj o) {
	return json_object_iter_key(json_object_iter(o));
}

inline JsonObjIterator iteratorNext(JsonObj o, JsonObjIterator i) {
	return json_object_iter_key(json_object_iter_next(o, json_object_key_to_iter(i)));
}

inline JsonObjIteratorKey iteratorKey(JsonObjIterator i) {
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
		bool loadFile(string path);
		void writeFile(string path);
		void load(string json);
		string write();
#ifdef USING_QT
		bool loadJsonObj(modelmaker::JsonObjIteratorVal &obj) { return loadJsonObj(obj); };
#endif
	protected:
		virtual modelmaker::JsonValOut buildJsonObj() = 0;
		virtual bool loadJsonObj(modelmaker::JsonVal obj) = 0;
};

class unknown: public Model {
	private:
		modelmaker::JsonValOut m_obj;
	public:
		unknown();
		unknown(Model *v);
		unknown(bool v);
		unknown(int v);
		unknown(double v);
		unknown(string v);
		virtual ~unknown();

		bool loaded();
		bool loadJsonObj(modelmaker::JsonVal obj);
		modelmaker::JsonValOut buildJsonObj();

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
};

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
using namespace ` + me.namespace + `::modelmaker;

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
`
	return out
}
