---
title: obsidian插件开发常用api
tags: [obsidian,'#Fleeting_N']
date: 2023-06-22 14:31:59
draft: true
hideInList: false
isTop: false
published: true
categories: [obsidian]
---

# 3.3 插件开发常用api

作者：阿三 
博客：[https://blog.asan123.top](https://blog.asan123.top) 
公众号：阿三爱吃瓜 

>持续不断记录、整理、分享，让自己和他人一起成长！😊



## vault

### 获取库的基础路径

```ts
import {
  FileSystemAdapter,
} from "obsidian";

const basePath = (
          this.app.vault.adapter as FileSystemAdapter
        ).getBasePath()
```

### 获取全局设置

```ts
const config = this.app.vault.config
```

## workspace

### 获取当前编辑器

```ts
import { MarkdownView } from 'obsidian';

// 获取md编辑器
getEditor() {
    const mdView = this.app.workspace.getActiveViewOfType(MarkdownView);
    if (mdView) {
      return mdView.editor;
    } else {
      return null;
    }
}

```

## editor

### 获取文件内容

```ts
const editor = this.getEditor()
const content = editor.getValue()
```

### 设置文件内容

```ts
const editor = this.getEditor()
// string为需要设置的字符串，采用\n换行
const content = editor.setValue("string")
```

### 获取某一行内容

```ts
const editor = this.getEditor()
// line类型为number
const content = editor.getLine(line)
```

### 设置某一行内容

```ts
const editor = this.getEditor()
// line参数为某一行的行数
const content = editor.setLine(line, "string")
```

### 获取当前文件

```ts
const file = this.app.workspace.getActiveFile();
```

### 获取光标位置

```ts
const editor = this.getEditor()
const cursor = editor.getCursor()
```

### 设置光标

```ts
const editor = this.getEditor()
// line是行数，ch是距离第一个字符串的位置
editor.setCursor({line: 'number', ch: 'number'});
```

### 获取滚动位置

```ts
const editor = this.getEditor()
const { left, top } = editor.getScrollInfo();
```

### 设置滚动位置

```ts
const editor = this.getEditor()
// left,top均为number
editor.scrollTo(left, top);
```

### 获取自定义 frontmatter

```ts
  getFrontmatterValue(key: string, defaultValue: any = undefined) {
    const file = this.app.workspace.getActiveFile();
    if (!file) {
      return undefined;
    }
    const path = file.path;
    const cache = this.app.metadataCache.getCache(path);

    let value = defaultValue;
    if (cache?.frontmatter && cache.frontmatter.hasOwnProperty(key)) {
      value = cache.frontmatter[key];
    }
    return value;
  }
```

### 获取选中的文字

```ts
  getSelectedText(editor: any) {
    if (editor.somethingSelected()) {
      let cursorStart = editor.getCursor(true);
      let cursorEnd = editor.getCursor(false);
      let content = editor.getRange(
        { line: cursorStart.line, ch: 0 },
        { line: cursorEnd.line, ch: editor.getLine(cursorEnd.line).length }
      );
      return {
        start: { line: cursorStart.line, ch: 0 },
        end: {
          line: cursorEnd.line,
          ch: editor.getLine(cursorEnd.line).length,
        },
        content: content,
      };
    } else {
      var lineNr = editor.getCursor().line;
      var contents = editor.getDoc().getLine(lineNr);
      let cursorStart = {
        line: lineNr,
        ch: 0,
      };
      let cursorEnd = {
        line: lineNr,
        ch: contents.length,
      };
      let content = editor.getRange(cursorStart, cursorEnd);
      return { start: cursorStart, end: cursorEnd, content: content };
    }
  }
```

### 发送通知

```ts
import {
  Notice,
} from "obsidian";

// 参数分别为message和timeout
new Notice("message", 1000);
```

## plugin

### 生命周期

```ts
// 挂载时的钩子，一般用于加载设置，设置命令行等初始化操作
async onload(){}

// 卸载插件时的钩子，处理插件中存在的一些副作用，如定时器等
async onunload(){}
```

### 设置快捷键

```ts
this.addCommand({
  id: "id",
  name: "name",
  checkCallback: (checking: boolean) =%3E {
    // do something...
  },
});
```

### 设置文件菜单

```ts
this.app.workspace.on(
  "file-menu",
  (menu: Menu, file: TFile, source: string) => {
    menu.addItem((item: MenuItem) => {
      item.setTitle("title").onClick(() => {
        console.log("dome something...")
      });
    });
  }
);
```

### 加载设置

```ts
async loadSettings() {
  // DEFAULT_SETTINGS 为默认设置
  this.settings = Object.assign(DEFAULT_SETTINGS, await this.loadData());
}
```

### 保存设置

```ts
async saveSettings() {
    await this.saveData(this.settings);
}
```

## 设置页面

### 基础

```ts
class SampleSettingTab extends PluginSettingTab {
	plugin: MyPlugin;

	constructor(app: App, plugin: MyPlugin) {
		super(app, plugin);
		this.plugin = plugin;
	}

	display(): void {
		let {containerEl} = this;

		containerEl.empty();

        // 设置二级标题
		containerEl.createEl('h2', {text: 'Settings for my awesome plugin.'});

		new Setting(containerEl)
			.setName('Setting #1')
			.setDesc('It\'s a secret')
			.addText(text => text // 添加文字输入框
				.setPlaceholder('Enter your secret')
				.setValue('')
				.onChange(async (value) => {
					console.log('Secret: ' + value);
					this.plugin.settings.mySetting = value;
					await this.plugin.saveSettings();
				}));
	}
}
```

### 添加文字输入控件

```ts
new Setting(containerEl)
    .setName('Setting #1')
    .setDesc('It\'s a secret')
    .addText(text => text // 添加文字输入框
        .setPlaceholder('Enter your secret')
        .setValue('') // 默认值
        .onChange(async (value) => {
            console.log('Secret: ' + value);
            this.plugin.settings.mySetting = value;
            await this.plugin.saveSettings();
        }));
```

### 添加开关控件

```ts
new Setting(containerEl)
  .setName("Setting #1")
  .setDesc(
    "It\'s a secret"
  )
  .addToggle(toggle =>
    toggle
      .setValue(this.plugin.settings.mySetting)
      .onChange(async value => {
        this.plugin.settings.mySetting = value;
        await this.plugin.saveSettings();
      })
  );
```

## 添加按钮控件

```ts
    new Setting(containerEl).setName("Setting #1").addButton(cb => {
      cb.setButtonText("buttom");
      cb.setCta();
      cb.onClick(() => {
        console.log("do something")
      });
    });
```

## 其他

### 网络请求

可以安装 `node-fetch` 进行请求