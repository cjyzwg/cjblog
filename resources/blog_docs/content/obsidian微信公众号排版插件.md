---
title: obsidian微信公众号排版插件
tags: [obsidian,博客，微信公众号,'#Fleeting_N']
date: 2023-08-21 23:58:51
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

平时写博客，直接markdown渲染到页面上，直接排版即可，由于也有一个公众号，想着可以直接将obsidian写的文章也能够放到公众号上，这样能省去了排版的部分，以前写过一篇关于无痛苦更新公众号的文章，**想要做到80%的时间用于写作，而不是用在排版**，现在有三个版本的无痛苦更新公众号。

### 第一个版本

免费开源地址在：[go插件版本](https://github.com/cjyzwg/markdown-wechat)，写作的流程：

> obsidian写作 -> go插件 -> 更新公众号文章

![](https://mmbiz.qpic.cn/mmbiz_png/H7fyhCeib6kpI4AtkNGglyEURERiaqr5mG8rLEv6tQOfkIiaSko2a3vfzCtDjLLUAN11XFInUjtyppt557LAk9lqA/640?wxfrom=5&wx_lazy=1&wx_co=1)

1、下载release版本或者关注【**阿三爱吃瓜**】公众号回复更新公众号即可获得

2、公众号需要开启白名单，开启基本配置获取appid和secret

### 第二个版本

`付费版本`，写作的流程：

> obsidian写作 -> obsidian插件 -> 全选复制样式 -> 更新公众号文章

![](https://weimgpub.oss-cn-hangzhou.aliyuncs.com/img/202308220143273.png)

选择相应的主题，自己定义css，放入插件的themes目录下，css要求：

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

虽然全选复制到剪贴板中很方便，但也有着很明显的缺点：

1、由于obsidian的自带或第三方主题是需要魔改才能用的，很多都是不能直接用的，因为是直接渲染来着，而且很多是带有自定义变量，微信公众号是不识别的，所以全选复制只能是自己重新定义css才可以用。

2、每次必须打开微信公众号草稿箱复制一遍到文本编辑器中，其实也很麻烦。

基于这两点，我其实开发了第三个版本的插件，也是我现在用的，ob用什么主题，那么上传到公众号里也是什么主题，我这里使用的是`typora-vue主题`同时自动添加封面图片+原文作者，只需要在手机订阅号里预览即可实现一键发布。当然了引出第三个版本啦～

### 第三个版本

`付费版本`，写作的流程：

> obsidian写作 -> obsidian插件 -> 一键更新公众号文章

1、公众号需要开启白名单，开启基本配置获取appid和secret，在插件的设置页面上配置相应的信息

2、无需打开公众号页面，只用打开手机订阅助手即可实现一键发布

3、当然微信公众号也有直接发布，我没有开发，希望是先检查下有没有问题再发布～

`注意：如果出现不匹配，说明obsidian自带主题的css会过于复杂，解析出现错误，我这里使用的是typora-vue主题`

| 版本       | 开源版 | 全选复制版 | 一键发布版 |
| ---------- | ------ | ---------- | ---------- |
| 价格       | 免费   | 免费      | ¥49.9      |
| 主题更换   | ❌     | ✅         | ✅         |
| 开启白名单 | ✅     | ❌         | ✅         |
| 公众号配置  | ✅    | ❌         | ✅         |
| 方便程度    | 一般   | 简单       | 最简单     |
| 自动配置封面图片 |  ❌   | ❌        | ✅         |
| 自动配置原文作者 |  ❌   | ❌        | ✅         |


 