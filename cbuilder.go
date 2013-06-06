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
	"strconv"
	"strings"
)

type Out struct {
	hppPrefix   string
	hpp         string
	constructor string
	reader      string
	writer      string
	namespace   string
}

func NewCOut(namespace string) Out {
	var out Out
	out.namespace = namespace
	out.hppPrefix = `#include <string>
#include <sstream>
#include <vector>
#include <map>
#include <json/json.h>
#include "modelmakerdefs.hpp"


using std::string;
using std::vector;
using std::map;
`
	return out
}

func (me *Out) typeMap(t string) string {
	switch t {
	case "float", "float32", "float64", "double":
		return "double"
	case "unknown":
		return "modelmaker::unknown"
	default:
		return t
	}
	return ""
}

func (me *Out) buildTypeDec(t string, index []string) string {
	array := ""
	out := ""
	for i := 0; i < len(index); i++ {
		if index[i] == "array" {
			out += "vector<"
			array += " >"
		} else if index[i][:3] == "map" {
			out += "map<" + index[i][4:] + ", "
			array += " >"
		}
	}
	out += t + array
	return out
}

func (me *Out) buildVar(v, t string, index []string) string {
	return me.buildTypeDec(t, index) + " " + v + ";"
}

func (me *Out) addVar(v string, index []string) {
	jsonV := v
	if len(v) > 0 && v[0] < 91 {
		v = string(v[0]+32) + v[1:]
	}
	t := me.typeMap(index[len(index)-1])
	index = index[:len(index)-1]
	me.hpp += "\t\t" + me.buildVar(v, t, index) + "\n"
	var reader CppCode
	reader.tabs += "\t"
	me.constructor += me.buildConstructor(v, t, len(index))
	me.reader += me.buildReader(&reader, v, jsonV, t, "", index, 0)
	me.writer += me.buildWriter(v, jsonV, t, index)
}

func (me *Out) addClass(v string) {
	me.hpp += "\nnamespace " + me.namespace + " {\n"
	me.hpp += "\nclass " + v + ": public modelmaker::Model {\n"
	me.hpp += "\n\tpublic:\n"
	me.hpp += "\n\t\t" + v + "();\n"
	me.hpp += "\n\t\tvoid load(string text);\n"
	me.hpp += "\n\t\tbool load(json_object *obj);\n"
	me.hpp += "\n\t\tjson_object* buildJsonObj();\n\n"
	me.reader += "void " + v + `::load(string json) {
	json_object *obj = json_tokener_parse(json.c_str());
	load(obj);
	json_object_put(obj);
}

`
	me.constructor += v + "::" + v + "() {\n"
	me.reader += "bool " + v + "::load(json_object *in) {"
	me.writer += "json_object* " + v + `::buildJsonObj() {
	json_object *obj = json_object_new_object();`
}

func (me *Out) closeClass() {
	me.hpp += "};\n\n"
	me.hpp += "}\n\n"
	me.constructor += "}\n\n"
	me.reader += "\n\treturn true;\n}\n\n"
	me.writer += "\n\treturn obj;\n}\n\n"
}

func (me *Out) header(fileName string) string {
	n := strings.ToUpper(fileName)
	n = strings.Replace(n, ".", "_", -1)
	return `//Generated Code
#ifndef ` + n + `
#define ` + n + `
` + me.hppPrefix + me.hpp + `
#endif`
}

func (me *Out) body(headername string) string {
	include := ""
	if headername != "" {
		include += `//Generated Code

#include "` + headername + `"

using namespace ` + me.namespace + `;
using std::stringstream;

`
	}
	return include + me.constructor + me.reader + me.writer[:len(me.writer)-1]
}

func (me *Out) endsWithClose() bool {
	return me.hpp[len(me.hpp)-3:] == "};\n"
}

func (me *Out) buildConstructor(v, t string, index int) string {
	if index < 1 {
		switch t {
		case "bool", "int", "double":
			return "\tthis->" + v + " = 0;\n"
		case "string":
			return "\tthis->" + v + " = \"\";\n"
		}
	}
	return ""
}

func (me *Out) buildReader(code *CppCode, v, jsonV, t, sub string, index []string, depth int) string {
	if depth == 0 {
		code.PushBlock()
		code.Insert("json_object *obj" + strconv.Itoa(depth) + " = json_object_object_get(in, \"" + jsonV + "\");")
	}
	if len(index) > 0 {
		is := "i"
		for i := 0; i < depth; i++ {
			is += "i"
		}
		if depth != 0 {
			//code.Insert("json_object *obj" + strconv.Itoa(depth) + " = json_object_array_get_idx(obj" + strconv.Itoa(depth-1) + ", i);")
		}
		if index[0] == "array" {
			code.PushIfBlock("obj" + strconv.Itoa(depth) + " != NULL && json_object_get_type(obj" + strconv.Itoa(depth) + ") == json_type_array")
			code.Insert("unsigned int size = json_object_array_length(obj" + strconv.Itoa(depth) + ");")
			code.Insert("this->" + v + sub + ".resize(size);")
			code.PushForBlock("unsigned int " + is + " = 0; " + is + " < size; " + is + "++")
			code.Insert("json_object *obj" + strconv.Itoa(depth+1) + " = json_object_array_get_idx(obj" + strconv.Itoa(depth) + ", " + is + ");")
			me.buildReader(code, v, jsonV, t, sub+"["+is+"]", index[1:], depth+1)
			code.PopBlock()
			code.PopBlock()
		} else if index[0][:3] == "map" {
			code.PushIfBlock("obj" + strconv.Itoa(depth) + " != NULL && json_object_get_type(obj" + strconv.Itoa(depth) + ") == json_type_object")
			code.PushPrefixBlock("json_object_object_foreach(obj" + strconv.Itoa(depth) + ", key, obj" + strconv.Itoa(depth+1) + ")")
			code.Insert(index[0][4:] + " " + is + ";")
			code.PushBlock()
			switch index[0][4:] {
			case "bool":
				code.Insert(is + " = key == \"true\";")
			case "double", "int", "string":
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
		i := 0
		if depth == 1 {
			code.PushBlock()
			code.Insert("json_object *obj" + strconv.Itoa(len(index)) + " = json_object_object_get(in, \"" + jsonV + "\");")
			code.PushIfBlock("obj" + strconv.Itoa(len(index)) + " != NULL")
		} else {
			i += 2
		}
		switch t { //type
		case "int", "double":
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(depth) + ") == " + "json_type_" + t)
			code.Insert("this->" + v + sub + " = json_object_get_" + t + "(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "bool":
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(depth) + ") == " + "json_type_boolean")
			code.Insert("this->" + v + sub + " = json_object_get_boolean(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "string":
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(depth) + ") == " + "json_type_string")
			code.Insert("this->" + v + sub + " = json_object_get_string(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "modelmaker::unknown":
			code.Insert("this->" + v + sub + ".load(obj" + strconv.Itoa(depth) + ");")
		default:
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(depth) + ") == " + "json_type_object")
			code.Insert("this->" + v + sub + ".load(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		}
		if depth == 1 {
			code.PopBlock()
			code.PopBlock()
		}
	}

	if depth == 0 {
		code.PopBlock()
	}

	return code.String()
}

func (me *Out) buildArrayWriter(code *CppCode, t, v, sub string, depth int, index []string) {
	is := "i"
	ns := "n"
	for i := 0; i < depth; i++ {
		is += "i"
		ns += "n"
	}

	if len(index) > depth {
		if index[depth] == "array" {
			code.Insert("json_object *out" + strconv.Itoa(len(index[depth:])) + " = json_object_new_array();")
			code.PushForBlock("unsigned int " + is + " = 0; " + is + " < this->" + v + sub + ".size(); " + is + "++")
			me.buildArrayWriter(code, t, v, sub+"["+is+"]", depth+1, index)
			code.Insert("json_object_array_add(out" + strconv.Itoa(len(index[depth:])) + ", out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.PopBlock()
		} else if index[depth][:3] == "map" {
			code.Insert("json_object *out" + strconv.Itoa(len(index[depth:])) + " = json_object_new_object();")
			code.PushForBlock(me.buildTypeDec(t, index[depth:]) + "::iterator " + ns + " = this->" + v + sub + ".begin(); " + ns + " != this->" + v + sub + ".end(); " + ns + "++")
			switch index[depth][4:] {
			case "bool":
				code.Insert("string key = " + ns + "->first ? \"true\" : \"false\";")
			case "string", "int", "double":
				code.Insert("std::stringstream s;")
				code.Insert("string key;")
				code.Insert("s << " + ns + "->first;")
				code.Insert("s >> key;")
			}
			me.buildArrayWriter(code, t, v, sub+"["+ns+"->first]", depth+1, index)
			code.Insert("json_object_object_add(out" + strconv.Itoa(len(index[depth:])) + ", key.c_str(), out" + strconv.Itoa(len(index[depth:])-1) + ");")
			code.PopBlock()
		}
	} else {
		switch t {
		case "int", "double":
			code.Insert("json_object *out0 = json_object_new_" + t + "(this->" + v + sub + ");")
		case "bool":
			code.Insert("json_object *out0 = json_object_new_boolean(this->" + v + sub + ");")
		case "string":
			code.Insert("json_object *out0 = json_object_new_string(this->" + v + sub + ".c_str());")
		default:
			code.Insert("json_object *out0 = this->" + v + sub + ".buildJsonObj();")
		}
	}
}

func (me *Out) buildWriter(v, jsonV, t string, index []string) string {
	var out CppCode
	out.tabs = "\t"
	out.PushBlock()
	if len(index) > 0 {
		me.buildArrayWriter(&out, t, v, "", 0, index)
		out.Insert("json_object_object_add(obj, \"" + jsonV + "\", out" + strconv.Itoa(len(index)) + ");")
	} else {
		switch t {
		case "int", "double":
			out.Insert("json_object *out0 = json_object_new_" + t + "(this->" + v + ");")
		case "bool":
			out.Insert("json_object *out0 = json_object_new_boolean(this->" + v + ");")
		case "string":
			out.Insert("json_object *out0 = json_object_new_string(this->" + v + ".c_str());")
		default:
			out.Insert("json_object *out0 = this->" + v + ".buildJsonObj();")
		}
		out.Insert("json_object_object_add(obj, \"" + jsonV + "\", out0);")
	}
	out.PopBlock()
	return out.String()
}

func (me *Out) buildModelmakerDefsHeader() string {
	out := `//Generated Code

#ifndef MODELMAKERDEFS_HPP
#define MODELMAKERDEFS_HPP

#include <string>
#include <json/json.h>

using std::string;

namespace modelmaker {

class unknown;

class Model {
	friend unknown;
	public:
		bool loadFile(string path);
		void writeFile(string path);
		string write();
	protected:
		virtual json_object* buildJsonObj() = 0;
		virtual bool load(json_object *obj) = 0;
};

class unknown: public Model {
	private:
		json_object *m_obj;
	public:
		unknown();
		~unknown();

		bool loaded();
		bool load(json_object *obj);
		json_object* buildJsonObj();

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

#endif
`
	return out
}

func (me *Out) buildModelmakerDefsBody() string {
	out := `//Generated Code

#include <fstream>
#include "modelmakerdefs.hpp"

using namespace modelmaker;

bool Model::loadFile(string path) {
	std::ifstream in;
	in.open(path.c_str());
	string json;
	if (in.is_open()) {
		while (in.good()) {
			in >> json;
		}
		in.close();
		return true;
	}
	return false;
}

void Model::writeFile(string path) {
	std::ofstream out;
	out.open(path.c_str());
	string json = write();
	out << json << "\n";
	out.close();
}

string Model::write() {
	json_object *obj = buildJsonObj();
	string out = json_object_to_json_string(obj);
	json_object_put(obj);
	return out;
}

unknown::unknown() {
	m_obj = 0;
}

unknown::~unknown() {
	json_object_put(m_obj);
}

bool unknown::load(json_object *obj) {
	//clone the input object because it will get deleted with its parent
	m_obj = json_tokener_parse(json_object_to_json_string(obj));
	return true;
}

json_object* unknown::buildJsonObj() {
	return m_obj;
}

bool unknown::loaded() {
	return m_obj;
}

bool unknown::isBool() {
	return m_obj && json_object_get_type(m_obj) == json_type_boolean;
}

bool unknown::isInt() {
	return m_obj && json_object_get_type(m_obj) == json_type_int;
}

bool unknown::isDouble() {
	return m_obj && json_object_get_type(m_obj) == json_type_double;
}

bool unknown::isObject() {
	return m_obj && json_object_get_type(m_obj) == json_type_object;
}

bool unknown::toBool() {
	return json_object_get_boolean(m_obj);
}

int unknown::toInt() {
	return json_object_get_int(m_obj);
}

double unknown::toDouble() {
	return json_object_get_double(m_obj);
}

string unknown::toString() {
	return json_object_get_string(m_obj);
}

void unknown::set(Model *v) {
	json_object *obj = v->buildJsonObj();
	json_object *old = m_obj;
	m_obj = obj;
	json_object_put(old);
}

void unknown::set(bool v) {
	json_object *obj = json_object_new_boolean(v);
	json_object *old = m_obj;
	m_obj = obj;
	json_object_put(old);
}

void unknown::set(int v) {
	json_object *obj = json_object_new_int(v);
	json_object *old = m_obj;
	m_obj = obj;
	json_object_put(old);
}

void unknown::set(double v) {
	json_object *obj = json_object_new_double(v);
	json_object *old = m_obj;
	m_obj = obj;
	json_object_put(old);
}

void unknown::set(string v) {
	json_object *obj = json_object_new_string(v.c_str());
	json_object *old = m_obj;
	m_obj = obj;
	json_object_put(old);
}
`
	return out
}
