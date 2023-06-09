---
title: 搭建git服务器超简单命令
categories:
  - 树莓派
---
### 搭建git私有服务器命令
#### 服务器：
1.sudo apt-get install git-core  
2.git --bare init /home/git/rep.git (rep为仓库的名字)
#### 本地：
1.git clone git@10.8.0.8:/home/git/rep.git
