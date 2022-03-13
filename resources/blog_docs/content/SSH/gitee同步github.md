## 背景
最近因为在做公众号自动同步的开发，每次一push代码，服务器的webhook就会不生效，同时会停止运行，导致博客网站停止，必须得手动取开启。

### 原因
初步判断是阿里云在git pull github的代码时，一直不执行，可能是外网的问题。

### 解决思路
想到之前国内的gitee也存在代码，考虑同步github,再通过webhook 更新


![](https://gitee.com/cjyzwg/img/raw/master/202203131218367.png)

### 步骤

#### **实时同步github到gitee**
1.ssh-keygen -t rsa -C "youremail@example.com"
生成的 id_rsa 是私钥，id_rsa.pub 是公钥。(注意此处不要设置密码，生成的公私钥用于下面 GitHub / Gitee 的配置，以保证公私钥成对，否则从 GitHub -%3E Gitee 的同步将会失败。)

![](https://gitee.com/cjyzwg/img/raw/master/202203131223319.png)

2.在 GitHub 项目的「**Settings** -> **Secrets** → **New repository secret**」路径下配置好命名为 **GITEE_RSA_PRIVATE_KEY**，value值填写 **id_rsa 私钥**的内容。

3.在 GitHub 的个人设置页面「**Settings **-> **SSH and GPG keys**」配置 SSH Keys 公钥（即：**id_rsa.pub**），命名随意。

![](https://gitee.com/cjyzwg/img/raw/master/202203131227576.png)

4.在 Gitee 的个人设置页面「**安全设置** -> **SSH 公钥**」配置 SSH 公钥（即：**id_rsa.pub**），命名随意。

![](https://gitee.com/cjyzwg/img/raw/master/202203131228907.png)

5.创建 workflow

```yaml
name: SyncToGitee

on:
  push:
    branches: [master]
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Github Sync to Gitee
        uses: wearerequired/git-mirror-action@master
        env:
          # 注意在 Settings->Secrets 配置 GITEE_RSA_PRIVATE_KEY
          SSH_PRIVATE_KEY: ${{ secrets.GITEE_RSA_PRIVATE_KEY }}
        with:
          # 注意替换为你的 GitHub 源仓库地址
          source-repo: git@github.com:cjyzwg/cjblog.git
          # 注意替换为你的 Gitee 目标仓库地址
          destination-repo: git@gitee.com:cjyzwg/cjblog.git
          
```

6.执行同步,修改代码（如修改 README），提交，成功触发同步！

![](https://gitee.com/cjyzwg/img/raw/master/202203131231561.png)

#### **配置gitee webhook**
1.在 Gitee 项目的「**管理** -> **WebHooks** 」路径下配置好url和webhook密码。

2.在 Github 项目的「**Settings** -> **Webhooks** 」路径下删除对应的webhook。

**如果克隆的地址是github，则要执行此步**

3.在对应的执行的项目目录下切换成gitee：

> git remote set-url origin https://gitee.com/cjyzwg/cjblog.git

##### 同步操作结束，愉快的执行go run main.go吧～