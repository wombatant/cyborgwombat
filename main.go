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
	"flag"
	"fmt"
	"github.com/gtalent/lex"
	"io/ioutil"
	"os"
)

type Var struct {
	Name string
	Type []string
}

type Model struct {
	Name string
	Vars []Var
}

func main() {
	out := flag.String("o", "stdout", "File or file set(languages with header files) to write the output to")
	in := flag.String("i", "", "The model file to generate JSON-C code for")
	namespace := flag.String("n", "models", "Namespace for the models")
	version := flag.Bool("v", false, "version")
	flag.Parse()

	if *version {
		fmt.Println("modelmaker version 0.7.3")
		return
	}
	parseFile(*in, *out, *namespace)
}

func parseFile(path, outFile, namespace string) {
	ss, err := ioutil.ReadFile(path)
	if err != nil {
		println("Could not find or open specified model")
		os.Exit(1)
	}
	var tokens []lex.Token

	input := string(ss)
	for input[len(input)-2] == '\n' && input[len(input)-1] == '\n' {
		input = input[:len(input)-1]
	}
	symbols := []string{"[", "]"}
	keywords := []string{}
	stringTypes := []lex.Pair{}
	commentTypes := []lex.Pair{}
	l := lex.NewAnalyzer(symbols, keywords, stringTypes, commentTypes, true)

	for point := 0; point < len(input); {
		var t lex.Token
		t.TokType, t.TokValue, point = l.NextToken(input, point)
		tokens = append(tokens, t)
	}

	var p Parser
	p.out = NewCOut(namespace, USING_JANSSON)
	out, err := p.parse(tokens)
	if err != nil {
		println(err)
		os.Exit(2)
		return
	} else {
		if outFile == "stdout" {
			fmt.Print(out.header(""))
			fmt.Print(out.body(""))
		} else {
			cout := out.(*CppJansson)
			ioutil.WriteFile(outFile+".hpp", []byte(cout.header(outFile+".hpp")), 0644)
			ioutil.WriteFile(outFile+".cpp", []byte(cout.body(outFile+".hpp")), 0644)
			ioutil.WriteFile("modelmakerdefs.hpp", []byte(cout.buildModelmakerDefsHeader()), 0644)
			ioutil.WriteFile("modelmakerdefs.cpp", []byte(cout.buildModelmakerDefsBody()), 0644)
		}
	}
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
			t := vv.Type[len(vv.Type)-1]
			if !isScalar(t) {
				if node, ok := m[t]; ok {
					node.dependents = append(node.dependents, v.model.Name)
					v.remainingDependancies++
				} else {
					println(fmt.Sprintf("Error: unrecognized type: %s", t))
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
		println(fmt.Sprintf("Error: cyclical dependency detected"))
		os.Exit(4)
	}
	return out
}
