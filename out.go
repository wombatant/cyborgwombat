/*
   Copyright 2013-2014 gtalent2@gmail.com

   This Source Code Form is subject to the terms of the Mozilla Public
   License, v. 2.0. If a copy of the MPL was not distributed with this
   file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package main

import (
	"github.com/gtalent/cyborgbear/parser"
)

const (
	USING_JANSSON = iota
	USING_QT      = iota
	USING_GO      = iota
)

type Out interface {
	write(string) string
	writeFile(string) error
	addClass(string)
	addVar(string, []parser.VarType)
	closeClass(string)
}
