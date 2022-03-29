# rmlibs-gopool
> 协程池
### 设计思路
池化协程，任务执行完成后进行回收
```
//指定容量、过期时间，初始化协程池
demoPool := gopool.NewPool(size int, expired time.Duration)
```
有待执行任务时，看池中是否有空闲的执行者worker<br/>
&emsp;&emsp;如果有则取出一个来执行，完成后回收<br/>
&emsp;&emsp;如果池中没有空闲worker，则判断worker数量是否达到容量上限<br/>
&emsp;&emsp;&emsp;&emsp;如果没有达到上限，则新建一个worker来执行任务<br/>
&emsp;&emsp;&emsp;&emsp;否则阻塞等待空闲的worker来执行
```
//task为func任务，params为func的参数
demoPool.Do(task func,params ...interface{})
```
### 执行流程
![avatar](https://cdn.nlark.com/yuque/0/2021/svg/22090258/1637564285530-50537a70-0bec-4892-8a7f-27ef1adae338.svg)