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
	"strings"
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
	out.types = []string{"int", "uint", "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "bool", "float32", "float64"}
	return out
}

func (me *Out) buildVar(v string, t string, index int) string {
	array := ""
	out := "\t"
	for i := 0; i < index; i++ {
		out += "vector<"
		array += ">"
	}
	out += t + array + " " + strings.ToLower(v) + ";\n"
	return out
}

func (me *Out) addVar(v string, t string, index int) {
	me.hpp += me.buildVar(v, t, index)
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

func (me *Out) addReader(v string, t string, index int, tabs int) string {
	i := func(is int) string {
		out := ""
		for i := 0; i < is; i++ {
			out += "i"
		}
		return out
	}
	tab := ""
	for n := 0; n < tabs; n++ {
		tab += "\t"
	}
	reader := "\t" + t + " out;\n{"
	if index > 0 {
		reader = "\tjson_object *array;\n"
		reader += tab + "int size = json_object_array_length(array);\n"
		reader += tab + me.buildVar("list", t, index)
		reader += tab + "for (int " + i(index) + " = 0; " + i(index) + " < size; " + i(index) + "++) {"
		reader += tab + "\t" + me.buildVar("var", t, index-1)
		reader += tab + "\t" + me.addReader(v, t, index-1, tabs+1)
		reader += tab + "\t" + "list.push_back(out);\n"
		reader += tab + "}"
	} else {
		switch t { //type
		case "int":
		}
	}
	reader += "}\n"
	return reader
}

func (me *Out) addWriter(v string, t string, index int) {
	switch t { //type
	case "int":
	}
}
