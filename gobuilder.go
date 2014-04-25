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
	"github.com/gtalent/cyborgbear/parser"
	"io/ioutil"
)

type Go struct {
	text string
	pkg  string
}

func NewGo(pkg string) *Go {
	out := new(Go)
	out.pkg = pkg
	return out
}

func (me *Go) addClass(name string) {
	me.text += `
type ` + name + ` struct {`
}

func (me *Go) addVar(name string, t []parser.VarType) {
	vtype := ""
	for _, v := range t {
		switch v.Type {
		case "array":
			vtype += "[" + v.Index + "]"
		case "map":
			vtype += "map[" + v.Index + "]"
		case "slice":
			vtype += "[]"
		case "unknown":
			vtype += "interface{}"
		default:
			vtype += v.Type
		}
	}
	me.text += "\n\t" + name + " " + vtype
}

func (me *Go) closeClass(name string) {
	me.text += `
}

func (me *` + name + `) FromJSON(text []byte) error {
	return json.Unmarshal(text, me)
}

func (me *` + name + `) ReadJSONFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		 return err
	}
	return json.Unmarshal(file, me)
}

func (me *` + name + `) ToJSON() []byte {
	out, _ := json.Marshal(me)
	return out
}

func (me *` + name + `) WriteJSONFile(path string) error {
	out, _ := json.Marshal(me)
	return ioutil.WriteFile(path, out, 0644)
}

func (me *` + name + `) ToGob() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(me)
	return buf.Bytes()
}

func (me *` + name + `) WriteGobFile(path string) error {
	out := me.ToGob()
	return ioutil.WriteFile(path, out, 0644)
}

func (me *` + name + `) FromGob(data []byte) error {
	var buf bytes.Buffer
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(me)
	if err != nil {
		return err
	}
	return err
}

func (me *` + name + `) ReadGobFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		 return err
	}
	return me.FromGob(file)
}

func (me *` + name + `) ReadObjFile(path string) error {
	switch path[strings.LastIndex(path, ".")+1:] {
	case "json":
		return me.ReadJSONFile(path)
	case "gob":
		return me.ReadGobFile(path)
	default:
		err := me.ReadGobFile(path+".gob")
		if err == nil {
			return nil
		}
		err = me.ReadJSONFile(path+".json")
		if err == nil {
			return nil
		}
		return err
	}
	return nil
}
`
}

func (me *Go) writeFile(s string) error {
	return ioutil.WriteFile(s, []byte(me.write("")), 0644)
}

func (me *Go) write(s string) string {
	out := `package ` + me.pkg + `

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"strings"
)

` + me.text
	return out
}
