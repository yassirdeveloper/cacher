package main

var GetCommand = NewCommand("GET").
	AddArgument(0, "string").
	AddOption('e', "expires-in", "int")
