package handler

import (
	"fmt"
	"math/rand"
	"time"
)

type AddParam struct {
	A int
	B int
}

func HandlerAddParam(arg interface{}) {
	param, ok := arg.(*AddParam)
	if !ok {
		fmt.Errorf("param error")
		return
	}

	c := param.A + param.B
	fmt.Printf("Add %d + %d = %d\n", param.A, param.B, c)
	time.Sleep(time.Duration(1)*time.Second)
	fmt.Println("Add finish")
}

func RunTaskAddServer() error {
	time.Sleep(1*time.Second)

	for i := 0; i < 100; i++ {
		arg := &AddParam{
			A: rand.Int()%100,
			B: rand.Int()%100,
		}
		GTaskMgr().AddTask(arg, HandlerAddParam)
	}

	return nil
}