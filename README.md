# pixiv-api

## 进度
开发中。。。

## 参考
本项目的实现细节参考：https://github.com/upbit/pixivpy

## 概述
这是一个pixiv的api，通过go编写。

## 功能
- 搜图
- 推荐
- ...

## 目前的问题
现在没办法通过username和password登录，必须先获取refreshToken，通过token调用p站的http api。

## 文档&变量说明
变量有一些约定如下：
- 所有的常量或变量取单词的第一个字母大写，如果有多个选项，就在后面加0 1 2 3 ...
  如，SearchTarget这个参数是个枚举，那么它的option就是ST0、ST1、ST2、ST3...