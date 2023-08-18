---
title: obsidian 配合 hexo、vercel：让发布博客简单到不可思议
tags: [obsidian,博客发布,'#Fleeting_N']
date: 2023-06-16 15:18:14
draft: true
hideInList: false
isTop: false
published: true
categories: [obsidian]
---

------


![image.png](https://s2.loli.net/2023/03/18/WzmvgDcoRysPAhi.png)

自从我使用obsidian之后，其实越来越少发博客了。原因除了写的少之外，就是我的博客和知识管理系统分属两个地儿，每次都容易忘记去写。

结合我自己的博客发布流程位：
1. obsidian创建文章并写作
2. 利用自己写的半插件使用快捷键`Alt+G`键同步博客文章，以及自己服务器上的博客文章
3. vercel和cloudflare自动拉取github仓库
4. 访问https://blog.asan123.top 和https://hexo.hexiefamily.xin 都可以访问的到
5. 同时博客小程序也可以访问的到

整体的过程就变得非常简单化，我只需要几个快捷键就可以搞定了。

使用过程如下。

### 使用到的 obsidian 插件

#### image auto upload plugin

用于自动上传图片到图床。需要配合 picgo 使用。

#### quickadd

用于快速创建一篇新博客。

下面是我的设置：

1.  创建一个 `Template` 类型的quickadd 命令

![image.png](https://s2.loli.net/2023/03/11/HXaVj2uZneSE9l6.png)

2.  在根目录新建一个 `_Templates` 文件夹，并创建模板文件 `hugo博客模板`

```auto
---
title: {{NAME}}
tags: [{{VALUE:标签？}},'#Fleeting_N']
date: {{DATE:YYYY-MM-DD HH:mm:ss}}
draft: true
hideInList: false
isTop: false
published: false
categories: [{{VALUE:分类？}}]
---
```

published 字段是发布，默认设置为false，草稿。等到文章写完并修改无误后，再修改为 true 进行发布。

4.  设置quickadd 命令：

![image.png](https://s2.loli.net/2023/03/11/Z9BDtVHJr2uaIyq.png)

#### obsidian-git

用于自动备份文件到 github。

插件设置修改如下：

![image.png](https://s2.loli.net/2023/03/11/cxTJiutPEfkHFW8.png)

`注意`：自动备份是个好东西，但是如果手机端也有的话，有可能会造成文件需要合并的问题，所以我是快捷键分开处理。
### 启用 Cloudflare

打开 [Cloudflare Dash](https://dash.cloudflare.com/) 导航栏 `Pages` ，点 `创建项目`，授权 Github 项目，选择博客所在仓库，选择正确的分支。

添加环境变量，指定高版本 `HEXO_VERSION`

然后部署即可。

部署完成后就可以使用 cloudflare 的二级域名访问博客了。如果你像我一样有自己的独立域名，那么可以进行域名绑定。

### Cloudflare 绑定独立域名

首先第一步，把自己的域名托管到 cloudflare。参考： [如何将域名托管到cloudflare](https://www.back2me.cn/skills/cloudflare.html) 这篇文章。

然后打开导航栏 Pages ，在右侧找到刚刚的博客站点，在设置或者部署中找到 【自定义域】，设置自定义域名，输入之前托管进来的域名，按照指引完成绑定。

![image.png](https://s2.loli.net/2023/03/11/mToq84ZpMhFjyGN.png)

以上，所有设置都已完成。
当然用vercel部署是一样的。

现在我要发布一篇新博客时，只要在 obsidian 中打开 hugo博客的这个库，然后使用 quickadd 新建一篇博客，写上内容，然后把yaml 中的 `published` 字段值改为 `true` 即可,✨一定要记得改，否则他就不发布了，（从草稿改为发布）。等待3分钟后 obsidian-git 插件自动同步到 github，博客就自动更新发布好了。

这篇文章就是使用新方式发布的，优雅不是一点点。

### 参考

[Hugo With Obsidian](https://immmmm.com/hugo-with-obsidian/)
[Hi , Cloudflare Pages](https://immmmm.com/hi-cloudflare/)
[Hugo 博客写作最佳实践](https://blog.zhangyingwei.com/posts/2022m4d11h19m42s28/)
[如何将域名托管到cloudflare](https://www.back2me.cn/skills/cloudflare.html)
[hexo部署vercel](https://hexo.io/zh-cn/docs/one-command-deployment#Vercel)
[莉莉蒙的少数派链接](https://sspai.com/u/4b8zstxp/updates)