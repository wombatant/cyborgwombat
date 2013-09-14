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
package parser

import (
	"fmt"
	"github.com/gtalent/lex"
	"os"
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

/*
  Return values:
    * number of tokens parsed
    * error
*/
func processVariable(model *Model, tokens []lex.Token) (int, error) {
	size := 2 // should be 1 less than the actual number parsed

	if len(tokens) < 4 {
		return 0, fmt.Errorf("Error: incomplete variable")
	}

	var t []VarType
	var variable string
	if len(tokens) > 2 {
		variable = tokens[1].String()
		tokens = tokens[2:]
	} else {
		return 0, fmt.Errorf("Error: unexpected end of file")
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
			return 0, fmt.Errorf("Error: unexpected end of file")
		}
	}
	for tokens[0].String() == "[" || tokens[0].String() == "map" {
		if tokens[0].String() == "[" {
			if len(tokens) > 2 && tokens[1].String() == "]" {
				size += 2
				t = append(t, VarType{Type: "slice"})
				tokens = tokens[2:]
			} else if len(tokens) > 3 && tokens[1].Type() == lex.IntLiteral && tokens[2].String() == "]" {
				size += 3
				t = append(t, VarType{Type: "array", Index: tokens[1].String()})
				tokens = tokens[3:]
			} else {
				if len(tokens) > 2 {
					return 0, fmt.Errorf("Error: unexpected token")
				} else {
					return 0, fmt.Errorf("Error: unexpected end of file")
				}
			}
		} else if tokens[0].String() == "map" {
			if len(tokens) > 4 && tokens[1].String() == "[" {
				switch tokens[2].String() {
				case "bool", "int", "float", "float32", "float64", "double", "string":
					size += 4
					t = append(t, VarType{Type: "map", Index: tokens[2].String()})
					if tokens[3].String() != "]" {
						return 0, fmt.Errorf("Error: expected token \"]\" atfter map index type \"%s\", got \"%s\"", tokens[2].String(), tokens[3].String())
					}
					if tokens[4].Type() != lex.Identifier {
						return 0, fmt.Errorf("Error: expected type after token \"]\", got \"%s\"", tokens[4].String())
					}
					tokens = tokens[4:]
				default:
					return 0, fmt.Errorf("Error: invalid map type, key must be bool, int, float, float32, float64, double, or string")
				}
			} else {
				if len(tokens) > 4 {
					return 0, fmt.Errorf("Error: unexpected token")
				} else {
					return 0, fmt.Errorf("Error: unexpected end of file")
				}
			}
		}
	}
	t = append(t, VarType{Type: tokens[0].String()})
	if len(tokens) < 1 {
		return 0, fmt.Errorf("Error: incomplete variable")
	}
	model.Vars = append(model.Vars, Var{Name: variable, Type: t})
	return size, nil
}

/*
  Parses a set of model specs into a list of model models.
*/
func Parse(input string) ([]*Model, error) {
	//parse into tokens
	var tokens []lex.Token

	//if len(input) < 2 {
	//	return
	//}
	//for input[len(input)-2] == '\n' && input[len(input)-1] == '\n' {
	//	input = input[:len(input)-1]
	//}
	symbols := []string{"[", "]", "#"}
	keywords := []string{}
	stringTypes := []lex.Pair{}
	commentTypes := []lex.Pair{}
	l := lex.NewAnalyzer(symbols, keywords, stringTypes, commentTypes, true)

	for point := 0; point < len(input); {
		var t lex.Token
		t.TokType, t.TokValue, point = l.NextToken(input, point)
		tokens = append(tokens, t)
	}

	line := 1
	var models []*Model
	var err error
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.TokType {
		case lex.Symbol:
			if t.String() == "#" {
				for ; i < len(tokens); i++ {
					if tokens[i].String() == "\n" {
						break
					}
				}
			}
		case lex.Whitespace:
			if t.String() == "\n" {
				line++
			} else if t.String() == "\t" {
				size := 1
				if len(tokens) > 1 && tokens[i+1].Type() == lex.Identifier {
					size, err = processVariable(models[len(models)-1], tokens[i:])
					if err != nil {
						return models, fmt.Errorf("Error on line %d: \n\t%s", line, err)
					}
				}
				i += size
			}
		case lex.Identifier:
			models = append(models, &Model{Name: t.String()})
		default:
			err = fmt.Errorf("Error on line %d: \n\tError: unexpected token", line+1)
			return models, err
		}
	}
	models = topSortModels(models)
	return models, err
}

func isScalar(v string) bool {
	switch v {
	case "bool", "int", "double", "float32", "float64", "string", "unknown":
		return true
	}
	return false
}

/*
  Topicologically sorts models to be sure they are declared
  in a workable order.
*/
func topSortModels(models []*Model) []*Model {
	type topSortNode struct {
		model                 *Model
		remainingDependancies int
		//indices of if dependents
		dependents []string
	}

	out := make([]*Model, len(models))
	m := make(map[string]*topSortNode)
	a := make([]*topSortNode, len(models))
	//build name map
	for i, v := range models {
		node := new(topSortNode)
		node.model = v
		a[i] = node
		m[v.Name] = node
	}
	//build dependency structure
	for _, v := range a {
		for _, vv := range v.model.Vars {
			t := vv.Type[len(vv.Type)-1].Type
			if !isScalar(t) {
				if node, ok := m[t]; ok {
					node.dependents = append(node.dependents, v.model.Name)
					v.remainingDependancies++
				} else {
					fmt.Println(fmt.Sprintf("Error: unrecognized type: %s", t))
					os.Exit(3)
				}
			}
		}
	}

	index := 0
	cyclicalDeps := false
	//sort
	for len(a) != 0 && !cyclicalDeps {
		cyclicalDeps = true
		for i := 0; i < len(a); i++ {
			v := a[i]
			if v.remainingDependancies < 1 {
				cyclicalDeps = false
				out[index] = v.model
				index++
				a[i] = a[len(a)-1]
				a = a[:len(a)-1]
				for _, vv := range v.dependents {
					m[vv].remainingDependancies--
				}
			}
		}
	}
	if cyclicalDeps {
		fmt.Println(fmt.Sprintf("Error: cyclical dependency detected"))
		os.Exit(4)
	}
	return out
}
