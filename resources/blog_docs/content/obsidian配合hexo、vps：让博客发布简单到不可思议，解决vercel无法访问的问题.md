---
title: obsidian配合hexo、vps：让博客发布简单到不可思议，解决vercel无法访问的问题
tags: [obsidian,博客发布,'#Fleeting_N']
date: 2023-07-16 09:08:25
draft: true
hideInList: false
isTop: false
published: true
categories: [obsidian,博客]
---

- 作者：阿三
- 博客：[https://hexo.hexiefamily.xin](https://hexo.hexiefamily.xin)
- 公众号：阿三爱吃瓜

> 持续不断记录、整理、分享，让自己和他人一起成长！😊

## 用vercel出现了墙的问题？
 
我用vercel，cloudfare搭建hexo博客只需要轻松地上传到github上，过个十秒钟左右就可以访问到博客最新的内容了，但不知道什么原因，vercel、cloudfare等其他类似厂家都被墙了，之前修改过DNS记录，临时解决这个问题，但这一段时间又无法访问了，虽然博客没啥访问量，但是不能不访问啊！！！
 
## 想到的解决思路

### 使用国内的类似平台

的确有，国内是[zeabur](https://zeabur.com/)，尝试过，用法差不多，如果不翻墙的话，访问的速度依旧是差强人意。

### 使用云服务器

前几年白嫖了阿里云服务器，就一直在用，本身也有[原始的博客](https://blog.hexiefamily.xin/)在用，主要给小程序提供接口来着，直接就在这台服务器重新部署hexo博客以及docksify文档说明。

#### 遇到了新问题

- 服务器是ubuntu 16.04版本，node安装以及npm安装就遇到`node: /lib/x86_64-linux-gnu/libc.so.6: version GLIBC_2.28 not found (required by node)`问题，解决了但没完全解决。
- 重新换思路，直接使用docker compose安装，版本还是1.18的老版本，算了，不想大动了，就用这个版本。

> **前提还是用obsidian写文章->github->服务器->访问，直接从服务器开搞**

## docker版本开始

我在github上找到了别人写的一个hexo docker compose 版本，因为我要兼容docsify，所以小小的修改了下，源码请[访问](https://github.com/cjyzwg ，后续如果把docsify改成vuepress也是类似的做法，照常修改即可。

### 1、修改docker-composer.yml

```yml

version: "2"
services:
  hexo-svc:
    build:
      context: .
      args:
        NODE_VERSION: 14
    container_name: cj-hexo
    image: hexo-docsify:latest
    restart: always
    environment:
      HEXO_SERVER_PORT: 4000
    ports:
      - '4000:4000'
    volumes:
      - ./data/blog:/root/blog
  docsify-svc:
    build:
      context: .
      args:
        NODE_VERSION: 14
    container_name: cj-docsify
    image: hexo-docsify:latest
    restart: always
    environment:
      DOCSIFY_SERVER_PORT: 3000
    ports:
      - '3000:3000'
    volumes:
      - ./data/docs:/root/docs
    entrypoint: ["docsify-entrypoint.sh"]
    
```

这里entrypoint里的sh脚本，需要同时在Dockerfile中声明到`/usr/local/bin`目录下，否则启动容器的时候会报错，这样docker-compose 启动会覆盖对应的entrypoint的脚本值。

```Dockerfile

COPY hexo-entrypoint.sh /usr/local/bin/
RUN chmod a+x /usr/local/bin/hexo-entrypoint.sh

COPY docsify-entrypoint.sh /usr/local/bin/
RUN chmod a+x /usr/local/bin/docsify-entrypoint.sh

# 容器启动的时候默认执行hexo的命令
ENTRYPOINT ["hexo-entrypoint.sh"]

```

### 2、shell文件夹下push_hexo.sh和push_docsify.sh修改对应的target文件夹

```push_hexo.sh

#!/bin/bash
target=/root/docker-hexo-docsify

```

### 3、启动docker compose

```sh

docker-compose up -d 

```

### 4、用nginx反向代理

```conf

server {
    listen 80;
    listen 443;
    server_name hexo.hexiefamily.xin;
    ssl_certificate /etc/nginx/ssl_cert/hexo.hexiefamily.xin/cert.pem;
    ssl_certificate_key /etc/nginx/ssl_cert/hexo.hexiefamily.xin/key.pem;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_prefer_server_ciphers on;
    location / {
        proxy_set_header HOST $host;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        proxy_pass http://127.0.0.1:4000/;
    }
}

```

`注意`：这里我加了ssl，ssl证书我用的是acme自动生成letsencript证书，大概每三个月过期一次，我本地已经写好对应的脚本，只要更新即可，这里就不延展开了，本地在`/Documents/code/shell/alidns/update.sh`中，有需要的关注公众号：`阿三爱吃瓜`，联系我即可。

### 5、本地联同obsidian运行的脚本

obsidian开了个插件，可详细[查看或者获取](https://docs.hexiefamily.xin/#/md/obsidian/docs/c-3%E6%8F%92%E4%BB%B6Api%E7%A4%BA%E4%BE%8B)，同时本地对应的脚本如下：

``` docs.sh

#!/bin/sh
vercelpath=/Users/cj/Documents/code/hexo/vercel_doc
indexpath=/md/idea-plugin
function add_git(){
    git status
    git pull origin master
    git add .
    git commit -m "更新文档"
    git push origin master
}
cd $vercelpath
add_git

# 更新服务器，不使用hook，安全性太低了
echo "服务器开始更新"
output=`ssh cjsrrfamily "bash /root/docker_hexo_docsify/shell/push_docsify.sh"`
echo $output


```

插件快捷键是`Alt/Option+g`即可同步到github上，间隔几十秒会同步到服务器上。


## 总结

原先使用obsidian配合vercel快速部署发布博客，很流畅，但是由于墙的问题解决就很费劲，所以该用上传到github的同时上传到服务器上，后续如果有时间改的话，需要改成ations，同时我这里也没有采用github的webhook模式，感觉安全系数太低了，所以没采用，还有涉及到github ssh public key的问题必应下就找到了，所以这个整体大致步骤是满足我的日常需求的，如果你还有任何帮助，请联系我。
