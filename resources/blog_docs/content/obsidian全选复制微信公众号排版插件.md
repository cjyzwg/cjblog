---
title: obsidian全选复制微信公众号排版插件
tags: [obsidian,博客，微信公众号,'#Fleeting_N']
date: 2023-08-23 23:58:51
draft: true
hideInList: false
isTop: false
published: false
categories: [obsidian,博客,微信公众号]
---

作者：阿三 
博客：[Nockygo](https://hexo.hexiefamily.xin) 
公众号：阿三爱吃瓜

> 持续不断记录、整理、分享，让自己和他人一起成长！😊


------
## 前言

平时写博客，直接markdown渲染到页面上，直接排版即可，因为有一个公众号，想着可以直接将obsidian写的文章也能够放到公众号上，这样能省去了排版的部分，以前写过一篇关于无痛苦更新公众号的文章，用的是[golang写的工具版本](https://github.com/cjyzwg/markdown-wechat)。

**80%的时间用于写作，而不是用在排版**

### 全选复制

写作的流程：

> obsidian写作 -> obsidian插件 -> 全选复制样式 -> 更新公众号文章

### 获取方式

`注意是付费版本哦` ，在公众号聊天界面回复  `排版` ，或者网页打开[面包多链接](https://mbd.pub/o/bread/mbd-ZJ2WlZxy)即可下载
#### 配置下默认主题

如果你想要设置默认css，请将你的css文件放到插件目录的themes目录下，同时配置这里，修改为你的css文件名，`注意：.css必须存在`

![](https://weimgpub.oss-cn-hangzhou.aliyuncs.com/img/202308241659931.png)

当obsidian写完，点击`全选复制`，打开微信公众号的草稿页面，粘贴就可以用啦～

下面是各个主题的样式：

![](https://weimgpub.oss-cn-hangzhou.aliyuncs.com/img/202308220143273.png)

![](https://weimgpub.oss-cn-hangzhou.aliyuncs.com/img/202308241840913.png)

![](https://weimgpub.oss-cn-hangzhou.aliyuncs.com/img/202308241840665.png)


#### css文件配置的定义

```css
/* 1~6 标题样式定义 */
h1 {} h2 {} h3 {} h4 {} h5 {} h6 {}
a { color: red; } /* 超链接样式定义 */
strong {} /* 加粗样式定义 */
del {} /* 删除线样式定义 */
em {}  /* 下划线样式定义 */
u {}   /* 下划线样式定义 */
p {}   /* 段落样式定义 */
ul {}  /* 无序列表样式定义 */
ol {}  /* 有序列表样式定义 */
li {}  /* 列表条目样式定义 */
blockquote {} /* 块级引用样式定义 */
table {}
td {}
th {}
pre {}        /* 样式定义 */
.code-highlight {} /* 代码块样式定义 */
.code-line {}      /* 代码块行样式定义 */
.code-spans {}     /* 代码块行样式定义 */

sup {} /* GFM 脚注样式定义 */
.footnotes-title {} /* GFM 脚注，参考标题样式定义 */
.footnotes-list {} /* GFM 脚注，参考列表样式定义 */

.image-warpper {} /* 图片父节点样式定义 */
.image {} /* 图片样式定义 */

/* 部分代码高亮样式 */
.comment {}
.property {}
.function {}
.keyword {}
.punctuation {}
.unit {}
.tag {}
.color {}
.selector {}
.quote {}
.number {}
.attr-name {}
.attr-value {}
```
### 总结

全选复制到剪贴板不仅仅支持微信公众号，其他内容编辑器比如面包多也是支持的。

我其实还开发了第三个版本的插件，也是我现在最常用的，ob用什么主题，那么上传到公众号里也是什么主题，同时自动添加封面图片+原文作者，只需要在手机订阅号里预览即可实现一键发布，这个下期再说～


