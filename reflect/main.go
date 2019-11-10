package main

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	ABC string
}

func main() {

	//list := []MyStruct{{ABC:"aflscdac "}, {"cececece"}}
	//Fetch(list, "MyStruct")
}

func Fetch(args []interface{}, name string) []string {
	for _, v := range args {
		value := reflect.TypeOf(v)
		fmt.Printf("type: %s", value)
	}
	return nil
}
