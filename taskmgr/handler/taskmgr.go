package handler

import (
	"github.com/panjf2000/ants"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"fmt"
)

var (
	gTaskMgr *TaskMgr
	gTaskId int
)

type pTaskFunc func(interface{})

type TaskItem struct {
	taskId int
	arg   interface{}
	pFunc pTaskFunc
}

type TaskMgr struct {
	p *ants.PoolWithFunc
	size int
	taskQueue []TaskItem
	lock sync.Mutex
}


func GTaskMgr() *TaskMgr {
	if gTaskMgr == nil {
		gTaskMgr = newTaskMgr(10)
	}
	return gTaskMgr
}

func newTaskMgr(size int) *TaskMgr {
	return &TaskMgr{
		size: size,
	}
}

//添加任务
func (t *TaskMgr) AddTask(arg interface{}, pFunc pTaskFunc) error {

	gTaskId++

	t.lock.Lock()
	item := TaskItem{
		taskId: gTaskId,
		arg: arg,
		pFunc: pFunc,
	}
	t.taskQueue = append(t.taskQueue, item)
	fmt.Printf("AddTask %d\n", gTaskId)
	t.lock.Unlock()

	//err := t.p.Serve(&TaskItem{
	//	arg:   arg,
	//	pFunc: pFunc,
	//})
	//if err != nil {
	//	logrus.Errorf("TaskMgr AddTask error: %s", err.Error())
	//	return err
	//}
	return nil
}

// 分派任务处理
func (t *TaskMgr) HandlerFunc(arg interface{}) {
	task, ok := arg.(*TaskItem)
	if ok {
		if task.pFunc != nil {
			task.pFunc(task.arg)
		}
	} else {
		logrus.Error("TaskItem error")
	}
}

func (t *TaskMgr) run(exit chan bool) {
	p, err := ants.NewPoolWithFunc(t.size, func(arg interface{}) {
		t.HandlerFunc(arg)
	})
	if err != nil {
		panic("TaskMgr NewPoolWithFunc error")
	}
	defer p.Release()

	t.p = p

Loop:
	for {
		select{
		case <-exit:
			fmt.Println("TaskMgr run exit!!!")
			break Loop
		case <-time.After(time.Duration(500) * time.Microsecond):
			fmt.Printf("running: %d free: %d\n", p.Running(), p.Free())
			//添加任务
			if p.Free() > 0 {
				var item *TaskItem
				t.lock.Lock()
				if len(t.taskQueue) > 0 {
					item = &t.taskQueue[0]
					t.taskQueue = t.taskQueue[1:]
				} else {
					item = nil
				}
				t.lock.Unlock()

				if item != nil {
					err = p.Serve(item)
					if err != nil {
						fmt.Println("TaskAdd fail")
					} else {
						fmt.Printf("TaskAdd suc %d\n", item.taskId)
					}
				}
			}
		}
	}


}

func RunTaskMgrServer() error {

	fmt.Println("TaskMgr running!")
	taskMgr := GTaskMgr()
	ex := make(chan bool)
	go taskMgr.run(ex)


	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-ch:
		fmt.Println("TaskMgr exit!")
	}

	close(ex)

	return nil
}