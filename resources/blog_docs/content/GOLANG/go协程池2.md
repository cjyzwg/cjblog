---
title:go协程池2
categories:
  - GOLANG
---
#### go协程池改编
-基于[<font color=#0099ff>go协程池</font>](http://blog.hexiefamily.xin/article?path=%2fGOLANG%2fgo%e5%8d%8f%e7%a8%8b%e6%b1%a0.md)  

#### 看下协程池的原理
![avatar](https://blog.hexiefamily.xin/assets/workerpool2.png)  

#### 协程池结构体,并且建立有缓存的任务channel
```go
type Taskpool struct {
	Tasknum     int
	Taskchannel chan *Task
}

func NewTaskPool(tasknum int, taskmaxlen int) *Taskpool {
	tp := &Taskpool{
		Tasknum:     tasknum,
		Taskchannel: make(chan *Task, taskmaxlen),
	}
	return tp
}
func (tp *Taskpool) AllocateTask() {
	defer close(tp.Taskchannel)
	for i := 0; i < tp.Tasknum; i++ {
		t := NewTask(func() error {
			fmt.Println("task ", i, " is allocate")
			return nil
		}, i)
		tp.Taskchannel <- t
	}
}
```
#### 图中goroutine2
```go
func (tp *Taskpool) ReceiveTask(jobchannel chan *Task) {
	defer close(jobchannel)
	for task := range tp.Taskchannel {
		fmt.Println("job channel get this task channel", task.Id)
		jobchannel <- task
	}
}
```

#### 图中goroutine3
```go
func (wp *Workerpool) ReceiveWork(jobchannel chan *Task) {
	defer close(wp.Workchannel)
	for work := range jobchannel {
		fmt.Println("work channel get this work from job channel", work.Id)
		wp.Workchannel <- work
	}
}
```

#### 完整代码如下
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	F  func() error
	Id int
}

func NewTask(argc func() error, id int) *Task {
	t := &Task{
		F:  argc,
		Id: id,
	}
	return t
}

func (t *Task) Execute() {
	fmt.Println("task ", t.Id, " is running ")
	t.F()
	time.Sleep(2 * time.Second)
}

//Taskpool 可以不包含channel不然只剩tasknum
type Taskpool struct {
	Tasknum     int
	Taskchannel chan *Task
}

func NewTaskPool(tasknum int, taskmaxlen int) *Taskpool {
	tp := &Taskpool{
		Tasknum:     tasknum,
		Taskchannel: make(chan *Task, taskmaxlen),
	}
	return tp
}
func (tp *Taskpool) AllocateTask() {
	defer close(tp.Taskchannel)
	for i := 0; i < tp.Tasknum; i++ {
		t := NewTask(func() error {
			fmt.Println("task ", i, " is allocate")
			return nil
		}, i)
		tp.Taskchannel <- t
	}
}

func (tp *Taskpool) ReceiveTask(jobchannel chan *Task) {
	defer close(jobchannel)
	for task := range tp.Taskchannel {
		fmt.Println("job channel get this task channel", task.Id)
		jobchannel <- task
	}
}

type Workerpool struct {
	Workernum   int
	Workchannel chan *Task
}

func (wp *Workerpool) ReceiveWork(jobchannel chan *Task) {
	defer close(wp.Workchannel)
	for work := range jobchannel {
		fmt.Println("work channel get this work from job channel", work.Id)
		wp.Workchannel <- work
	}
}
func NewWorkerPool(workernum int, workmaxlen int) *Workerpool {
	w := &Workerpool{
		Workernum:   workernum,
		Workchannel: make(chan *Task, workmaxlen),
	}
	return w
}
func (w *Workerpool) CreateWork() {
	var wg sync.WaitGroup
	fmt.Println("work num is:", w.Workernum)
	for i := 0; i < w.Workernum; i++ {
		wg.Add(1)
		//waitgroup 需要加goroutine
		go w.GetWork(&wg, i)
	}
	wg.Wait()
}
func (w *Workerpool) GetWork(wg *sync.WaitGroup, id int) {
	for job := range w.Workchannel {
		job.Execute()
		fmt.Println("worker ", id, " works with task ", job.Id)
	}
	wg.Done()
}
func main() {
	startTime := time.Now()
	channelLen := 10
	workerNum := 20
	taskNum := 100
	jobChannel := make(chan *Task, channelLen)
	worker := NewWorkerPool(workerNum, channelLen)
	task := NewTaskPool(taskNum, channelLen)
	go task.AllocateTask()
	go task.ReceiveTask(jobChannel)
	go worker.ReceiveWork(jobChannel)
	worker.CreateWork()
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time is :", diff.Seconds(), " seconds")
}

```
