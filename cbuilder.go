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
	types  []string
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
	out.types = []string{"int", "bool", "float", "float32", "float64", "double"}
	return out
}

func (me *Out) buildVar(v string, t string, index int) string {
	array := ""
	out := ""
	for i := 0; i < index; i++ {
		out += "vector<"
		array += " >"
	}
	out += t + array + " " + v + ";\n"
	return out
}

func (me *Out) addVar(v string, t string, index int) {
	if len(v) > 0 && v[0] < 91 {
		v = string(v[0]+32) + v[1:]
	}
	me.hpp += "\t\t" + me.buildVar(v, t, index)
	me.reader += me.buildReader(v, t, index, 0)
	me.writer += me.buildWriter(v, t, index)
}

func (me *Out) addClass(v string) {
	me.hpp += "\nclass " + v + " {\n"
	me.hpp += "\n\tpublic:\n"
	me.reader += "void " + v + `::load(string json) {
	json_object *obj = json_tokener_parse(json.c_str());
	load(obj);
	free(obj);
}

`
	me.reader += "bool " + v + "::load(json_object *in) {\n"
	me.writer += "json_object* " + v + `::write() {
	json_object *obj = json_object_new_object();`
}

func (me *Out) closeClass() {
	me.hpp += "\n\t\tvoid load(string path);\n"
	me.hpp += "\n\t\tbool load(json_object *obj);\n"
	me.hpp += "\n\t\tjson_object* write();\n"
	me.hpp += "};\n\n"
	me.reader += "\treturn true;\n}\n\n"
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

func (me *Out) buildReader(v, t string, index, depth int) string {
	depth += 1
	out := "out" + strconv.Itoa(index)
	tab := ""
	tabs := depth * 2
	for n := 0; n < tabs; n++ {
		tab += "\t"
	}
	reader := ""

	if index > 0 {
		if depth != 1 {
			reader += tab[1:] + me.buildVar(out, t, index)
		}
		reader += tab[1:] + "{\n"
		if depth == 1 {
			reader += tab[0:] + "json_object *obj" + strconv.Itoa(index) + " = json_object_object_get(in, \"" + v + "\");\n"
		}
		reader += tab[0:] + "if (json_object_get_type(obj" + strconv.Itoa(index) + ") != json_type_array) {\n"
		reader += tab + "\treturn false;\n"
		reader += tab + "}\n"
		reader += tab + "int size = json_object_array_length(obj" + strconv.Itoa(index) + ");\n"
		if depth == 1 {
			reader += tab + me.buildVar(out, t, index)
		}
		reader += tab + "for (int i = 0; i < size; i++) {\n"
		reader += tab + "\tjson_object *obj" + strconv.Itoa(index-1) + " = json_object_array_get_idx(obj" + strconv.Itoa(index) + ", i);\n"
		reader += me.buildReader(v, t, index-1, depth)
		if depth == 1 {
			reader += tab + "\tthis->" + v + ".push_back(out" + strconv.Itoa(index-1) + ");\n"
		} else {
			reader += tab + "\t" + out + ".push_back(out" + strconv.Itoa(index-1) + ");\n"
		}
		reader += tab[0:] + "}\n"
		reader += tab[1:] + "}\n"
	} else {
		primitive := true
		i := 0
		if depth == 1 {
			reader += tab[1:] + "{\n"
			reader += tab + "json_object *obj" + strconv.Itoa(index) + " = json_object_object_get(in, \"" + v + "\");\n"
			reader += tab + "{\n"
		} else {
			i += 2
		}
		reader += tab[i:] + "\tif (json_object_get_type(obj" + strconv.Itoa(index) + ") != "
		switch t { //type
		case "float", "float32", "float64", "double":
			reader += "json_type_double) return false;\n"
			reader += tab[i:] + "\tdouble out0 = json_object_get_double(obj" + strconv.Itoa(index) + ");\n"
		case "int":
			reader += "json_type_int) return false;\n"
			reader += tab[i:] + "\tint out0 = json_object_get_int(obj" + strconv.Itoa(index) + ");\n"
		case "bool":
			reader += "json_type_bool) return false;\n"
			reader += tab[i:] + "\tbool out0 = json_object_get_bool(obj" + strconv.Itoa(index) + ");\n"
		case "string":
			reader += "json_type_string) return false;\n"
			reader += tab[i:] + "\tstring out0 = json_object_get_string(obj" + strconv.Itoa(index) + ");\n"
		default:
			primitive = false
			reader += "json_type_object) return false;\n"
			reader += tab[i:] + "\t" + t + " out0;\n"
			reader += tab[i:] + "\tout0.load(obj" + strconv.Itoa(index) + ");\n"
		}
		if depth == 1 && primitive {
			reader += tab + "\tthis->" + v + " = out0;\n"
		}
		if depth == 1 {
			reader += tab + "}\n"
			reader += tab[1:] + "}\n"
		}
	}
	return reader
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
		case "float", "float32", "float64", "double":
			out.Insert("json_object *out0 = json_object_new_double(this->" + v + sub + ");")
		case "int":
			out.Insert("json_object *out0 = json_object_new_int(this->" + v + sub + ");")
		case "bool":
			out.Insert("json_object *out0 = json_object_new_bool(this->" + v + sub + ");")
		case "string":
			out.Insert("json_object *out0 = json_object_new_string(this->" + v + sub + ".c_str());")
		default:
			out.Insert("json_object *out0 = this->" + v + sub + ".write();")
		}
		out.Insert("json_object_array_add(array" + strconv.Itoa(depth) + ", out0);")
		out.Insert("free(out0);")
	}
	out.PopBlock()
}

func (me *Out) buildWriter(v, t string, index int) string {
	var out CppCode
	out.tabs = "\t"
	out.PushBlock()
	if index > 0 {
		me.buildArrayWriter(&out, t, v, 0, index-1)
		out.Insert("json_object_object_add(obj, \"" + v + "\", array0);")
	} else {
		switch t {
		case "float", "float32", "float64", "double":
			out.Insert("json_object *out0 = json_object_new_double(this->" + v + ");")
		case "int":
			out.Insert("json_object *out0 = json_object_new_int(this->" + v + ");")
		case "bool":
			out.Insert("json_object *out0 = json_object_new_bool(this->" + v + ");")
		case "string":
			out.Insert("json_object *out0 = json_object_new_string(this->" + v + ".c_str());")
		default:
		}
		out.Insert("json_object_object_add(obj, \"" + v + "\", out0);")
		out.Insert("free(out0);")
	}
	out.PopBlock()
	return out.String()
}
