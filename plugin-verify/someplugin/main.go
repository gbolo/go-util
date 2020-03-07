package main

import "fmt"

const version = "R13_2"

type usercc string

func (u usercc) Version() string {
	return version
}

func (u usercc) DoSomething() {
	fmt.Println("Doing Something...")
}

var SomePlugin usercc
