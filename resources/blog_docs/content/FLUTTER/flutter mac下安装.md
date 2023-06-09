---
categories:
  - FLUTTER
---
#### 参考资料
**https://ninghao.net/video/6364**  
**https://flutterchina.club/setup-macos/**

##### 启动命令：
>$ flutter create cjios  
>$ flutter emulator --launch apple_ios_simulator  
如果出现no device 首先要装android studio  

##### 如果出现UIApplication.LaunchOptionsKey问题
解决办法：

nano  cjios/ios/Runner/AppDelegate.swift
[UIApplication.LaunchOptionsKey: Any]  
改成  
[UIApplicationLaunchOptionsKey: Any]

##### 运行命令  
>$ flutter run   
##### 安装命令
>$ flutter devices  
iPhone • xxxxxx • ios • iOS 13.2
>$ flutter run -d   xxxxxx 
>即可安装到手机上了

#### 工具栏：
1、https://pub.dev/packages/zefyr  
2、https://javiercbk.github.io/json_to_dart/
3、获取列表数据：https://blog.csdn.net/weixin_30867015/article/details/98277080
4、展示markdown文件：https://github.com/flutter/flutter_markdown/blob/master/example/lib/main.dart