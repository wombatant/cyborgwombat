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
	"flag"
	"fmt"
	"github.com/gtalent/lex"
	"io/ioutil"
	"os"
)

func main() {
	out := flag.String("o", "stdout", "File or file set(languages with header files) to write the output to")
	in := flag.String("i", "", "The model file to generate JSON-C code for")
	version := flag.Bool("v", false, "version")
	flag.Parse()

	if *version {
		fmt.Println("modelmaker version 0.1")
		return
	}
	parseFile(*in, *out)
}

func parseFile(path, outFile string) {
	ss, err := ioutil.ReadFile(path)
	if err != nil {
		println("Could not find or open specified model")
		os.Exit(0)
	}
	var tokens []lex.Token

	input := string(ss)
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
	p.out = NewCOut()
	out, err := p.processObject(tokens)
	if err != nil {
		return
	} else {
		if outFile == "stdout" {
			fmt.Print(out.header())
			fmt.Print(out.body(""))
		} else {
			ioutil.WriteFile(outFile+".hpp", []byte(out.header()), 0644)
			ioutil.WriteFile(outFile+".cpp", []byte(out.body(outFile+".hpp")), 0644)
			ioutil.WriteFile("modelmakerdefs.hpp", []byte(out.buildModelmakerDefsHeader()), 0644)
			ioutil.WriteFile("modelmakerdefs.cpp", []byte(out.buildModelmakerDefsBody()), 0644)
		}
	}
}

type Parser struct {
	out Out
}

func (me *Parser) processVariable(tokens []lex.Token) (int, error) {
	size := 3 // should be 1 less than the actual number parsed

	if len(tokens) < 4 {
		return 0, errors.New("Incomplete variable")
	}

	variable := tokens[1].String()
	index := 0
	tokens = tokens[3:]
	for tokens[0].String() == "[" {
		if tokens[1].String() != "]" {
			return 0, errors.New("] expected")
		}
		size += 2
		tokens = tokens[2:]
		index++
	}
	t := tokens[0].String()
	if len(tokens) < 1 {
		return 0, errors.New("Incomplete variable")
	}
	me.out.addVar(variable, t, index)
	return size, nil
}

func (me *Parser) processObject(tokens []lex.Token) (Out, error) {
	line := 1
	var err error
	prev := ""
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.TokType {
		case lex.Whitespace:
			if t.String() == "\n" {
				line++
				if prev == "\n" {
					me.out.closeClass()
				}
			} else if t.String() == "\t" {
				var size int
				size, err = me.processVariable(tokens[i:])
				i += size
			}
		case lex.Identifier:
			me.out.addClass(t.String())
		default:
			err = errors.New("Unidentified token")
		}
		prev = t.String()
		if err != nil {
			println(fmt.Sprintf("Error: world ended on line %d\n", line))
			break
		}
	}
	if me.out.endsWithClose() {
		me.out.closeClass()
	}
	return me.out, err
}
