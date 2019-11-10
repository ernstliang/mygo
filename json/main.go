package main

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
)

type RecItem struct {
	Count  int     `json:"Count"`
	Module int     `json:"Module"`
	Ids    []int64 `json:"Ids"`
}

type Category struct {
	Primary   string `json:"Primary"`
	Secondary string `json:"Secondary"`
}

func main() {

	categories := []Category{
		{Primary: "150951177368911141", Secondary: "150858136633228584"},
		{"101", "201"},
	}
	b, err := ffjson.Marshal(&categories)
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(b))

	var categories2 []Category
	err = ffjson.Unmarshal(b, &categories2)
	if err != nil {
		panic(err)
	}
	fmt.Println(categories2)


	//测试int64的json格式化
	ids := []int64{150951177368911141, 150858136633228584, 150952798366743848}
	var r RecItem
	r.Count = 10
	r.Module = 1
	r.Ids = ids

	buf, err := ffjson.Marshal(&r)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf)

	var r2 RecItem
	err = ffjson.Unmarshal(buf, &r2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("r2: %+v", r2)
}