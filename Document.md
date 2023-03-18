# 文档

### 数据包类

```go
type DataPack struct {
    Id    int   // 操作序号
    Opt   uint8 // 操作指令
    X, Y  int   // 人物位置
    Mymap [8][8]int // 地图
}
```

用来通讯和保存地图的类。

### 主要函数

### Login

client 端：

```go
func Login(client *rpc.Client)
```

传入一个 `rpc.Client` ， `Login` 函数将会请求输入用户名和密码，并将 用户名 和 密码 使用 `md5`  加密 并和 `version` 组成一个字符串，发送到服务端。

若输入错误，则重复此过程。

server 端：

```go
func (this *McServer) Login(md5str string, reply *DataPack) error
```

接受用户发来的字符串 `md5str` ，并与本地的字符串对比，若相同，令 `reply` 为地图，`reply.Opt = 'A'`。否则 `reply.Opt = 'W'` 

这种登录方式的好处是，即使服务端的数据被盗，入侵者也无法知道 登录者 的密码。

#### FetchClient

```go
func (this *McServer) FetchClient(ClientMap DataPack, reply *DataPack) error
```

服务端的 函数，接受 客户端 发来的数据包，并根据 操作序号 `Id` 和 操作指令 `Opt` 来处理请求。

若 `ClientMap.Opt = 'Q'` ，即请求断开连接，则保存信息，并调用 `CloseConnection` 函数断开连接。

否则 若操作序号小于服务器的操作序号，说明这是过去的信息，不处理。

若操作序号等于服务器的操作序号，则调用 `Compare(data, ClientMap)` 对比数据，若相同 则令 `reply.Opt = 'C' // Correct!` ，否则 `reply.Opt = 'F' // Fault` ，并返回正确的数据包（令 `*reply = data` ）

若操作序号大于服务器的序号，则根据指令移动服务器端的数据，并与客户端的对比。

#### TouchServer

客户端的函数

```go
func TouchServer(client *rpc.Client, cur DataPack, option uint8)
```

参数依次为 `*rpc.client` ，要同步的数据包 和 操作指令。

与服务器同步数据，并等待返回。

#### KeepInTouch

客户端的函数

```go
func KeepInTouch(client *rpc.Client)
```

调用 `TouchServer` ，每隔 0.05s 与服务器端通讯一次。

#### Move

服务器端和客户端都有

```go
func Move(ch uint8)
```

传入移动指令，改变客户端或服务端的 `data` 地图数据。

### 总结与反思

- 你觉得解决这个任务的过程有意思吗？
  
  我觉得解决这个任务的过程非常有意思。

- 你在网上找到了哪些资料供你学习？你觉得去哪里/用什么方式搜索可以比较有效的获得自己想要的资料？
  
  - 参考资料：
    
    - [TCP/IP 介绍 | 菜鸟教程](https://www.runoob.com/tcpip/tcpip-intro.html)
    
    - [GitHub - QSCTech/2023-spring-round-two: 求是潮2023春季纳新二面题仓库](https://github.com/QSCTech/2023-spring-round-two)
    
    - [2023-spring-round-two/Minecraft Sever (Lite Version) at main · QSCTech/2023-spring-round-two · GitHub](https://github.com/QSCTech/2023-spring-round-two/tree/main/Minecraft%20Sever%20(Lite%20Version))
    
    - [Go gob - 简书](https://www.jianshu.com/p/b208cf559c41)
    
    - [Go 语言多维数组 | 菜鸟教程](https://www.runoob.com/go/go-multi-dimensional-arrays.html#:~:text=%E4%BA%8C%E7%BB%B4%E6%95%B0%E7%BB%84%E6%98%AF%E6%9C%80%E7%AE%80%E5%8D%95%E7%9A%84%E5%A4%9A%E7%BB%B4%E6%95%B0%E7%BB%84%EF%BC%8C%E4%BA%8C%E7%BB%B4%E6%95%B0%E7%BB%84%E6%9C%AC%E8%B4%A8%E4%B8%8A%E6%98%AF%E7%94%B1%E4%B8%80%E7%BB%B4%E6%95%B0%E7%BB%84%E7%BB%84%E6%88%90%E7%9A%84%E3%80%82.%20%E4%BA%8C%E7%BB%B4%E6%95%B0%E7%BB%84%E5%AE%9A%E4%B9%89%E6%96%B9%E5%BC%8F%E5%A6%82%E4%B8%8B%EF%BC%9A.%20variable_type%20%E4%B8%BA%20Go%20%E8%AF%AD%E8%A8%80%E7%9A%84%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%EF%BC%8CarrayName%20%E4%B8%BA%E6%95%B0%E7%BB%84%E5%90%8D%EF%BC%8C%E4%BA%8C%E7%BB%B4%E6%95%B0%E7%BB%84%E5%8F%AF%E8%AE%A4%E4%B8%BA%E6%98%AF%E4%B8%80%E4%B8%AA%E8%A1%A8%E6%A0%BC%EF%BC%8Cx%20%E4%B8%BA%E8%A1%8C%EF%BC%8Cy,a%20%5B%20i%20%5D%20%5B%20j%20%5D%20%E6%9D%A5%E8%AE%BF%E9%97%AE%E3%80%82.)
    
    - [网络游戏中服务器端与客户端分别处理哪些事情 - 小 楼 一 夜 听 春 雨 - 博客园](https://www.cnblogs.com/kex1n/archive/2012/05/29/2523992.html)
    
    - 基于 socket 的方法
      
      - https://zhuanlan.zhihu.com/p/143346084
      
      - [go socket编程（详细）_tianlongtc的博客-CSDN博客](https://blog.csdn.net/tianlongtc/article/details/80163661)
      
      - [Go语言基于Socket编写服务器端与客户端通信的实例 - 腾讯云开发者社区-腾讯云](https://cloud.tencent.com/developer/article/1073200)
      
      - [Go语言实现TCP服务端和客户端 - 腾讯云开发者社区-腾讯云](https://cloud.tencent.com/developer/article/1733034)    
    
    - RPC
      
      - https://zhuanlan.zhihu.com/p/139384493
      
      - https://zhuanlan.zhihu.com/p/143961275
    
    - md5加密
      
      - https://zhuanlan.zhihu.com/p/457859814
      
      - [go byte类型-慢慢理解 - 走走停停走走 - 博客园](https://www.cnblogs.com/zccst/p/14054009.html)
    
    - 文件读写
      
      - [Go语言纯文本文件的读写操作](http://c.biancheng.net/view/4556.html)
      
      - https://zhuanlan.zhihu.com/p/259826174
      
      - [Go语言fmt.Printf使用指南 - 简书](https://www.jianshu.com/p/40cbdc02e4b5)
    
    - 结构体
      
      - [GO语音gob包的系列化和反序列化使用和遇到的错误_gob.register_朱鑫烨的博客-CSDN博客](https://blog.csdn.net/weixin_44001557/article/details/102816811?spm=1001.2014.3001.5502) 
      
      - [Go语言采坑记录gob序列化坑_gob: duplicate type received_可爱飞行猪的博客-CSDN博客](https://blog.csdn.net/GeMarK/article/details/89357013)
      
      - [go语言的结构体指针 - 海龟先生 - 博客园](https://www.cnblogs.com/haiguixiansheng/p/10613754.html)
    
    - 锁
      
      - [sync.Mutex互斥锁 - Go语言圣经](https://gopl-zh.github.io/ch9/ch9-02.html)
    
    - 清屏
      
      - [如何在Go中清除终端屏幕？ - 问答 - 腾讯云开发者社区-腾讯云](https://cloud.tencent.com/developer/ask/sof/142335)
      
      - https://rosettacode.org/wiki/Terminal_control/Clear_the_screen#Go
    
    - strconv
      
      - [转换 | strconv (strconv) - Go 中文开发手册 - 开发者手册 - 腾讯云开发者社区-腾讯云](https://cloud.tencent.com/developer/section/1144302)
    
    - time
      
      - [golang sleep函数 休眠延时_go sleep_whatday的博客-CSDN博客](https://blog.csdn.net/whatday/article/details/102637905)
    
    - 常量
      
      - [Go 语言常量 | 菜鸟教程](https://www.runoob.com/go/go-constants.html)
  
  - 你觉得去哪里/用什么方式搜索可以比较有效的获得自己想要的资料？
    
    - 如果是 Go 的基础知识，可以看 Go 语言圣经，菜鸟教程，或者文档。
    
    - 如果是 报错 信息，直接搜索 报错即可。
    
    - 其他的也可以用以上两种方法。

- 在过程中，你遇到最大的困难是什么？你是怎么解决的？
  
  主要是不熟悉这个语言，不知道这个项目该用什么方式实现。有想过自己实现 数据包 的收发，或者 使用 之前了解的 `thrift` 框架实现。
  
  然后在搜 `go thrift` 教程时，发现了 `go`  语言内置的 `rpc` ，实现方式可能比较简单，然后就使用 `rpc` 了。

- 完成任务之后，再回去阅读你写下的代码和文档，有没有看不懂的地方？如果再过一年，你觉得那时你还可以看懂你的代码吗？
  
  - 目前没有看不懂的地方。
  
  - 我觉得一年后我还可以看懂。

- 其他想说的想法或者建议？
  
  - 我想想先

以上。
