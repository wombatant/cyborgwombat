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
	"errors"
	"fmt"
	"github.com/gtalent/lex"
)

type VarType struct {
	Type  string
	Index string
}

type Var struct {
	Name string
	Type []VarType
}

type Model struct {
	Name string
	Vars []Var
}

type Parser struct {
	out    Out
	models []*Model
}

func (me *Parser) processVariable(tokens []lex.Token) (int, error) {
	size := 3 // should be 1 less than the actual number parsed

	if len(tokens) < 4 {
		return 0, errors.New("Incomplete variable")
	}

	var t []VarType
	var variable string
	if len(tokens) > 2 {
		variable = tokens[1].String()
		tokens = tokens[2:]
	} else {
		return 0, errors.New("Error: unexpected end of file")
	}
	for {
		if len(tokens) > 0 {
			if tokens[0].String() == " " || tokens[0].String() == "\t" {
				tokens = tokens[1:]
				size++
			} else {
				break
			}
		} else {
			return 0, errors.New("Error: unexpected end of file")
		}
	}
	for tokens[0].String() == "[" || tokens[0].String() == "map" {
		if tokens[0].String() == "[" {
			if tokens[1].String() == "]" {
				size += 2
				t = append(t, VarType{Type: "slice"})
				tokens = tokens[2:]
			} else if tokens[1].Type() == lex.IntLiteral && tokens[2].String() == "]" {
				size += 3
				t = append(t, VarType{Type: "array", Index: tokens[1].String()})
				tokens = tokens[3:]
			} else {
				return 0, errors.New("Unexpected token")
			}
		} else if tokens[0].String() == "map" {
			if tokens[1].String() == "[" {
				switch tokens[2].String() {
				case "bool", "int", "float", "float32", "float64", "double", "string":
					size += 4
					t = append(t, VarType{Type: "map", Index: tokens[2].String()})
					tokens = tokens[4:]
				default:
					return 0, errors.New("Invalid map type, key must be bool, int, float, float32, float64, double, or string")
				}
			} else {
				return 0, errors.New("Unexpected token")
			}
		}
	}
	t = append(t, VarType{Type: tokens[0].String()})
	if len(tokens) < 1 {
		return 0, errors.New("Incomplete variable")
	}
	m := me.models[len(me.models)-1]
	m.Vars = append(m.Vars, Var{Name: variable, Type: t})
	return size, nil
}

func (me *Parser) parse(tokens []lex.Token) (Out, error) {
	line := 1
	var err error
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.TokType {
		case lex.Whitespace:
			if t.String() == "\n" {
				line++
			} else if t.String() == "\t" {
				var size int
				size, err = me.processVariable(tokens[i:])
				if err != nil {
					return me.out, err
				}
				i += size
			}
		case lex.Identifier:
			me.models = append(me.models, &Model{Name: t.String()})
		default:
			err = errors.New(fmt.Sprintf("Error: world ended on line %d\n       Unexpected token", line+1))
			return me.out, err
		}
	}
	me.models = topSortModels(me.models)
	for _, v := range me.models {
		me.out.addClass(v.Name)
		for _, vv := range v.Vars {
			me.out.addVar(vv.Name, vv.Type)
		}
		me.out.closeClass()
	}
	return me.out, err
}
