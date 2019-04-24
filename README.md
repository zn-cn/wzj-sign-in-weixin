新版微助教自动签到

开发环境：

- 语言：Golang:1.11.3
- 借助展现平台：微信公众号

## 功能

- 多用户自动签到（二维码签到，GPS签到，普通签到均可）

- - 输入一次openid或者含有openid的链接之后，两小时内可自动签到所有签到
  - 签到成功邮件提醒

- 支持自定义坐标，如：114.440465,30.517877

- 支持自定义坐标标签，如：东十二

- 获取自己设置的所有标签和坐标

- 课程讨论开启邮件提醒

- 获取所有课程

- 设置讨论课程

- 设置通知邮件

注：

- 输入的带有openid的链接有效期为两小时

- 重新进入学生栏目的网页会使之前的链接失效

- - 这时需要重新进公众号输入带有openid的链接 

- 考虑到在微助教那边设置讨论课程要进入微助教网页会更新openid，这里开发一个设置接口

- 设置讨论课程之后需要到**微助教服务**号参与讨论

- **邮件提醒可能在垃圾邮件里面！**

- 自动签到间隔为20 s，防止排名太靠前被发现了

## 使用指南

关注公众号：阿楠技术

![qrcode](./dist/qrcode.jpg)

![img](https://mmbiz.qpic.cn/mmbiz_png/o0CwfJSKW4xiacguVibaxCVDHupDicXysSESJicpy5OlGCFMK98iafzKHofY9ficcJJznahzErF9lkbXficTJ3mK3UbHQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

![img](https://mmbiz.qpic.cn/mmbiz_jpg/o0CwfJSKW4xiacguVibaxCVDHupDicXysSESuQygic7Of8LR3rxfenssHuYlFBRibM1RNgw3ScLMicYJ6c0YgiaIHQibCg/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

下图操作手速要快，不然链接出不来的！

![img](https://mmbiz.qpic.cn/mmbiz_jpg/o0CwfJSKW4xiacguVibaxCVDHupDicXysSE3ER03JzFyxQsp7fmMOu95IqrbqVkEoZ0YyI3sIJMr9Z7qqCL1bXvxA/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

![img](https://mmbiz.qpic.cn/mmbiz_jpg/o0CwfJSKW4xiacguVibaxCVDHupDicXysSEp8lLHF9T7s2SnVQL8088W2HkF0JGxAWGjr4KSIODMSwVibFz1nhAAfQ/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

坐标拾取：https://lbs.amap.com/console/show/picker

常用经纬度

- 东十二:114.440465,30.517877
  - F楼:114.440799,30.517263
- 东九
  - A区:114.433653,30.519743
  - B区:114.433165,30.519453
  - C区:114.432995,30.519014
  - D区:114.433286,30.518387
- 西十二
  - 北门:114.413702,30.514882
  - 南门:114.413811,30.514386
  - 东门:114.414505,30.514723
  - 西门:114.413069,30.514554

效果如下：

![screenshot1](./dist/screenshot1.png)

![screenshot7](./dist/screenshot7.png)

![img](https://mmbiz.qpic.cn/mmbiz_jpg/o0CwfJSKW4xiacguVibaxCVDHupDicXysSEFwEmBHaAZ28OwCJEFyxJllHY94uvY47U3rvHEJGspKwEpuLiaYquJXA/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

![screenshot11](./dist/screenshot11.jpg)

## 原理

### 签到原理

1. 首先自己注册一个微助教教师账号，开启签到
2. 打开微助教签到页面，我们能发现签到API如下

![img](https://mmbiz.qpic.cn/mmbiz_png/o0CwfJSKW4xiacguVibaxCVDHupDicXysSETc7M9tZu9bmLRbx2K7U5AcVB05qJxq5rAlyIYx8hnIcasiaF9Yrcgicw/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

3. 对整个网页文件搜索 course-id 或者 sign-id，发现如下页面

![img](https://mmbiz.qpic.cn/mmbiz_png/o0CwfJSKW4xiacguVibaxCVDHupDicXysSEgD1MiadsyyYHBrtibZ6VjczjfDMH6KZveoX3zEBjXiaqyohyEMJnL5E2A/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

4. 开发签到自动化接口流程如下

- 工具用户输入的openid拼接签到页面链接

- 下载签到页面链接，提取 hidden 的 input 组件信息 

- 使用接口进行签到，发送数据如下

  Form Data

![img](https://mmbiz.qpic.cn/mmbiz_png/o0CwfJSKW4xiacguVibaxCVDHupDicXysSEKMxdPEntxdVUV5sWwG9WgDJDC4IvkgmibOjAd398VtFsibPQQ2crmBicQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

​	Header

![img](https://mmbiz.qpic.cn/mmbiz_png/o0CwfJSKW4xiacguVibaxCVDHupDicXysSEJqicAy9iaNwu4fJbZNSExXG7ew9x6tBhpibmZkFwtM3ibOObFsHzehubAw/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

注：sign-id 为签到id，下图为微助教教师签到页面

![img](https://mmbiz.qpic.cn/mmbiz_png/o0CwfJSKW4xiacguVibaxCVDHupDicXysSE4h54ibewCr4KBzp635lrYgibOv0uUeLVbnHqCZAELRrHcICKgqFlZ40A/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

### 微信后台开发

介绍：微信公众号提供开发者模式，使得用户在微信公众号发送的消息都会被转发到开发者的服务器，开发者可立即回复用户，同时开发者可在48小时内调用客服接口给用户发送其他消息。更多详情请查看微信公众号开发文档。

项目开发技术：

- 使用 MongoDB 持久化存储用户信息
- 使用 Redis 存储临时监控队列，临时状态信息等
- 开发语言：Golang

自动签到原理：

1. 定时器

   - 每 5s 对监控队列对象进行一次检测
   - 检测是否有签到开启，有则进行签到
   - 检测是否有讨论开启，有则进行提醒

2. 监控队列实现

   - 用户输入 openid 时，在 Redis 设置 key/value，存储openid，并设置为两小时过期
   - key格式：user:task:$openid, value: openid
   - 注：第一个openid为开发者公众号给用户的一个 openid，第二个 openid 为用户输入的微助教的 openid
   - 每次检测时，取出 Redis 命名空间 user:task 中的 key，进行相应检测

3. 签到数据处理

   坐标来源：用户在公众号设置坐标标签时，后台会将坐标存储到 Redis 中，用户签到时从 Redis 中取出。

   坐标处理：签到时对坐标进行随机化处理和截断处理，防止坐标惊人的一致，同时与微助教签到数据格式

   代码如下：

   ```go
   // 随机化处理，防止一致
   coordinate.Lon += float64(rand.Intn(40)-20) * 0.000001
   coordinate.Lat += float64(rand.Intn(40)-20) * 0.000001
   data.Set("lon", strconv.FormatFloat(coordinate.Lon, 'f', 5, 64)) // 5 表示截断为5位小数
   data.Set("lat", strconv.FormatFloat(coordinate.Lat, 'f', 5, 64))
   ```

   注：图中 5 表示截断为 5 位小树

更多细节请阅读项目源码

## 注意

- 本项目仅供学习和个人使用
- 侵删