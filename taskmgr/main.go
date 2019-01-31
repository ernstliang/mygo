package main

import (
	"golang.org/x/sync/errgroup"
	"github.com/ernstliang/mygo/taskmgr/handler"
)

func main() {

	var (
		err error
		g errgroup.Group
	)

	g.Go(handler.RunTaskMgrServer)

	g.Go(handler.RunTaskAddServer)

	if err = g.Wait(); err != nil {
		panic(err)
	}
}