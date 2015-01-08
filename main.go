/*
   Copyright 2013-2014 gtalent2@gmail.com

   This Source Code Form is subject to the terms of the Mozilla Public
   License, v. 2.0. If a copy of the MPL was not distributed with this
   file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/gtalent/cyborgbear/parser"
	"io/ioutil"
	"os"
)

type parseFileArgs struct {
	out        string
	in         string
	namespace  string
	outputType string
	include    string
	lowerCase  bool
	version    bool
}

func main() {
	var args parseFileArgs
	flag.StringVar(&args.out, "o", "stdout", "File or file set(languages with header files) to write the output to")
	flag.StringVar(&args.in, "i", "", "The model file to generate JSON serialization code for")
	flag.StringVar(&args.namespace, "n", "models", "Namespace for the models")
	flag.StringVar(&args.outputType, "t", "cpp-jansson", "Output type(cpp-jansson, cpp-qt)")
	flag.StringVar(&args.include, "include", "", "header file to include")
	flag.BoolVar(&args.lowerCase, "lc", false, "Make variable names lowercase in output models")
	flag.BoolVar(&args.version, "v", false, "version")
	flag.Parse()

	if args.version {
		fmt.Println("cyborgwombat version " + cyborgbear_version)
		return
	}
	parseFile(args)
}

func parseFile(args parseFileArgs) {
	ss, err := ioutil.ReadFile(args.in)
	if err != nil {
		fmt.Println("Could not find or open specified model file")
		os.Exit(1)
	}
	input := string(ss)

	ioutputType := USING_JANSSON
	switch args.outputType {
	case "cpp-jansson":
		ioutputType = USING_JANSSON
	case "cpp-qt":
		ioutputType = USING_QT
	}

	var out Out
	switch ioutputType {
	case USING_JANSSON, USING_QT:
		out = NewCOut(args.namespace, ioutputType, args.lowerCase)
	}

	models, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
		return
	} else {
		for _, v := range models {
			out.addClass(v.Name)
			for _, vv := range v.Vars {
				out.addVar(vv.Name, vv.Type)
			}
			out.closeClass(v.Name)
		}

		if args.out == "stdout" {
			fmt.Print(out.write(""))
		} else {
			out.writeFile(args.out)
		}
	}
}
