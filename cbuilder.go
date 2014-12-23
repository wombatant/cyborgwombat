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
	"strings"
)

type Cpp struct {
	hppPrefix   string
	hpp         string
	constructor string
	reader      string
	writer      string
	equals      string
	notEquals   string
	namespace   string
	lowerCase   bool
	lib         int
}

func NewCOut(namespace string, lib int, lowerCase bool) *Cpp {
	out := new(Cpp)
	out.lowerCase = lowerCase
	out.namespace = namespace
	out.lib = lib
	out.hppPrefix = `

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
		return "unknown"
	default:
		return me.namespace + "::" + t
	}
	return t
}

func (me *Cpp) namespaceOpen() string {
	out := ""
	list := strings.Split(me.namespace, "::")
	for _, v := range list {
		out += "namespace " + v + " {\n"
	}
	return out
}

func (me *Cpp) namespaceClose() string {
	out := ""
	list := strings.Split(me.namespace, "::")
	for _, _ = range list {
		out += "}\n"
	}
	return out + "\n"
}

func (me *Cpp) buildTypeDec(t string, index []parser.VarType) string {
	array := ""
	out := ""
	for i := 0; i < len(index); i++ {
		if index[i].Type == "slice" {
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

func (me *Cpp) buildVar(v, t string, index []parser.VarType) string {
	return me.buildTypeDec(t, index) + " " + v + ";"
}

func (me *Cpp) addVar(v string, index []parser.VarType) {
	if me.lowerCase && len(v) > 0 && v[0] < 91 {
		v = string(v[0]+32) + v[1:]
	}
	t := me.typeMap(index[len(index)-1].Type)
	index = index[:len(index)-1]
	me.hpp += "\t\t" + me.buildVar(v, t, index) + "\n"
	me.constructor += me.buildConstructor(v, t, index)
	me.reader += me.buildReader(v)
	me.writer += me.buildWriter(v)
	me.equals += me.buildEquals(v)
	me.notEquals += me.buildNotEquals(v)
}

func (me *Cpp) addClass(v string) {
	me.hpp += "\nclass " + v + " {\n"
	me.hpp += "\n\tpublic:\n"
	me.hpp += "\n\t\t" + v + "();\n"
	me.hpp += "\n\t\tbool operator==(const " + v + "&) const;\n"
	me.hpp += "\n\t\tbool operator!=(const " + v + "&) const;\n"

	me.constructor += v + "::" + v + "() {\n"
	me.reader += `inline Error fromJson(` + v + ` *model, json_t *jo) {
	Error err = Error::Ok;
`
	me.writer += "inline Error toJson(" + v + ` model, json_t *jo) {
	Error err = Error::Ok;
`
	me.equals += "bool " + v + "::operator==(const " + v + " &o) const {\n"
	me.notEquals += "bool " + v + "::operator!=(const " + v + " &o) const {\n"
}

func (me *Cpp) closeClass(v string) {
	me.hpp += "};\n\n"
	me.constructor += "}\n\n"
	me.reader += "\treturn err;\n}\n\n"
	me.writer += "\treturn err;\n}\n\n"
	me.equals += "\n\treturn true;\n}\n\n"
	me.notEquals += "\n\treturn false;\n}\n\n"
}

func (me *Cpp) header(fileName string) string {
	n := strings.ToUpper(fileName)
	n = strings.Replace(n, ".", "_", -1)
	out := `//Generated Code

#ifndef ` + n + `
#define ` + n + `

#include "json_read.hpp"
#include "json_write.hpp"

` + me.hppPrefix + me.namespaceOpen() +
		me.hpp + me.writer + me.reader +
		me.namespaceClose() + `
#endif`
	return out
}

func (me *Cpp) body(headername string) string {
	include := ""
	if headername != "" {
		include += `//Generated Code
#include "string.h"
#include "` + headername + `"

`
	}
	return include + me.namespaceOpen() + me.constructor + me.equals + me.notEquals + me.namespaceClose()
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
	}
	return ""
}

func (me *Cpp) buildWriter(v string) string {
	return "\terr |= writeVal(jo, \"" + v + "\", model." + v + ");\n"
}

func (me *Cpp) buildReader(v string) string {
	return "\terr |= readVal(jo, \"" + v + "\", &model->" + v + ");\n"
}

func (me *Cpp) buildEquals(v string) string {
	return "\tif (" + v + " != o." + v + ") return false;\n"
}

func (me *Cpp) buildNotEquals(v string) string {
	return "\tif (" + v + " != o." + v + ") return true;\n"
}
