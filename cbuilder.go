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
#include "types.h"

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
	me.hpp += "\t" + me.buildVar(v, t, index)
	me.reader += me.buildReader(v, t, index, 0)
}

func (me *Out) addClass(v string) {
	me.hpp += "\nclass " + v + " {\n"
	me.hpp += "\tvoid read(string json);\n"
	me.hpp += "\tstring write();\n"
}

func (me *Out) closeClass() {
	me.hpp += "};\n"
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
	var jfunc string
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
		reader += tab + "\tjson_object obj* = json_object_array_get_idx(obj, i);\n"
		reader += me.buildReader(v, t, index-1, tabs)
		reader += tab + "\t" + out + ".push_back(out" + strconv.Itoa(index-1) + ");\n"
		reader += tab + "}\n"
		if tabs == 2 {
			reader += tab + "this." + v + " = " + out + ";\n"
		}
		reader += tab[1:] + "}\n"
	} else {
		switch t { //type
		case "float", "float32", "float64", "double":
			t = "double"
			jfunc = "json_object_get_double(obj);"
		case "int":
			t = "int"
			jfunc = "json_object_get_int(obj);"
		case "bool":
			t = "bool"
			jfunc = "json_object_get_bool(obj);"
		case "string":
			t = "string"
			jfunc = "json_object_get_string(obj);"
		}
		reader += tab[1:] + t + " " + out + " = " + jfunc + "\n"
	}
	return reader
}

func (me *Out) addWriter(v string, t string, index int) {
	switch t { //type
	case "int":
	}
}
