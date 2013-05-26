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
)

type Out struct {
	hpp         string
	constructor string
	reader      string
	writer      string
}

func NewCOut() Out {
	var out Out
	out.hpp = `//Generated Code

#include <string>
#include <vector>
#include <json/json.h>
#include "modelmakerdefs.hpp"


using std::string;
using std::vector;
`
	return out
}

func (me *Out) typeMap(t string) string {
	switch t {
	case "float", "float32", "float64", "double":
		return "double"
	default:
		return t
	}
	return ""
}

func (me *Out) buildVar(v, t string, index int) string {
	array := ""
	out := ""
	for i := 0; i < index; i++ {
		out += "vector<"
		array += " >"
	}
	out += t + array + " " + v + ";"
	return out
}

func (me *Out) addVar(v, t string, index int) {
	jsonV := v
	if len(v) > 0 && v[0] < 91 {
		v = string(v[0]+32) + v[1:]
	}
	t = me.typeMap(t)
	me.hpp += "\t\t" + me.buildVar(v, t, index) + "\n"
	var reader CppCode
	reader.tabs += "\t"
	me.constructor += me.buildConstructor(v, t, index)
	me.reader += me.buildReader(&reader, v, jsonV, t, index, 0)
	me.writer += me.buildWriter(v, jsonV, t, index)
}

func (me *Out) addClass(v string) {
	me.hpp += "\nnamespace models {\n"
	me.hpp += "\nclass " + v + ": public Model {\n"
	me.hpp += "\n\tpublic:\n"
	me.hpp += "\n\t\t" + v + "();\n"
	me.hpp += "\n\t\tvoid load(string text);\n"
	me.hpp += "\n\t\tstring write();\n"
	me.hpp += "\n\t\tbool load(json_object *obj);\n"
	me.hpp += "\n\t\tjson_object* buildJsonObj();\n\n"
	me.reader += "void " + v + `::load(string json) {
	json_object *obj = json_tokener_parse(json.c_str());
	load(obj);
	json_object_put(obj);
}

`
	me.writer += "string " + v + `::write() {
	json_object *obj = buildJsonObj();
	string out = json_object_to_json_string(obj);
	json_object_put(obj);
	return out;
}

`
	me.constructor += v + "::" + v + "() {\n"
	me.reader += "bool " + v + "::load(json_object *in) {"
	me.writer += "json_object* " + v + `::buildJsonObj() {
	json_object *obj = json_object_new_object();`
}

func (me *Out) closeClass() {
	me.hpp += "};\n\n"
	me.hpp += "};\n\n"
	me.constructor += "}\n\n"
	me.reader += "\n\treturn true;\n}\n\n"
	me.writer += "\n\treturn obj;\n}\n\n"
}

func (me *Out) header() string {
	return me.hpp
}

func (me *Out) body(headername string) string {
	include := ""
	if headername != "" {
		include += `//Generated Code

#include "` + headername + `"

using namespace models;

`
	}
	return include + me.constructor + me.reader + me.writer[:len(me.writer)-1]
}

func (me *Out) endsWithClose() bool {
	return (me.hpp)[len(me.hpp)-3:] != "};\n"
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

func (me *Out) buildReader(code *CppCode, v, jsonV, t string, index, depth int) string {
	depth += 1
	out := "out" + strconv.Itoa(index)

	if index > 0 {
		if depth != 1 {
			code.Insert(me.buildVar(out, t, index))
		}
		code.PushBlock()
		if depth == 1 {
			code.Insert("json_object *obj" + strconv.Itoa(index) + " = json_object_object_get(in, \"" + jsonV + "\");")
		}
		code.PushIfBlock("obj" + strconv.Itoa(index) + " != NULL")
		code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(index) + ") != json_type_array")
		code.Insert("return false;")
		code.PopBlock()
		code.Insert("int size = json_object_array_length(obj" + strconv.Itoa(index) + ");")
		if depth == 1 {
			code.Insert(me.buildVar(out, t, index))
		}
		code.PushForBlock("int i = 0; i < size; i++")
		code.Insert("json_object *obj" + strconv.Itoa(index-1) + " = json_object_array_get_idx(obj" + strconv.Itoa(index) + ", i);")

		me.buildReader(code, v, jsonV, t, index-1, depth)
		if depth == 1 {
			code.Insert("this->" + v + ".push_back(out" + strconv.Itoa(index-1) + ");")
		} else {
			code.Insert("" + out + ".push_back(out" + strconv.Itoa(index-1) + ");")
		}
		code.PopBlock()
		code.PopBlock()
		code.PopBlock()
	} else {
		primitive := true
		i := 0
		if depth == 1 {
			code.PushBlock()
			code.Insert("json_object *obj" + strconv.Itoa(index) + " = json_object_object_get(in, \"" + jsonV + "\");")
			code.PushIfBlock("obj" + strconv.Itoa(index) + " != NULL")
		} else {
			i += 2
		}
		switch t { //type
		case "int", "double":
			code.Insert(t + " out0;")
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(index) + ") == " + "json_type_" + t)
			code.Insert("out0 = json_object_get_" + t + "(obj" + strconv.Itoa(index) + ");")
			code.PopBlock()
		case "bool":
			code.Insert(t + " out0;")
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(index) + ") == " + "json_type_boolean")
			code.Insert("out0 = json_object_get_boolean(obj" + strconv.Itoa(index) + ");")
			code.PopBlock()
		case "string":
			code.Insert("string out0;")
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(index) + ") == " + "json_type_string")
			code.Insert("out0 = json_object_get_string(obj" + strconv.Itoa(index) + ");")
			code.PopBlock()
		default:
			primitive = false
			code.Insert(t + " out0;")
			code.PushIfBlock("json_object_get_type(obj" + strconv.Itoa(index) + ") == " + "json_type_object")
			code.Insert("out0.load(obj" + strconv.Itoa(index) + ");")
			code.PopBlock()
		}
		if depth == 1 && primitive {
			code.Insert("this->" + v + " = out0;")
		}
		if depth == 1 {
			code.PopBlock()
			code.PopBlock()
		}
	}
	return code.String()
}

func (me *Out) buildArrayWriter(out *CppCode, t, v string, depth, index int) {
	sub := "[i]"
	is := "i"
	list := "this->" + v
	for i := 0; i < depth; i++ {
		if i < depth {
			list += "[" + is + "]"
		}
		is += "i"
		sub += "[" + is + "]"
	}

	out.Insert("json_object *array" + strconv.Itoa(depth) + " = json_object_new_array();")
	out.PushForBlock("int " + is + " = 0; " + is + " < " + list + ".size(); " + is + "++")
	if index != 0 {
		me.buildArrayWriter(out, t, v, depth+1, index-1)
		out.Insert("json_object_array_add(array" + strconv.Itoa(depth) + ", array" + strconv.Itoa(depth+1) + ");")
	} else {
		switch t {
		case "int", "double":
			out.Insert("json_object *out0 = json_object_new_" + t + "(this->" + v + sub + ");")
		case "bool":
			out.Insert("json_object *out0 = json_object_new_boolean(this->" + v + sub + ");")
		case "string":
			out.Insert("json_object *out0 = json_object_new_string(this->" + v + sub + ".c_str());")
		default:
			out.Insert("json_object *out0 = this->" + v + sub + ".buildJsonObj();")
		}
		out.Insert("json_object_array_add(array" + strconv.Itoa(depth) + ", out0);")
	}
	out.PopBlock()
}

func (me *Out) buildWriter(v, jsonV, t string, index int) string {
	var out CppCode
	out.tabs = "\t"
	out.PushBlock()
	if index > 0 {
		me.buildArrayWriter(&out, t, v, 0, index-1)
		out.Insert("json_object_object_add(obj, \"" + jsonV + "\", array0);")
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

#include <string>
#include <json/json.h>

using std::string;

namespace models {

class unknown;

class Model {
	friend unknown;
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
`
	return out
}

func (me *Out) buildModelmakerDefsBody() string {
	out := `//Generated Code

#include "modelmakerdefs.hpp"

using namespace models;

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
