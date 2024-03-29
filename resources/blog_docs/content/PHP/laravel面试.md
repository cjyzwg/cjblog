---
title: laravel面试
date: 2020-12-10 10:44:17
categories:
  - PHP
---

### 一、背景
最近去面试，去了一家游戏公司面试，由于公司不用框架，并且没有深入理解laravel，导致面试吃亏：
* 1.解释下laravel的服务容器，控制反转IOC？
* 2.说下laravel的队列？
* 3.之前公司是如何解耦的？


### 二、答案
* 1.
依赖注入:通过构造注入，函数调用或者属性的设置来提供组件的依赖关系
一开始我们代码依赖关系可能是如图这样的:
![avatar](https://blog.hexiefamily.xin/assets/ioc1.jpg)
由于不可抗拒的原因，PHP版本升级，需求更改等等，要把α对象替换掉，把β对象删掉
![avatar](https://blog.hexiefamily.xin/assets/ioc2.jpg)

IOC容器在全局维持一个对象实例集合和类名集合，我们在写某个类的时候把这个类依赖的对象注册到容器里，调用这个类的时候再实例化拿出来。
这个就是IOC的思想，一个系统通过组织控制和对象的完全分离来实现”控制反转”。对于依赖注入，这就意味着通过在系统的其他地方控制和实例化依赖对象，从而实现了解耦。  

* 2.
Laravel可配置多种队列驱动，包括 "sync", "database", "beanstalkd", "sqs", "redis", "null"
比如向用户发送邮件的场景：现在有10w封邮件需要发送，最简单的，我们需要有一个方法将邮件的收件人、内容等，拆分成10w条任务放在队列中，同时需要设置一个回调方法负责处理每条任务。当队列中有邮件发送任务时，队列会主动调用回调方法，并传递任务详情进去。回调方法处理完成后，单条邮件即发送完毕。其他邮件依样处理。
![avatar](https://blog.hexiefamily.xin/assets/ioc3.jpg)
![avatar](https://blog.hexiefamily.xin/assets/ioc4.jpg)


##### 结论
* 好好总结