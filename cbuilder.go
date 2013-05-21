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

using namespace std::string;
using namespace std::vector;
`
	out.types = []string{"int", "bool", "float", "float32", "float64", "double"}
	return out
}

func (me *Out) buildVar(v string, t string, index int) string {
	array := ""
	out := ""
	for i := 0; i < index; i++ {
		out += "vector<"
		array += ">"
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
}

func (me *Out) addClass(v string) {
	me.hpp += "\nclass " + v + " {\n"
	me.hpp += "\n\tpublic:\n"
	me.reader += "void " + v + "::load(json_object *obj) {\n"
	me.writer += "string " + v + "::read() {\n"
}

func (me *Out) closeClass() {
	me.hpp += "\n\t\tvoid load(json_object *obj);\n"
	me.hpp += "\t\tstring write();\n"
	me.hpp += "};\n\n"
	me.reader += "}\n\n"
	me.writer += "}\n\n"
}

func (me *Out) header() string {
	return me.hpp
}

func (me *Out) body(headername string) string {
	include := ""
	if headername != "" {
		include = `#include "` + headername + "\"\n\n"
	}
	return include + me.reader + me.writer[:len(me.writer)-1]
}

func (me *Out) endsWithClose() bool {
	return (me.hpp)[len(me.hpp)-3:] != "};\n"
}

func (me *Out) buildReader(v string, t string, index int, tabs int) string {
	tabs += 2
	out := "out" + strconv.Itoa(index)
	tab := ""
	for n := 0; n < tabs; n++ {
		tab += "\t"
	}
	reader := ""
	//var jtype string
	if index > 0 {
		if tabs != 2 {
			reader += tab[1:] + me.buildVar(out, t, index)
		}
		reader += tab[1:] + "{\n"
		reader += tab + "int size = json_object_array_length(obj);\n"
		if tabs == 2 {
			reader += tab + me.buildVar(out, t, index)
		}
		reader += tab + "for (int i = 0; i < size; i++) {\n"
		reader += tab + "\tjson_object *obj = json_object_array_get_idx(obj, i);\n"
		reader += me.buildReader(v, t, index-1, tabs)
		if tabs == 2 {
			reader += tab + "\tthis->" + v + ".push_back(out" + strconv.Itoa(index-1) + ");\n"
		} else {
			reader += tab + "\t" + out + ".push_back(out" + strconv.Itoa(index-1) + ");\n"
		}
		reader += tab + "}\n"
		reader += tab[1:] + "}\n"
	} else {
		primitive := true
		i := 0
		if tabs == 2 {
			reader += tab[1:] + "{\n"
		} else {
			i += 1
		}
		switch t { //type
		case "float", "float32", "float64", "double":
			reader += tab[i:] + "double out0 = json_object_get_double(obj);\n"
		case "int":
			reader += tab[i:] + "int out0 = json_object_get_int(obj);\n"
		case "bool":
			reader += tab[i:] + "bool out0 = json_object_get_bool(obj);\n"
		case "string":
			reader += tab[i:] + "string out0 = json_object_get_string(obj);\n"
		default:
			primitive = false
			reader += tab + "this->" + v + ".load(obj);\n"
		}
		if tabs == 2 && primitive {
			reader += tab + "this->" + v + " = out0;\n"
		}
		if tabs == 2 {
			reader += tab[1:] + "}\n"
		}
	}
	return reader
}

func (me *Out) addWriter(v string, t string, index int) {
	switch t { //type
	case "int":
	}
}
