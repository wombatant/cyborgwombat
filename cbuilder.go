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

type CppJansson struct {
	hppPrefix   string
	hpp         string
	constructor string
	reader      string
	writer      string
	namespace   string
}

func NewCOut(namespace string) *CppJansson {
	out := new(CppJansson)
	out.namespace = namespace
	out.hppPrefix = `#include <string>
#include <sstream>
#include <vector>
#include <map>
#include <jansson.h>
#include "modelmakerdefs.hpp"


using std::string;
using std::vector;
using std::map;
`
	return out
}

func (me *CppJansson) typeMap(t string) string {
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

func (me *CppJansson) buildTypeDec(t string, index []string) string {
	array := ""
	out := ""
	for i := 0; i < len(index); i++ {
		if index[i] == "array" {
			out += "vector< "
			array += " >"
		} else if index[i][:3] == "map" {
			out += "map< " + index[i][4:] + ", "
			array += " >"
		}
	}
	out += t + array
	return out
}

func (me *CppJansson) buildVar(v, t string, index []string) string {
	return me.buildTypeDec(t, index) + " " + v + ";"
}

func (me *CppJansson) addVar(v string, index []string) {
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

func (me *CppJansson) addClass(v string) {
	me.hpp += "\nnamespace " + me.namespace + " {\n"
	me.hpp += "\nclass " + v + ": public modelmaker::Model {\n"
	me.hpp += "\n\tpublic:\n"
	me.hpp += "\n\t\t" + v + "();\n"
	me.hpp += "\n\t\tbool load_json_t(json_t *obj);\n"
	me.hpp += "\n\t\tjson_t* buildJsonObj();\n\n"
	me.constructor += v + "::" + v + "() {\n"
	me.reader += "bool " + v + "::load_json_t(json_t *in) {"
	me.writer += "json_t* " + v + `::buildJsonObj() {
	json_t *obj = json_object();`
}

func (me *CppJansson) closeClass() {
	me.hpp += "};\n\n"
	me.hpp += "}\n\n"
	me.constructor += "}\n\n"
	me.reader += "\n\treturn true;\n}\n\n"
	me.writer += "\n\treturn obj;\n}\n\n"
}

func (me *CppJansson) header(fileName string) string {
	n := strings.ToUpper(fileName)
	n = strings.Replace(n, ".", "_", -1)
	return `//Generated Code
#ifndef ` + n + `
#define ` + n + `
` + me.hppPrefix + me.hpp + `
#endif`
}

func (me *CppJansson) body(headername string) string {
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

func (me *CppJansson) endsWithClose() bool {
	return me.hpp[len(me.hpp)-3:] == "};\n"
}

func (me *CppJansson) buildConstructor(v, t string, index int) string {
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

func (me *CppJansson) buildReader(code *CppCode, v, jsonV, t, sub string, index []string, depth int) string {
	if depth == 0 {
		code.PushBlock()
		code.Insert("json_t *obj0 = json_object_get(in, \"" + jsonV + "\");")
	}
	if len(index) > 0 {
		is := "i"
		for i := 0; i < depth; i++ {
			is += "i"
		}
		if index[0] == "array" {
			code.PushIfBlock("obj" + strconv.Itoa(depth) + " != NULL && json_typeof(obj" + strconv.Itoa(depth) + ") == JSON_ARRAY")
			code.Insert("unsigned int size = json_array_size(obj" + strconv.Itoa(depth) + ");")
			code.Insert("this->" + v + sub + ".resize(size);")
			code.PushForBlock("unsigned int " + is + " = 0; " + is + " < size; " + is + "++")
			code.Insert("json_t *obj" + strconv.Itoa(depth+1) + " = json_array_get(obj" + strconv.Itoa(depth) + ", " + is + ");")
			me.buildReader(code, v, jsonV, t, sub+"["+is+"]", index[1:], depth+1)
			code.PopBlock()
			code.PopBlock()
		} else if index[0][:3] == "map" {
			code.PushIfBlock("obj" + strconv.Itoa(depth) + " != NULL && json_typeof(obj" + strconv.Itoa(depth) + ") == JSON_OBJECT")
			code.Insert("const char *key;")
			code.Insert("json_t *obj" + strconv.Itoa(depth+1) + ";")
			code.PushPrefixBlock("json_object_foreach(obj" + strconv.Itoa(depth) + ", key, obj" + strconv.Itoa(depth+1) + ")")
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

			//initialize value in map
			code.Insert(me.buildTypeDec(t, index[1:]) + " val;")
			code.Insert("this->" + v + sub + ".insert(std::make_pair(" + is + ", val));")

			me.buildReader(code, v, jsonV, t, sub+"["+is+"]", index[1:], depth+1)
			code.PopBlock()
			code.PopBlock()
		}
	} else {
		code.PushBlock()
		switch t { //type
		case "int":
			code.PushIfBlock("json_is_integer(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = (int) json_integer_value(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "double":
			code.PushIfBlock("json_is_real(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = json_real_value" + t + "(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "bool":
			code.PushIfBlock("json_is_boolean(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = json_is_true(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "string":
			code.PushIfBlock("json_is_string(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + " = json_string_value(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		case "modelmaker::unknown":
			code.Insert("this->" + v + sub + ".load_json_t(obj" + strconv.Itoa(depth) + ");")
		default:
			code.PushIfBlock("json_is_object(obj" + strconv.Itoa(depth) + ")")
			code.Insert("this->" + v + sub + ".load_json_t(obj" + strconv.Itoa(depth) + ");")
			code.PopBlock()
		}
		code.PopBlock()
	}

	if depth == 0 {
		code.PopBlock()
	}

	return code.String()
}

func (me *CppJansson) buildArrayWriter(code *CppCode, t, v, sub string, depth int, index []string) {
	is := "i"
	ns := "n"
	for i := 0; i < depth; i++ {
		is += "i"
		ns += "n"
	}

	if len(index) > depth {
		if index[depth] == "array" {
			code.Insert("json_t *out" + strconv.Itoa(len(index[depth:])) + " = json_array();")
			code.PushForBlock("unsigned int " + is + " = 0; " + is + " < this->" + v + sub + ".size(); " + is + "++")
			me.buildArrayWriter(code, t, v, sub+"["+is+"]", depth+1, index)
			code.Insert("json_array_append(out" + strconv.Itoa(len(index[depth:])) + ", out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.Insert("json_decref(out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.PopBlock()
		} else if index[depth][:3] == "map" {
			code.Insert("json_t *out" + strconv.Itoa(len(index[depth:])) + " = json_object();")
			code.PushForBlock(me.buildTypeDec(t, index[depth:]) + "::iterator " + ns + " = this->" + v + sub + ".begin(); " + ns + " != this->" + v + sub + ".end(); ++" + ns)
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
			code.Insert("json_object_set(out" + strconv.Itoa(len(index[depth:])) + ", key.c_str(), out" + strconv.Itoa(len(index[depth:])-1) + ");")
			code.Insert("json_decref(out" + strconv.Itoa(len(index[depth+1:])) + ");")
			code.PopBlock()
		}
	} else {
		switch t {
		case "int":
			code.Insert("json_t *out0 = json_integer(this->" + v + sub + ");")
		case "double":
			code.Insert("json_t *out0 = json_real(this->" + v + sub + ");")
		case "bool":
			code.Insert("json_t *out0 = json_boolean(this->" + v + sub + ");")
		case "string":
			code.Insert("json_t *out0 = json_string(this->" + v + sub + ".c_str());")
		default:
			code.Insert("json_t *out0 = this->" + v + sub + ".buildJsonObj();")
		}
	}
}

func (me *CppJansson) buildWriter(v, jsonV, t string, index []string) string {
	var out CppCode
	out.tabs = "\t"
	out.PushBlock()
	if len(index) > 0 {
		me.buildArrayWriter(&out, t, v, "", 0, index)
		out.Insert("json_object_set(obj, \"" + jsonV + "\", out" + strconv.Itoa(len(index)) + ");")
		out.Insert("json_decref(out" + strconv.Itoa(len(index)) + ");")
	} else {
		switch t {
		case "int":
			out.Insert("json_t *out0 = json_integer(this->" + v + ");")
		case "double":
			out.Insert("json_t *out0 = json_real(this->" + v + ");")
		case "bool":
			out.Insert("json_t *out0 = json_boolean(this->" + v + ");")
		case "string":
			out.Insert("json_t *out0 = json_string(this->" + v + ".c_str());")
		default:
			out.Insert("json_t *out0 = this->" + v + ".buildJsonObj();")
		}
		out.Insert("json_object_set(obj, \"" + jsonV + "\", out0);")
		out.Insert("json_decref(out0);")
	}
	out.PopBlock()
	return out.String()
}

func (me *CppJansson) buildModelmakerDefsHeader() string {
	out := `//Generated Code

#ifndef MODELMAKERDEFS_HPP
#define MODELMAKERDEFS_HPP

#include <string>
#include <jansson.h>

using std::string;

namespace modelmaker {

class unknown;

class Model {
	friend class unknown;
	public:
		bool loadFile(string path);
		void writeFile(string path);
		void load(string json);
		string write();
	protected:
		virtual json_t* buildJsonObj() = 0;
		virtual bool load_json_t(json_t *obj) = 0;
};

class unknown: public Model {
	private:
		json_t *m_obj;
	public:
		unknown();
		unknown(Model *v);
		unknown(bool v);
		unknown(int v);
		unknown(double v);
		unknown(string v);
		virtual ~unknown();

		bool loaded();
		bool load_json_t(json_t *obj);
		json_t* buildJsonObj();

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

func (me *CppJansson) buildModelmakerDefsBody() string {
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
			string s;
			in >> s;
			json += s;
		}
		in.close();
		load(json);
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

void Model::load(string json) {
	json_t *obj = json_loads(json.c_str(), 0, NULL);
	load_json_t(obj);
	json_decref(obj);
}

string Model::write() {
	json_t *obj = buildJsonObj();
	char *tmp = json_dumps(obj, JSON_COMPACT);
	if (!tmp)
		return "{}";
	string out = tmp;
	free(tmp);
	json_decref(obj);
	return out;
}

unknown::unknown() {
	m_obj = 0;
}

unknown::unknown(Model *v) {
	m_obj = 0;
	set(v);
}

unknown::unknown(bool v) {
	m_obj = 0;
	set(v);
}

unknown::unknown(int v) {
	m_obj = 0;
	set(v);
}

unknown::unknown(double v) {
	m_obj = 0;
	set(v);
}

unknown::unknown(string v) {
	m_obj = 0;
	set(v);
}

unknown::~unknown() {
	json_decref(m_obj);
}

bool unknown::load_json_t(json_t *obj) {
	m_obj = json_incref(obj);
	return obj != 0;
}

json_t* unknown::buildJsonObj() {
	return json_incref(m_obj);
}

bool unknown::loaded() {
	return m_obj;
}

bool unknown::isBool() {
	return m_obj && json_is_boolean(m_obj);
}

bool unknown::isInt() {
	return m_obj && json_is_integer(m_obj);
}

bool unknown::isDouble() {
	return m_obj && json_is_real(m_obj);
}

bool unknown::isString() {
	return m_obj && json_is_string(m_obj);
}

bool unknown::isObject() {
	return m_obj && json_is_object(m_obj);
}

bool unknown::toBool() {
	return json_is_true(m_obj);
}

int unknown::toInt() {
	return json_integer_value(m_obj);
}

double unknown::toDouble() {
	return json_real_value(m_obj);
}

string unknown::toString() {
	return json_string_value(m_obj);
}

void unknown::set(Model *v) {
	json_t *obj = v->buildJsonObj();
	json_t *old = m_obj;
	m_obj = obj;
	if (old) {
		json_decref(old);
	}
}

void unknown::set(bool v) {
	json_t *obj = json_boolean(v);
	json_t *old = m_obj;
	m_obj = obj;
	if (old) {
		json_decref(old);
	}
}

void unknown::set(int v) {
	json_t *obj = json_integer(v);
	json_t *old = m_obj;
	m_obj = obj;
	if (old) {
		json_decref(old);
	}
}

void unknown::set(double v) {
	json_t *obj = json_real(v);
	json_t *old = m_obj;
	m_obj = obj;
	if (old) {
		json_decref(old);
	}
}

void unknown::set(string v) {
	json_t *obj = json_string(v.c_str());
	json_t *old = m_obj;
	m_obj = obj;
	if (old) {
		json_decref(old);
	}
}
`
	return out
}
