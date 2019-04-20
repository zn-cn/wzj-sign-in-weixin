新版微助教自动签到

开发环境：

- 语言：Golang:1.11.3
- 借助展现平台：微信公众号

## 功能

- 多用户自动签到（二维码签到，GPS签到，普通签到均可）
- 支持自定义坐标，如：114.440465,30.517877
- 支持自定义坐标标签，如：东十二
- 输入一次openid或者含有openid的链接之后，两小时内可自动签到所有签到
- 获取自己设置的所有标签和坐标

效果如下：

![screenshot1](./dist/screenshot1.png)

![screenshot7](./dist/screenshot7.png)

## 使用

关注公众号：阿楠的公众号

![qrcode](./dist/qrcode.jpg)

使用说明：

四种状态：

- 0 -> 默认状态
  - 输入0也会将状态重置
- 1 -> 设置为输入链接状态
  - 输入1之后直接输入微助教任意页面复制后的链接（带openid）
  - 或者直接输入openid
  - 注：此openid仅有两小时有效时间，因此每次上课前复制一次在公众号输入即可，每次进入微助教网页openid均会更新，如果重新进入了微助教页面，请更新openid
- 2 -> 设置为输入坐标状态
  - 输入格式如下：`东十二:114.440465,30.517877` ，左侧为标签，右侧为经纬度
- 3 -> 设置当前坐标标签状态
  - 输入格式如下：`东十二`
  - 设置之后之后所有的签到将使用此标签代表的坐标进行签到
- 4 -> 获取设置的所有标签和坐标

示例如下：

![screenshot5](./dist/screenshot5.jpg)

![screenshot6](./dist/screenshot6.jpg)

## 原理

签到原理：

- 通过拉取https://v18.teachermate.cn/wechat/wechat/guide/signin?openid=$openid 页面，如下：

  ![screenshot3](./dist/screenshot3.png)

- 注意上图中hidden的input组件，其中包含了所需签到的所有信息，如：签到id：sign-id，课程id：course-id, 临时openid: openid

- 利用HTML分析包，提取出其中信息，若没有相关信息则暂停此次签到

- 使用接口：https://v18.teachermate.cn/wechat-api/v1/class-attendance/student-sign-in 进行签到

  form data: 

  ```json
  {
      "openid":"6d10c8781d8b1dc19120915880bce8a1",
  	"lon":114.440465,
  	"lat":30.517877,
  	"courseId":1098194,
  	"signId":1514058,
  	"wx_csrf_name":"5cd27ada5a6bdb83bd0fdb23fab1abcf"
  }

  ```

  Header:

  ```json
  {
  	"User-Agent":"Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 MicroMessenger/6.2.3 NetType/WIFI Language/zh_CN",
  	"Content-Type":"application/x-www-form-urlencoded; charset=UTF-8"
  }

  ```

公众号自动签到原理：

- 设置一个定时器，如：every 5s, 定时检测队列中是否存在监控对象
- 取出队列中监控的对象，监控对象可存储openid（微助教openid）
- 对象监控自动过期删除（两小时）
- 遍历监控对象进行签到


## TODO

- [ ] 帮助文档完善
- [ ] 日志记录
- [ ] 。。。