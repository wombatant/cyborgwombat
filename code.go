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
	"strings"
)

type CppCode struct {
	tabs string
	text string
}

func (me *CppCode) String() string {
	return me.text
}

func (me *CppCode) Insert(text string) {
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		v := lines[i]
		me.text += "\n" + me.tabs + v
	}
}

func (me *CppCode) PushBlock() {
	me.text += "\n" + me.tabs + "{"
	me.tabs += "\t"
}

func (me *CppCode) PushIfBlock(cond string) {
	me.text += "\n" + me.tabs + "if (" + cond + ") {"
	me.tabs += "\t"
}

func (me *CppCode) Else() {
	me.text += "\n" + me.tabs[:len(me.tabs)-1] + "} else {"
}

func (me *CppCode) PushForBlock(cond string) {
	me.text += "\n" + me.tabs + "for (" + cond + ") {"
	me.tabs += "\t"
}

func (me *CppCode) PopBlock() {
	me.tabs = me.tabs[:len(me.tabs)-1]
	me.text += "\n" + me.tabs + "}"
}
