package main

import (
	"errors"
	"fmt"
	"github.com/gtalent/lex"
	"io/ioutil"
	"os"
)

type Out struct {
	txt   string
	types []string
}

func NewCOut() Out {
	var out Out
	out.txt = `//Generated Code
#include <string>
#include "types.h"

using namespace std;
`
	out.types = []string{}
	return out
}


func (me *Out) addVar(v string, t string) {
	me.txt += "\t" + t + " " + v + ";\n"
}

func (me *Out) addClass(v string) {
	me.txt += "\nclass " + v + " {\n"
}

func (me *Out) closeClass() {
	me.txt += "};\n"
}

func (me *Out) endsWithClose() bool {
	return (me.txt)[len(me.txt)-2:] != "};"
}

func main() {
	getTokens("image.txt")
}

func getTokens(path string) []lex.Token {
	ss, err := ioutil.ReadFile(path)
	if err != nil {
		println("Could not find or open sconscript")
		os.Exit(0)
	}
	var tokens []lex.Token

	input := string(ss)
	symbols := []string{}
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
		return nil
	}
	fmt.Print(out.txt)
	return tokens
}

type Parser struct {
	out Out
}

func (me *Parser) processVariable(tokens []lex.Token) error {
	if len(tokens) < 4 {
		return errors.New("Incomplete variable")
	}
	me.out.addVar(tokens[1].String(), tokens[3].String())
	return nil
}

func (me *Parser) processObject(tokens []lex.Token) (Out, error) {
	line := 0
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
				err = me.processVariable(tokens[i:])
				i += 3
			}
		case lex.Identifier:
			me.out.addClass(t.String())
		default:
			err = errors.New("Unidentified token")
		}
		prev = t.String()
		if err != nil {
			println(fmt.Sprintf("Error: World ended on line %d.\n", line))
			break
		}
	}
	if me.out.endsWithClose() {
		me.out.closeClass()
	}
	return me.out, err
}
