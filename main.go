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
	boost      bool
	lowerCase  bool
	version    bool
}

func main() {
	var args parseFileArgs
	flag.StringVar(&args.out       , "o", "stdout", "File or file set(languages with header files) to write the output to")
	flag.StringVar(&args.in        , "i", "", "The model file to generate JSON serialization code for")
	flag.StringVar(&args.namespace , "n", "models", "Namespace for the models")
	flag.StringVar(&args.outputType, "t", "cpp-jansson", "Output type(cpp-jansson, cpp-qt, go)")
	flag.StringVar(&args.include   , "include", "", "header file to include")
	flag.BoolVar(&args.boost    , "cpp-boost", false, "Boost serialization enabled")
	flag.BoolVar(&args.lowerCase, "lc", false, "Make variable names lowercase in output models")
	flag.BoolVar(&args.version  , "v", false, "version")
	flag.Parse()

	if args.version {
		fmt.Println("cyborgbear version " + cyborgbear_version)
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
	case "go", "Go":
		ioutputType = USING_GO
	}

	var out Out
	switch ioutputType {
	case USING_JANSSON, USING_QT:
		out = NewCOut(args.namespace, ioutputType, args.boost, args.lowerCase)
	case USING_GO:
		out = NewGo(args.namespace)
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
