---
title: obsidian插件开发教程
tags: [obsidian,plugin,ts,'#Fleeting_N']
date: 2023-06-24 08:20:07
draft: true
hideInList: false
isTop: false
published: false
categories: [obsidian]
---

作者：阿三 
博客：[https://blog.asan123.top](https://blog.asan123.top) 
公众号：阿三爱吃瓜 

 >持续不断记录、整理、分享，让自己和他人一起成长！😊

------
## 将setting用html元素来写

我需要添加多段desc的文本，还必须换行，setDesc 只满足于添加文本，再加一个setDesc也没有用，所以直接截取到需要的部分，然后.descEl创建createDiv就可以了

```ts

//在利用setting这个class进行创建
const block = new Setting(contentEl).setName("Name").setDesc("版本号:1.0.0");
    block.descEl.createDiv({text:"作者:cj"})
    block.descEl.createDiv({text:"描述:这是一段很长很长的描述"})
    block.addButton((btn) =%3E
    btn
      .setButtonText("翻译")
      .setCta()
      .onClick(() => {
        this.close();
        this.onSubmit(this.result);
      }))
    .addButton((btn) =>
      btn
        .setButtonText("还原")
        .setWarning()
        .onClick(() => {
          this.close();
          this.onSubmit(this.result);
        }));


```

```ts

//这段是自行创建的html元素
const setting_item = contentEl.createDiv({cls:"setting-item"});
    setting_item.createDiv({cls:"setting-item-info"})
    .createDiv({text:"Name",cls:"setting-item-name"})
    .createDiv({cls:"setting-item-description"})
    .createDiv({text:"作者:cj"})
    .createDiv({text:"版本号:1.0.0"})
    .createDiv({text:"描述:这是一段很长的描述"});
    const setting_item_control = setting_item.createDiv({cls:"setting-item-control"});
    setting_item_control.createEl("button",{text:"翻译",cls:"mod-cta"});
    setting_item_control.createEl("button",{text:"还原",cls:"mod-warning"});

```

main.ts初始化modal时，如果需要获取到onsubmit结果就在main.ts中追加
```ts

new PluginModal(this.app, (result) => {
				new Notice(`Hello, ${result}!`);
			}).open();

```


settings 用store来设置

最后添加赞赏码的svg格式

golang开发暂时不提


#### **什么是 Markdown**
 
Markdown 是一种方便记忆、书写的纯文本标记语言，用户可以使用这些标记符号以最小的输入代价生成极富表现力的文档：譬如您正在阅读的这份文档。它使用简单的符号标记不同的标题，分割不同的段落，**粗体** 或者 *斜体* 某些文字，更棒的是，**Cmd Markdown** 是我们给出的答案 —— 我们为记录思想和分享知识提供更专业的工具。 您可以使用 Cmd Markdown：
 
> 请保留此份 Cmd Markdown 的欢迎稿兼使用说明，如需撰写新稿件，点击顶部工具栏右侧的 <i class="icon-file"></i> **新文稿** 或者使用快捷键 `Ctrl+Alt+N`。
------
 

### 1. 实时同步预览
 
我们将 Cmd Markdown 的主界面一分为二，左边为**编辑区**，右边为**预览区**，在编辑区的操作会实时地渲染到预览区方便查看最终的版面效果，并且如果你在其中一个区拖动滚动条，我们有一个巧妙的算法把另一个区的滚动条同步到等价的位置，超酷！
 
### 2. 编辑工具栏
 
也许您还是一个 Markdown 语法的新手，在您完全熟悉它之前，我们在 **编辑区** 的顶部放置了一个如下图所示的工具栏，您可以使用鼠标在工具栏上调整格式，不过我们仍旧鼓励你使用键盘标记格式，提高书写的流畅度。
 
![tool-editor](https://www.zybuluo.com/static/img/toolbar-editor.png)
 
### 3. 编辑模式
 
完全心无旁骛的方式编辑文字：点击 **编辑工具栏** 最右侧的拉伸按钮或者按下 `Ctrl + M`，将 Cmd Markdown 切换到独立的编辑模式，这是一个极度简洁的写作环境，所有可能会引起分心的元素都已经被挪除，超清爽！
 
### 4. 实时的云端文稿
 
为了保障数据安全，Cmd Markdown 会将您每一次击键的内容保存至云端，同时在 **编辑工具栏** 的最右侧提示 `已保存` 的字样。无需担心浏览器崩溃，机器掉电或者地震，海啸——在编辑的过程中随时关闭浏览器或者机器，下一次回到 Cmd Markdown 的时候继续写作。
 
### 5. 离线模式
 
在网络环境不稳定的情况下记录文字一样很安全！在您写作的时候，如果电脑突然失去网络连接，Cmd Markdown 会智能切换至离线模式，将您后续键入的文字保存在本地，直到网络恢复再将他们传送至云端，即使在网络恢复前关闭浏览器或者电脑，一样没有问题，等到下次开启 Cmd Markdown 的时候，她会提醒您将离线保存的文字传送至云端。简而言之，我们尽最大的努力保障您文字的安全。
 
### 6. 高亮一段代码[^code]

 
```python
@requires_authorization
class SomeClass:
    pass
if __name__ == '__main__':
    # A comment
    print 'hello world'
```
 

 
 
### 7. 绘制表格
 
| 项目        | 价格   |  数量  |
| --------   | -----:  | :----:  |
| 计算机     | \$1600 |   5     |
| 手机        |   \$12   |   12   |
| 管线        |    \$1    |  234  |
 
### 8. 更详细语法说明
 
 
总而言之，不同于其它 *所见即所得* 的编辑器：你只需使用键盘专注于书写文本内容，就可以生成印刷级的排版格式，省却在键盘和工具栏之间来回切换，调整内容和格式的麻烦。**Markdown 在流畅的书写和印刷级的阅读体验之间找到了平衡。** 目前它已经成为世界上***最大***的技术分享网站 GitHub 和 技术问答网站 StackOverFlow 的御用书写格式。
 
---

 