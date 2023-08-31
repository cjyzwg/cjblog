---
title: obsidian多端同步(完美解决)
tags: [obsidian,博客,'#Fleeting_N']
date: 2023-08-28 08:26:47
draft: true
hideInList: false
isTop: false
published: false
categories: [obsidian,博客]
---

作者：阿三 
博客：[Nockygo](https://hexo.hexiefamily.xin) 
公众号：阿三爱吃瓜

> 持续不断记录、整理、分享，让自己和他人一起成长！😊


------
## 背景

看过我视频的小伙伴们应该都知道视频里介绍的同步是remotely-save插件，但是我自己本身不是在用这个，因为同步到坚果云其实会有缺陷。

1、坚果云免费版是有限制的，一旦你笔记很多就会造成这个困扰，可能会超过上限，那就必须得付费了，`建议：图片尽量上传到云上，使用picogo插件这样能减少点容量🤣`。
2、还有个问题是一旦出现某个文件有冲突了，那么就得删除整个文件夹再重新来，即便是有其他的云同步，还是无可避免的出现这个问题。

> 那怎么解决这两个问题呢，我的做法是用： **git**

## 操作

### 安装git

### 开启obsidian-git插件

当然你熟悉git步骤，直接命令行也是可以的，当然用插件更好，可以设置定时更新，也可以设置快捷键，push/pull 操作，当然我这里设置的是快捷键，我不太相信自动同步，总会存在文件diff的问题。

### 导入文件

同之前的视频所讲，用sendanywhere 软件将所有的文件导入到手机上，同时打开即可，如果不懂的童鞋可以看我之前同步的视频。
### ios上使用对应的命令

1、下载ish 命令行工具，安装下列命令
2、添加几个常用的命令脚本，wget即可 命令

```ts

sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

apk update

apk add openssh git nano

gitee上添加文件

ssh-keygen -t ed25519 -C "cj@example.com"

公钥上传到gitee的设置/配置钥匙管理里

添加完：ssh -T git@gitee.com

obsidian ios创建一个新的valut，假设命名为note

mkdir -p /mnt/cj/Obsidian

mount -t ios null /mnt/cj/Obsidian ios会弹出权限，赋予obsidian文件夹权限即可

cd /mnt/cj/Obsidian/note

git clone git@gitee.com:cjyzwg/obsidian_markdown.git markdown

nano commit.sh

./update.sh && cd /mnt/cj/Obsidian/note/markdown && git add . && git commit -m "sync" && git push origin master

chmod a+x commit.sh

nano update.sh

cd /mnt/cj/Obsidian/note/markdown && git pull

chmod a+x update.sh


mac上更新文件了，记得手机上sh update.sh

如果ios上更新文件了，记得手机上sh commit.sh



