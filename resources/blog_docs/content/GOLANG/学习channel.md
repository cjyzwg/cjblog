---
title: 学习channel
date: 2022-08-17 21:55:22
categories:
  - GOLANG
---
```go
package main
import(
	"fmt"
	"time"
)
func producer(nums ...int) <-chan int{// 括号里输入参数，括号外返回的参数类型
	out :=make(chan int)
	go func() {
		defer close(out)
		for _,n :=range nums{
			out<-n
		}
	}()
	return out
}
func cal(inch <-chan int) <-chan int{//括号里从chan读取数据，括号外返回的参数类型
	out :=make(chan int)
	go func() {
		defer close(out)
		for n :=range inch{
			out<-n*n
		}
	}()
	return out
}
func main(){
	t := time.Now()
	in :=producer(1000000)
	ch :=cal(in)
	//consumer
	for ret :=range ch{
		fmt.Printf("%3d", ret)
	}
	fmt.Println()
	elapsed := time.Since(t)
    fmt.Println("app elapsed:", elapsed)
}
```

##### 这个是首尾各一管道，读取数据用协程读取管道，写入数据用协程写入管道  
1.无缓冲的  就是一个送信人去你家门口送信 ，你不在家 他不走，你一定要接下信，他才会走。  
2.无缓冲保证信能到你手上
3.有缓冲的 就是一个送信人去你家仍到你家的信箱 转身就走 ，除非你的信箱满了 他必须等信箱空下来。
4.有缓冲的 保证 信能进你家的邮箱
5.如果这里没有缓冲的话，就会出现阻塞，因为无缓冲，他需要等你接到信，main 方法执行完了，你还没接，阻塞了，针对协程，要加个有缓冲，这样直接到邮箱里，到时候自己去取，这里的缓冲池，在retch中。
```go
package main
import(
	"fmt"
	"time"
)
func workerpool(n int,jobch <-chan int,retch chan<- string){// 没有返回的参数，jobch 读取chan中的数据，retch 写数据到chan中
	for i :=0;i<n;i++{
		go worker(i,jobch,retch);
	}
}
func worker(id int,jobch <-chan int,retch chan<- string){
	for job :=range jobch{
		ret := fmt.Sprintf("worker %d process job %d",id,job)
		retch <-ret
	}
}
func genjob(n int) <-chan int{
	jobch :=make(chan int,200)
	go func(){
		for i :=0;i<n;i++{
			jobch <-i
		}
		close(jobch)
	}()
	return jobch
}
func main(){
	jobch :=genjob(10)
	retch :=make(chan string,200)
	workerpool(5,jobch,retch)
	time.Sleep(time.Second)
	close(retch)
	for ret :=range retch{
		fmt.Println(ret)
	}
}
```






