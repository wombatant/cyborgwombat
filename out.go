package main

type Out interface {
	header(string) string
	body(string) string
	addClass(string)
	addVar(string, []string)
	closeClass()
}
