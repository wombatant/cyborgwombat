/*
   Copyright 2013-2014 gtalent2@gmail.com

   This Source Code Form is subject to the terms of the Mozilla Public
   License, v. 2.0. If a copy of the MPL was not distributed with this
   file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package models

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"strings"
)


type Model1 struct {
	Field1 string
	Field2 interface{}
	Field3 [4]int
	Field4 [][]string
	Field5 map[string]string
}

func (me *Model1) FromJSON(text []byte) error {
	return json.Unmarshal(text, me)
}

func (me *Model1) ReadJSONFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		 return err
	}
	return json.Unmarshal(file, me)
}

func (me *Model1) ToJSON() []byte {
	out, _ := json.Marshal(me)
	return out
}

func (me *Model1) WriteJSONFile(path string) error {
	out, _ := json.Marshal(me)
	return ioutil.WriteFile(path, out, 0644)
}

func (me *Model1) ToGob() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(me)
	return buf.Bytes()
}

func (me *Model1) WriteGobFile(path string) error {
	out := me.ToGob()
	return ioutil.WriteFile(path, out, 0644)
}

func (me *Model1) FromGob(data []byte) error {
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

func (me *Model1) ReadGobFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		 return err
	}
	return me.FromGob(file)
}

func (me *Model1) ReadObjFile(path string) error {
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
