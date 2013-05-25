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
	hpp    string
	reader string
	writer string
}

func NewCOut() Out {
	var out Out
	out.hpp = `//Generated Code

#include <string>
#include <vector>
#include <json/json.h>

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

func (me *Out) buildVar(v string, t string, index int) string {
	array := ""
	out := ""
	for i := 0; i < index; i++ {
		out += "vector<"
		array += " >"
	}
	out += t + array + " " + v + ";"
	return out
}

func (me *Out) addVar(v string, t string, index int) {
	jsonV := v
	if len(v) > 0 && v[0] < 91 {
		v = string(v[0]+32) + v[1:]
	}
	t = me.typeMap(t)
	me.hpp += "\t\t" + me.buildVar(v, t, index)
	var reader CppCode
	reader.tabs += "\t"
	me.reader += me.buildReader(&reader, v, jsonV, t, index, 0)
	me.writer += me.buildWriter(v, jsonV, t, index)
}

func (me *Out) addClass(v string) {
	me.hpp += "\nclass " + v + " {\n"
	me.hpp += "\n\tpublic:\n"
	me.hpp += "\n\t\tvoid load(string text);\n"
	me.hpp += "\n\t\tstring write();\n"
	me.hpp += "\n\t\tbool load(json_object *obj);\n"
	me.hpp += "\n\t\tjson_object* buildJsonObj();\n"
	me.reader += "void " + v + `::load(string json) {
	json_object *obj = json_tokener_parse(json.c_str());
	load(obj);
	free(obj);
}

`
	me.writer += "string " + v + `::write() {
	json_object *obj = buildJsonObj();
	string out = json_object_to_json_string(obj);
	free(obj);
	return out;
}

`
	me.reader += "bool " + v + "::load(json_object *in) {"
	me.writer += "json_object* " + v + `::buildJsonObj() {
	json_object *obj = json_object_new_object();`
}

func (me *Out) closeClass() {
	me.hpp += "};\n\n"
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

#include "` + headername + "\"\n\n"
	}
	return include + me.reader + me.writer[:len(me.writer)-1]
}

func (me *Out) endsWithClose() bool {
	return (me.hpp)[len(me.hpp)-3:] != "};\n"
}

func (me *Out) buildReader(code *CppCode, v, jsonV, t string, index, depth int) string {
	depth += 1
	out := "out" + strconv.Itoa(index)
	tab := ""
	tabs := depth * 2
	for n := 0; n < tabs; n++ {
		tab += "\t"
	}

	if index > 0 {
		if depth != 1 {
			code.Insert(me.buildVar(out, t, index))
		}
		code.PushBlock()
		if depth == 1 {
			code.Insert("json_object *obj" + strconv.Itoa(index) + " = json_object_object_get(in, \"" + jsonV + "\");")
		}
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
	} else {
		primitive := true
		i := 0
		if depth == 1 {
			code.PushBlock()
			code.Insert("json_object *obj" + strconv.Itoa(index) + " = json_object_object_get(in, \"" + jsonV + "\");")
			code.PushBlock()
		} else {
			i += 2
		}
		switch t { //type
		case "bool", "int", "double":
			code.Insert("if (json_object_get_type(obj" + strconv.Itoa(index) + ") != " + "json_type_" + t + ") return false;")
			code.Insert(t + " out0 = json_object_get_" + t + "(obj" + strconv.Itoa(index) + ");")
		case "string":
			code.Insert("if (json_object_get_type(obj" + strconv.Itoa(index) + ") != " + "json_type_string) return false;")
			code.Insert("string out0 = json_object_get_string(obj" + strconv.Itoa(index) + ");")
		default:
			primitive = false
			code.Insert("if (json_object_get_type(obj" + strconv.Itoa(index) + ") != " + "json_type_object) return false;")
			code.Insert(t + " out0;")
			code.Insert("out0.load(obj" + strconv.Itoa(index) + ");")
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
		case "bool", "int", "double":
			out.Insert("json_object *out0 = json_object_new_" + t + "(this->" + v + sub + ");")
		case "string":
			out.Insert("json_object *out0 = json_object_new_string(this->" + v + sub + ".c_str());")
		default:
			out.Insert("json_object *out0 = this->" + v + sub + ".buildJsonObj();")
		}
		out.Insert("json_object_array_add(array" + strconv.Itoa(depth) + ", out0);")
		out.Insert("free(out0);")
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
		case "bool", "int", "double":
			out.Insert("json_object *out0 = json_object_new_" + t + "(this->" + v + ");")
		case "string":
			out.Insert("json_object *out0 = json_object_new_string(this->" + v + ".c_str());")
		default:
			out.Insert("json_object *out0 = this->" + v + ".buildJsonObj();")
		}
		out.Insert("json_object_object_add(obj, \"" + jsonV + "\", out0);")
		out.Insert("free(out0);")
	}
	out.PopBlock()
	return out.String()
}
