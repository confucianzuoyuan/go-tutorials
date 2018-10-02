# 1, 实现的 API 功能

## 实现的功能

本教程通过实现一个账号系统，来演示如何构建一个真实的 API 服务器。通过实战展示了 API 构建过程中各个流程（准备 -> 设计 -> 开发 -> 测试 -> 部署）的实现方法，内容如下：

[](./images/neirong.png)

详细为：

- 设计阶段
    - API 构建技术选型
    - API 基本原理
    - API 规范设计
- 开发阶段
    - 如何读取配置文件
    - 如何管理和记录日志
    - 如何做数据库的 CURD 操作
    - 如何自定义错误 Code
    - 如何读取和返回 HTTP 请求
    - 如何进行业务逻辑开发
    - 如何对请求插入自己的处理逻辑
    - 如何进行 API 身份验证
    - 如何进行 HTTPS 加密
    - 如何用 Makefile 管理 API 源码
    - 如何给 API 命令添加版本功能
    - 如何管理 API 命令
    - 如何生成 Swagger 在线文档
- 测试阶段
    - 如何进行单元测试
    - 如何进行性能测试（函数性能）
    - 如何做性能分析
    - API 性能测试和调优
- 部署阶段
    - 如何用 Nginx 部署 API 服务
    - 如何做 API 高可用

通过以上各功能的介绍，我们可以完整、系统地学习 API 构建方法和技巧。

## 账号系统(apiserver)业务功能

- API 服务器状态检查
- 登录用户
- 新增用户
- 删除用户
- 更新用户
- 获取指定用户的详细信息
- 获取用户列表

# 2, RESTful API 介绍

## 什么是 API

API（Application Programming Interface，应用程序编程接口）是一些预先定义的函数或者接口，目的是提供应用程序与开发人员基于某软件或硬件得以访问一组例程的能力，而又无须访问源码，或理解内部工作机制的细节。

要实现一个 API 服务器，首先要考虑两个方面：API 风格和媒体类型。Go 语言中常用的 API 风格是 RPC 和 REST，常用的媒体类型是 JSON、XML 和 Protobuf。在 Go API 开发中常用的组合是 gRPC + Protobuf 和 REST + JSON。

## REST 简介

REST 代表表现层状态转移（REpresentational State Transfer），由 Roy Fielding 在他的 论文 中提出。REST 是一种软件架构风格，不是技术框架，REST 有一系列规范，满足这些规范的 API 均可称为 RESTful API。REST 规范中有如下几个核心：

- 1, REST 中一切实体都被抽象成资源，每个资源有一个唯一的标识 —— URI，所有的行为都应该是在资源上的 CRUD 操作
- 2, 使用标准的方法来更改资源的状态，常见的操作有：资源的增删改查操作
- 3, 无状态：这里的无状态是指每个 RESTful API 请求都包含了所有足够完成本次操作的信息，服务器端无须保持 Session

>无状态对于服务端的弹性扩容是很重要的。

REST 风格虽然适用于很多传输协议，但在实际开发中，REST 由于天生和 HTTP 协议相辅相成，因此 HTTP 协议已经成了实现 RESTful API 事实上的标准。在 HTTP 协议中通过 POST、DELETE、PUT、GET 方法来对应 REST 资源的增、删、改、查操作，具体的对应关系如下：

| HTTP | 方法 | 行为 | URI | 示例说明 |
|------|-----|------|-----|--------|
| GET  | 获取资源列表 | /users | 获取用户列表 |
| GET  | 获取一个具体的资源 | /users/admin | 获取 admin 用户的详细信息 |
| POST | 创建一个新的资源 | /users | 创建一个新用户 |
| PUT  | 以整体的方式更新一个资源 | /users/1 | 更新 id 为 1 的用户 |
| DELETE | 删除服务器上的一个资源 | /users/1 | 删除 id 为 1 的用户 |

## RPC 简介

根据维基百科的定义：远程过程调用（Remote Procedure Call，RPC）是一个计算机通信协议。该协议允许运行于一台计算机的程序调用另一台计算机的子程序，而程序员无须额外地为这个交互作用编程。

通俗来讲，就是服务端实现了一个函数，客户端使用 RPC 框架提供的接口，调用这个函数的实现，并获取返回值。RPC 屏蔽了底层的网络通信细节，使得开发人员无须关注网络编程的细节，而将更多的时间和精力放在业务逻辑本身的实现上，从而提高开发效率。

RPC 的调用过程如下：

[](./images/rpc.png)

- 1, Client 通过本地调用，调用 Client Stub
- 2, Client Stub 将参数打包（也叫 Marshalling）成一个消息，然后发送这个消息
- 3, Client 所在的 OS 将消息发送给 Server
- 4, Server 端接收到消息后，将消息传递给 Server Stub
- 5, Server Stub 将消息解包（也叫 Unmarshalling）得到参数
- 6, Server Stub 调用服务端的子程序（函数），处理完后，将最终结果按照相反的步骤返回给 Client

>Stub 负责调用参数和返回值的流化（serialization）、参数的打包解包，以及负责网络层的通信。Client 端一般叫 Stub，Server 端一般叫 Skeleton。

## REST vs RPC

在做 API 服务器开发时，很多人都会遇到这个问题 —— 选择 REST 还是 RPC。RPC 相比 REST 的优点主要有 3 点：

- 1, RPC+Protobuf 采用的是 TCP 做传输协议，REST 直接使用 HTTP 做应用层协议，这种区别导致 REST 在调用性能上会比 RPC+Protobuf 低
- 2, RPC 不像 REST 那样，每一个操作都要抽象成对资源的增删改查，在实际开发中，有很多操作很难抽象成资源，比如登录操作。所以在实际开发中并不能严格按照 REST 规范来写 API，RPC 就不存在这个问题
- 3, RPC 屏蔽网络细节、易用，和本地调用类似

>这里的易用指的是调用方式上的易用性。在做 RPC 开发时，开发过程很烦琐，需要先写一个 DSL 描述文件，然后用代码生成器生成各种语言代码，当描述文件有更改时，必须重新定义和编译，维护性差。

但是 REST 相较 RPC 也有很多优势：

- 轻量级，简单易用，维护性和扩展性都比较好
- REST 相对更规范，更标准，更通用，无论哪种语言都支持 HTTP 协议，可以对接外部很多系统，只要满足 HTTP 调用即可，更适合对外，RPC 会有语言限制，不同语言的 RPC 调用起来很麻烦
- JSON 格式可读性更强，开发调试都很方便
- 在开发过程中，如果严格按照 REST 规范来写 API，API 看起来更清晰，更容易被大家理解

>在实际开发中，严格按照 REST 规范来写很难，只能尽可能 RESTful 化。

其实业界普遍采用的做法是，内部系统之间调用用 RPC，对外用 REST，因为内部系统之间可能调用很频繁，需要 RPC 的高性能支撑。对外用 REST 更易理解，更通用些。当然以现有的服务器性能，如果两个系统间调用不是特别频繁，对性能要求不是非常高，以笔者的开发经验来看，REST 的性能完全可以满足。本小册不是讨论微服务，所以不存在微服务之间的高频调用场景，此外 REST 在实际开发中，能够满足绝大部分的需求场景，所以 RPC 的性能优势可以忽略，相反基于 REST 的其他优势，笔者更倾向于用 REST 来构建 API 服务器，本小册正是用 REST 风格来构建 API 的。

## 媒体类型选择

媒体类型是独立于平台的类型，设计用于分布式系统间的通信，媒体类型用于传递信息，一个正式的规范定义了这些信息应该如何表示。HTTP 的 REST 能够提供多种不同的响应形式，常见的是 XML 和 JSON。JSON 无论从形式上还是使用方法上都更简单。相比 XML，JSON 的内容更加紧凑，数据展现形式直观易懂，开发测试都非常方便，所以在媒体类型选择上，选择了 JSON 格式，这也是很多大公司所采用的格式。

# 3, API 流程和代码结构

为了使大家在开始实战之前对 API 开发有个整体的了解，这里选择了两个流程来介绍：

- HTTP API 服务器启动流程
- HTTP 请求处理流程

## HTTP API 服务器启动流程

[](./images/httpstart.png)

如上图，在启动一个 API 命令后，API 命令会首先加载配置文件，根据配置做后面的处理工作。通常会将日志相关的配置记录在配置文件中，在解析完配置文件后，就可以加载日志包初始化函数，来初始化日志实例，供后面的程序调用。接下来会初始化数据库实例，建立数据库连接，供后面对数据库的 CRUD 操作使用。在建立完数据库连接后，需要设置 HTTP，通常包括 3 方面的设置：

1. 设置 Header
2. 注册路由
3. 注册中间件

之后会调用`net/http`包的`ListenAndServe()`方法启动 HTTP 服务器。

在启动 HTTP 端口之前，程序会 go 一个协程，来 ping HTTP 服务器的 `/sd/health` 接口，如果程序成功启动，ping 协程在 timeout 之前会成功返回，如果程序启动失败，则 ping 协程最终会 timeout，并终止整个程序。

>解析配置文件、初始化 Log 、初始化数据库的顺序根据自己的喜好和需求来排即可。

## HTTP 请求处理流程

[](./images/httphandle.png)

一次完整的 HTTP 请求处理流程如上图所示。

### 1. 建立连接

客户端发送 HTTP 请求后，服务器会根据域名进行域名解析，就是将网站名称转变成 IP 地址：localhost -> 127.0.0.1，Linux hosts文件、DNS 域名解析等可以实现这种功能。之后通过发起 TCP 的三次握手建立连接。TCP 三次连接请参考 TCP 三次握手详解及释放连接过程，建立连接之后就可以发送 HTTP 请求了。

### 2. 接收请求

HTTP 服务器软件进程，这里指的是 API 服务器，在接收到请求之后，首先根据 HTTP 请求行的信息来解析到 HTTP 方法和路径，在上图所示的报文中，方法是 GET，路径是 /index.html，之后根据 API 服务器注册的路由信息（大概可以理解为：HTTP 方法 + 路径和具体处理函数的映射）找到具体的处理函数。

### 3. 处理请求

在接收到请求之后，API 通常会解析 HTTP 请求报文获取请求头和消息体，然后根据这些信息进行相应的业务处理，HTTP 框架一般都有自带的解析函数，只需要输入 HTTP 请求报文，就可以解析到需要的请求头和消息体。通常情况下，业务逻辑处理可以分为两种：包含对数据库的操作和不包含对数据的操作。大型系统中通常两种都会有：

1. 包含对数据库的操作：需要访问数据库（增删改查），然后获取指定的数据，对数据处理后构建指定的响应结构体，返回响应包。数据库通常用的是 MySQL，因为免费，功能和性能也都能满足企业级应用的要求。
2. 不包含对数据库的操作：进行业务逻辑处理后，构建指定的响应结构体，返回响应包。

### 4. 记录事务处理过程

在业务逻辑处理过程中，需要记录一些关键信息，方便后期 Debug 用。在 Go 中有各种各样的日志包可以用来记录这些信息。

## HTTP 请求和响应格式介绍

一个 HTTP 请求报文由请求行（request line）、请求头部（header）、空行和请求数据四部分组成，下图是请求报文的一般格式。

[](./images/httpformat.png)

- 第一行必须是一个请求行（request line），用来说明请求类型、要访问的资源以及所使用的 HTTP 版本
- 紧接着是一个头部（header）小节，用来说明服务器要使用的附加信息
- 之后是一个空行
- 再后面可以添加任意的其他数据（称之为主体：body）

>HTTP 响应格式跟请求格式类似，也是由 4 个部分组成：状态行、消息报头、空行和响应数据。

## 目录结构

```
├── admin.sh                     # 进程的start|stop|status|restart控制文件
├── conf                         # 配置文件统一存放目录
│   ├── config.yaml              # 配置文件
│   ├── server.crt               # TLS配置文件
│   └── server.key
├── config                       # 专门用来处理配置和配置文件的Go package
│   └── config.go                 
├── db.sql                       # 在部署新环境时，可以登录MySQL客户端，执行source db.sql创建数据库和表
├── docs                         # swagger文档，执行 swag init 生成的
│   ├── docs.go
│   └── swagger
│       ├── swagger.json
│       └── swagger.yaml
├── handler                      # 类似MVC架构中的C，用来读取输入，并将处理流程转发给实际的处理函数，最后返回结果
│   ├── handler.go
│   ├── sd                       # 健康检查handler
│   │   └── check.go 
│   └── user                     # 核心：用户业务逻辑handler
│       ├── create.go            # 新增用户
│       ├── delete.go            # 删除用户
│       ├── get.go               # 获取指定的用户信息
│       ├── list.go              # 查询用户列表
│       ├── login.go             # 用户登录
│       ├── update.go            # 更新用户
│       └── user.go              # 存放用户handler公用的函数、结构体等
├── main.go                      # Go程序唯一入口
├── Makefile                     # Makefile文件，一般大型软件系统都是采用make来作为编译工具
├── model                        # 数据库相关的操作统一放在这里，包括数据库初始化和对表的增删改查
│   ├── init.go                  # 初始化和连接数据库
│   ├── model.go                 # 存放一些公用的go struct
│   └── user.go                  # 用户相关的数据库CURD操作
├── pkg                          # 引用的包
│   ├── auth                     # 认证包
│   │   └── auth.go
│   ├── constvar                 # 常量统一存放位置
│   │   └── constvar.go
│   ├── errno                    # 错误码存放位置
│   │   ├── code.go
│   │   └── errno.go
│   ├── token
│   │   └── token.go
│   └── version                  # 版本包
│       ├── base.go
│       ├── doc.go
│       └── version.go
├── README.md                    # API目录README
├── router                       # 路由相关处理
│   ├── middleware               # API服务器用的是Gin Web框架，Gin中间件存放位置
│   │   ├── auth.go
│   │   ├── header.go
│   │   ├── logging.go
│   │   └── requestid.go
│   └── router.go
├── service                      # 实际业务处理函数存放位
│   └── service.go
├── util                         # 工具类函数存放目录
│   ├── util.go
│   └── util_test.go
└── vendor                         # vendor目录用来管理依赖包
    ├── github.com
    ├── golang.org
    ├── gopkg.in
    └── vendor.json
```

Go API 项目中，一般都会包括这些功能项：Makefile 文件、配置文件目录、RESTful API 服务器的 handler 目录、model 目录、工具类目录、vendor 目录，以及实际处理业务逻辑函数所存放的 service 目录。这些都在上述的代码结构中有列出，新加功能时将代码放入对应功能的目录/文件中，可以使整个项目代码结构更加清晰，非常有利于后期的查找和维护。

# 4, 启动一个最简单的 RESTful API 服务器

## 本节核心内容

- 启动一个最简单的 RESTful API 服务器
- 设置 HTTP Header
- API 服务器健康检查和状态查询
- 编译并测试 API

## REST Web 框架选择

要编写一个 RESTful 风格的 API 服务器，首先需要一个 RESTful Web 框架，经过调研选择了 GitHub star 数最多的 Gin。采用轻量级的 Gin 框架，具有如下优点：高性能、扩展性强、稳定性强、相对而言比较简洁。

## 加载路由，并启动 HTTP 服务

main.go 中的 main() 函数是 Go 程序的入口函数，在 main() 函数中主要做一些配置文件解析、程序初始化和路由加载之类的事情，最终调用 http.ListenAndServe() 在指定端口启动一个 HTTP 服务器。本小节是一个简单的 HTTP 服务器，仅初始化一个 Gin 实例，加载路由并启动 HTTP 服务器。

### 编写入口函数

编写 `main()` 函数，main.go 代码：

```go
package main

import (
    "log"
    "net/http"

    "apiserver/router"

    "github.com/gin-gonic/gin"
)

func main() {
    // Create the Gin engine.
    g := gin.New()

    // gin middlewares
    middlewares := []gin.HandlerFunc{}

    // Routes.
    router.Load(
        // Cores.
        g,

        // Middlewares.
        middlewares...,
    )

    log.Printf("Start to listening the incoming requests on http address: %s", ":8080")
    log.Printf(http.ListenAndServe(":8080", g).Error())
}
```

### 加载路由

`main()` 函数通过调用 `router.Load` 函数来加载路由（函数路径为 router/router.go，具体函数实现参照 demo01/router/router.go）：

```go
"apiserver/handler/sd"

    ....

    // The health check handlers
    svcd := g.Group("/sd")
    {
        svcd.GET("/health", sd.HealthCheck)
        svcd.GET("/disk", sd.DiskCheck)
        svcd.GET("/cpu", sd.CPUCheck)
        svcd.GET("/ram", sd.RAMCheck)
    }
```

该代码块定义了一个叫 sd 的分组，在该分组下注册了 `/health`、`/disk`、`/cpu`、`/ram` HTTP 路径，分别路由到 `sd.HealthCheck`、`sd.DiskCheck`、`sd.CPUCheck`、`sd.RAMCheck` 函数。sd 分组主要用来检查 API Server 的状态：健康状况、服务器硬盘、CPU 和内存使用量。具体函数实现参照 demo01/handler/sd/check.go。

### 设置 HTTP Header

`router.Load` 函数通过 `g.Use()` 来为每一个请求设置 Header，在 router/router.go 文件中设置 Header：

```go
    g.Use(gin.Recovery())
    g.Use(middleware.NoCache)
    g.Use(middleware.Options)
    g.Use(middleware.Secure)
```

- `gin.Recovery()`：在处理某些请求时可能因为程序 bug 或者其他异常情况导致程序 panic，这时候为了不影响下一次请求的调用，需要通过 gin.Recovery()来恢复 API 服务器
- `middleware.NoCache`：强制浏览器不使用缓存
- `middleware.Options`：浏览器跨域 OPTIONS 请求设置
- `middleware.Secure`：一些安全设置

>`middleware`包的实现见 demo01/router/middleware。

## API 服务器健康状态自检

有时候 API 进程起来不代表 API 服务器正常，问题：API 进程存在，但是服务器却不能对外提供服务。因此在启动 API 服务器时，如果能够最后做一个自检会更好些。apiserver 中也添加了自检程序，在启动 HTTP 端口前 go 一个 `pingServer` 协程，启动 HTTP 端口后，该协程不断地 ping `/sd/health` 路径，如果失败次数超过一定次数，则终止 HTTP 服务器进程。通过自检可以最大程度地保证启动后的 API 服务器处于健康状态。自检部分代码位于 main.go 中：

```go
func main() {
    ....

    // Ping the server to make sure the router is working.
    go func() {
        if err := pingServer(); err != nil {
            log.Fatal("The router has no response, or it might took too long to start up.", err)
        }
        log.Print("The router has been deployed successfully.")
    }()
    ....
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
    for i := 0; i < 10; i++ {
        // Ping the server by sending a GET request to `/health`.
        resp, err := http.Get("http://127.0.0.1:8080" + "/sd/health")
        if err == nil && resp.StatusCode == 200 {
            return nil
        }

        // Sleep for a second to continue the next ping.
        log.Print("Waiting for the router, retry in 1 second.")
        time.Sleep(time.Second)
    }  
    return errors.New("Cannot connect to the router.")
}
```

在 `pingServer()` 函数中，`http.Get` 向 `http://127.0.0.1:8080/sd/health` 发送 HTTP GET 请求，如果函数正确执行并且返回的 HTTP StatusCode 为 200，则说明 API 服务器可用，`pingServer` 函数输出部署成功提示；如果超过指定次数，`pingServer` 直接终止 API Server 进程，如下图所示。

[](./images/terminal.png)

>`/sd/health` 路径会匹配到 `handler/sd/check.go` 中的 `HealthCheck` 函数，该函数只返回一个字符串：OK。

## 编译源码

1. 将`vendor`文件夹中的包拷贝到相应位置。
2. 做检查然后编译。

```sh
$ gofmt -w .
$ go tool vet .
$ go build -v .
```

>建议每次编译前对 Go 源码进行格式化和代码静态检查，以发现潜在的 Bug 或可疑的构造。

## cURL 工具测试 API

### cURL 工具简介

我们采用 cURL 工具来测试 RESTful API，标准的 Linux 发行版都安装了 cURL 工具。cURL 可以很方便地完成对 REST API 的调用场景，比如：设置 Header，指定 HTTP 请求方法，指定 HTTP 消息体，指定权限认证信息等。通过 -v 选项也能输出 REST 请求的所有返回信息。cURL 功能很强大，有很多参数，这里列出 REST 测试常用的参数：

```
-X/--request [GET|POST|PUT|DELETE|…]  指定请求的 HTTP 方法
-H/--header                           指定请求的 HTTP Header
-d/--data                             指定请求的 HTTP 消息体（Body）
-v/--verbose                          输出详细的返回信息
-u/--user                             指定账号、密码
-b/--cookie                           读取 cookie
```

典型的测试命令为：

```sh
$ curl -v -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/user -d'{"username":"admin","password":"admin1234"}'
```

### 启动 API Server

```sh
$ ./apiserver
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /sd/health                --> apiserver/handler/sd.HealthCheck (5 handlers)
[GIN-debug] GET    /sd/disk                  --> apiserver/handler/sd.DiskCheck (5 handlers)
[GIN-debug] GET    /sd/cpu                   --> apiserver/handler/sd.CPUCheck (5 handlers)
[GIN-debug] GET    /sd/ram                   --> apiserver/handler/sd.RAMCheck (5 handlers)
Start to listening the incoming requests on http address: :8080
The router has been deployed successfully.
```

### 发送 HTTP GET 请求

```sh
$ curl -XGET http://127.0.0.1:8080/sd/health
OK

$ curl -XGET http://127.0.0.1:8080/sd/disk
OK - Free space: 16321MB (15GB) / 51200MB (50GB) | Used: 31%

$ curl -XGET http://127.0.0.1:8080/sd/cpu
CRITICAL - Load average: 2.39, 2.13, 1.97 | Cores: 2

$ curl -XGET http://127.0.0.1:8080/sd/ram
OK - Free space: 455MB (0GB) / 8192MB (8GB) | Used: 5%
```

可以看到 HTTP 服务器均能正确响应请求。

# 5, 配置文件读取

## 本节核心内容

- 介绍 apiserver 所采用的配置解决方案
- 介绍如何配置 apiserver 并读取其配置，以及配置的高级用法

本小节使用`demo2`中的代码。

## Viper 简介

Viper 是国外大神 spf13 编写的开源配置解决方案，具有如下特性:

- 设置默认值
- 可以读取如下格式的配置文件：JSON、TOML、YAML、HCL
- 监控配置文件改动，并热加载配置文件
- 从环境变量读取配置
- 从远程配置中心读取配置（etcd/consul），并监控变动
- 从命令行 flag 读取配置
- 从缓存中读取配置
- 支持直接设置配置项的值

Viper 配置读取顺序：

- `viper.Set()` 所设置的值
- 命令行 flag
- 环境变量
- 配置文件
- 配置中心：etcd/consul
- 默认值

从上面这些特性来看，Viper 毫无疑问是非常强大的，而且 Viper 用起来也很方便，在初始化配置文件后，读取配置只需要调用 `viper.GetString()`、`viper.GetInt()` 和 `viper.GetBool()` 等函数即可。

Viper 也可以非常方便地读取多个层级的配置，比如这样一个 YAML 格式的配置：

```yaml
common:
  database:
    name: test
    host: 127.0.0.1
```

如果要读取 host 配置，执行 `viper.GetString("common.database.host")` 即可。

apiserver 采用 YAML 格式的配置文件，采用 YAML 格式，是因为 YAML 表达的格式更丰富，可读性更强。

## 初始化配置

### 主函数中增加配置初始化入口

```go
package main

import (
    "errors"
    "log"
    "net/http"
    "time"

    "apiserver/config"

     ...

    "github.com/spf13/pflag"
)

var (
    cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
    pflag.Parse()

    // init config
    if err := config.Init(*cfg); err != nil {
        panic(err)
    }

    // Create the Gin engine.
    g := gin.New()

    ...
}
```

在 `main` 函数中增加了 `config.Init(*cfg)` 调用，用来初始化配置，cfg 变量值从命令行 flag 传入，可以传值，比如 `./apiserver -c config.yaml`，也可以为空，如果为空会默认读取 `conf/config.yaml`。

### 解析配置

`main` 函数通过 `config.Init` 函数来解析并 `watch` 配置文件（函数路径：`config/config.go`），`config.go` 源码为：

```go
package config

import (
    "log"
    "strings"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

type Config struct {
    Name string
}

func Init(cfg string) error {
    c := Config {
        Name: cfg,
    }

    // 初始化配置文件
    if err := c.initConfig(); err != nil {
        return err
    }

    // 监控配置文件变化并热加载程序
    c.watchConfig()

    return nil
}

func (c *Config) initConfig() error {
    if c.Name != "" {
        viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
    } else {
        viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
        viper.SetConfigName("config")
    }
    viper.SetConfigType("yaml") // 设置配置文件格式为YAML
    viper.AutomaticEnv() // 读取匹配的环境变量
    viper.SetEnvPrefix("APISERVER") // 读取环境变量的前缀为APISERVER
    replacer := strings.NewReplacer(".", "_") 
    viper.SetEnvKeyReplacer(replacer)
    if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
        return err
    }

    return nil
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config file changed: %s", e.Name)
    })
}
```

`config.Init()` 通过 `initConfig()` 函数来解析配置文件，通过 `watchConfig()` 函数来 watch 配置文件，两个函数解析如下：

1. `func (c *Config) initConfig() error`

设置并解析配置文件。如果指定了配置文件 `*cfg` 不为空，则解析指定的配置文件，否则解析默认的配置文件 `conf/config.yaml`。通过指定配置文件可以很方便地连接不同的环境（开发环境、测试环境）并加载不同的配置，方便开发和测试。

通过如下设置:

```go
viper.AutomaticEnv()
viper.SetEnvPrefix("APISERVER")
replacer := strings.NewReplacer(".", "_")
```

可以使程序读取环境变量，具体效果稍后会演示。

`config.Init` 函数中的 `viper.ReadInConfig()` 函数最终会调用 Viper 解析配置文件。

2. `func (c *Config) watchConfig()`

通过该函数的 viper 设置，可以使 viper 监控配置文件变更，如有变更则热更新程序。所谓热更新是指：可以不重启 API 进程，使 API 加载最新配置项的值。

## 配置并读取配置

API 服务器端口号可能经常需要变更，API 服务器启动时间可能会变长，自检程序超时时间需要是可配的（通过设置次数），另外 API 需要根据不同的开发模式（开发、生产、测试）来匹配不同的行为。开发模式也需要是可配置的，这些都可以在配置文件中配置，新建配置文件 `conf/config.yaml`（默认配置文件名字固定为 `config.yaml`），`config.yaml` 的内容为：

```yaml
runmode: debug               # 开发模式, debug, release, test
addr: :6663                  # HTTP绑定端口
name: apiserver              # API Server的名字
url: http://127.0.0.1:6663   # pingServer函数请求的API服务器的ip:port
max_ping_count: 10           # pingServer函数尝试的次数
```

在 main 函数中将相应的配置改成从配置文件读取，需要替换的配置见下图中红框部分。

[](./images/tihuanqian.png)

替换后

[](./images/tihuanhou.png)

另外根据配置文件的 runmode 调用 gin.SetMode 来设置 gin 的运行模式：

```go
func main() { 
    pflag.Parse()

    // init config
    if err := config.Init(*cfg); err != nil {
        panic(err)
    }

    // Set gin mode.
    gin.SetMode(viper.GetString("runmode"))

    ....

}
```

gin 有 3 种运行模式：debug、release 和 test，其中 debug 模式会打印很多 debug 信息。

## 编译并运行

修改 `conf/config.yaml` 将端口修改为 `8888`，并启动 apiserver

修改后配置文件为：

```yaml
runmode: debug               # 开发模式, debug, release, test
addr: :8888                  # HTTP绑定端口
name: apiserver              # API Server的名字
url: http://127.0.0.1:8888   # pingServer函数请求的API服务器的ip:port
max_ping_count: 10           # pingServer函数try的次数
```

修改后启动 apiserver：

[](./images/qidong.png)

可以看到，启动 apiserver 后端口为配置文件中指定的端口。

## Viper 高级用法

### 从环境变量读取配置

在本节第一部分介绍过，Viper 可以从环境变量读取配置，这是个非常有用的功能。现在越来越多的程序是运行在 Kubernetes 容器集群中的，在 API 服务器迁移到容器集群时，可以直接通过 Kubernetes 来设置环境变量，然后程序读取设置的环境变量来配置 API 服务器。读者不需要了解如何通过 Kubernetes 设置环境变量，只需要知道 Viper 可以直接读取环境变量即可。

例如，通过环境变量来设置 API Server 端口：

```
$ export APISERVER_ADDR=:7777
$ export APISERVER_URL=http://127.0.0.1:7777
$ ./apiserver 
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /sd/health                --> apiserver/handler/sd.HealthCheck (5 handlers)
[GIN-debug] GET    /sd/disk                  --> apiserver/handler/sd.DiskCheck (5 handlers)
[GIN-debug] GET    /sd/cpu                   --> apiserver/handler/sd.CPUCheck (5 handlers)
[GIN-debug] GET    /sd/ram                   --> apiserver/handler/sd.RAMCheck (5 handlers)
Start to listening the incoming requests on http address: :7777
The router has been deployed successfully.
```

从输出可以看到，设置 `APISERVER_ADDR=:7777` 和 `APISERVER_URL=http://127.0.0.1:7777` 后，启动 apiserver，API 服务器的端口变为 `7777`。

环境变量名格式为 `config/config.go` 文件中 `viper.SetEnvPrefix("APISERVER")` 所设置的前缀和配置名称大写，二者用 `_` 连接，比如 `APISERVER_RUNMODE`。如果配置项是嵌套的，情况可类推，比如

```yaml
....
max_ping_count: 10           # pingServer函数try的次数
db:
  name: db_apiserver
```

对应的环境变量名为 `APISERVER_DB_NAME`。

### 热更新

在 `main` 函数中添加如下测试代码（`for {}` 部分，循环打印 `runmode` 的值）：

```go
import (
    "fmt"
    ....
)

var (
    cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
    pflag.Parse()

    // init config
    if err := config.Init(*cfg); err != nil {
        panic(err)
    }

    for {
        fmt.Println(viper.GetString("runmode"))
        time.Sleep(4*time.Second)
    }
    ....
}
```

编译并启动 `apiserver` 后，修改配置文件中 `runmode` 为 `test`，可以看到 `runmode` 的值从 `debug` 变为 `test`：

[](./images/peizhi.png)

# 6, 记录和管理 API 日志

## 本节核心内容

- Go 日志包数量众多，功能不同、性能不同，我们介绍一个比较好的日志库，并给出原因
- 介绍如何初始化日志包
- 介绍如何调用日志包
- 介绍如何转存（rotate）日志文件

>本节源码为 `demo3`

## 日志包介绍

apiserver 所采用的日志包 lexkong/log 是调研 GitHub 上的 开源log 包后封装的一个日志包。它参考华为 paas-lager，做了一些便捷性的改动，功能完全一样，只不过更为便捷。相较于 Go 的其他日志包，该日志包有如下特点：

- 支持日志输出流配置，可以输出到 stdout 或 file，也可以同时输出到 stdout 和 file
- 支持输出为 JSON 或 plaintext 格式
- 支持彩色输出
- 支持 log rotate 功能
- 高性能

## 初始化日志包

在 `conf/config.yaml` 中添加 log 配置

[](./images/初始化日志.png)

在 `config/config.go` 中添加日志初始化代码

```go
package config

import (
    ....
    "github.com/lexkong/log"
    ....
)
....
func Init(cfg string) error {
    ....
    // 初始化配置文件
    if err := c.initConfig(); err != nil {
        return err
    }

    // 初始化日志包
    c.initLog()
    ....
}

func (c *Config) initConfig() error {
    ....
}

func (c *Config) initLog() {
    passLagerCfg := log.PassLagerCfg {
        Writers:        viper.GetString("log.writers"),
        LoggerLevel:    viper.GetString("log.logger_level"),
        LoggerFile:     viper.GetString("log.logger_file"),
        LogFormatText:  viper.GetBool("log.log_format_text"),
        RollingPolicy:  viper.GetString("log.rollingPolicy"),
        LogRotateDate:  viper.GetInt("log.log_rotate_date"),
        LogRotateSize:  viper.GetInt("log.log_rotate_size"),
        LogBackupCount: viper.GetInt("log.log_backup_count"),
    }

    log.InitWithConfig(&passLagerCfg)
}  

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
    ....
}
```

这里要注意，日志初始化函数 `c.initLog()` 要放在配置初始化函数 `c.initConfig()` 之后，因为日志初始化函数要读取日志相关的配置。`func (c *Config) initLog()` 是日志初始化函数，会设置日志包的各项参数，参数为：

- `writers`：输出位置，有两个可选项 —— file 和 stdout。选择 file 会将日志记录到 logger_file 指定的日志文件中，选择 stdout 会将日志输出到标准输出，当然也可以两者同时选择
- `logger_level`：日志级别，DEBUG、INFO、WARN、ERROR、FATAL
- `logger_file`：日志文件
- `log_format_text`：日志的输出格式，JSON 或者 plaintext，true 会输出成非 JSON 格式，false 会输出成 JSON 格式
- `rollingPolicy`：rotate 依据，可选的有 daily 和 size。如果选 daily 则根据天进行转存，如果是 size 则根据大小进行转存
- `log_rotate_date`：rotate 转存时间，配 合rollingPolicy: daily 使用
- `log_rotate_size`：rotate 转存大小，配合 rollingPolicy: size 使用
- `log_backup_count`：当日志文件达到转存标准时，log 系统会将该日志文件进行压缩备份，这里指定了备份文件的最大个数

## 调用日志包

日志初始化好了，将 demo02 中的 log 用 lexkong/log 包来替换。替换前（这里 grep 出了需要替换的行，可自行确认替换后的效果）：

```go
$ grep log * -R
config/config.go:	"log"
config/config.go:		log.Printf("Config file changed: %s", e.Name)
main.go:	"log"
main.go:			log.Fatal("The router has no response, or it might took too long to start up.", err)
main.go:		log.Print("The router has been deployed successfully.")
main.go:	log.Printf("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
main.go:	log.Printf(http.ListenAndServe(viper.GetString("addr"), g).Error())
main.go:		log.Print("Waiting for the router, retry in 1 second.")
```

替换后的源码文件见 demo03。

## 编译并运行

启动后，可以看到 apiserver 有 JSON 格式的日志输出：

[](./images/json格式日志.png)

## 管理日志文件

这里将日志转存策略设置为 size，转存大小设置为 1 MB

```yaml
  rollingPolicy: size
  log_rotate_size: 1
```

并在 main 函数中加入测试代码：

[](./images/日志测试代码.png)

启动 apiserver 后发现，在当前目录下创建了 log/apiserver.log 日志文件：

```
$ ls log/
apiserver.log
```

程序运行一段时间后，发现又创建了 zip 文件：

```
$ ls log/
apiserver.log  apiserver.log.20180531134509631.zip
```

该 zip 文件就是当 apiserver.log 大小超过 1MB 后，日志系统将之前的日志压缩成 zip 文件后的文件。

# 7, 初始化表

## 创建示例需要的数据库和表

1. 创建 db.sql，内容为：

```sql
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `db_apiserver` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `db_apiserver`;

--
-- Table structure for table `tb_users`
--

DROP TABLE IF EXISTS `tb_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tb_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `createdAt` timestamp NULL DEFAULT NULL,
  `updatedAt` timestamp NULL DEFAULT NULL,
  `deletedAt` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `idx_tb_users_deletedAt` (`deletedAt`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tb_users`
--

LOCK TABLES `tb_users` WRITE;
/*!40000 ALTER TABLE `tb_users` DISABLE KEYS */;
INSERT INTO `tb_users` VALUES (0,'admin','$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG','2018-05-27 16:25:33','2018-05-27 16:25:33',NULL);
/*!40000 ALTER TABLE `tb_users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-05-28  0:25:41
```

2. 登录 MySQL:

```
$ mysql -uroot -p 
```

3. source db.sql

```
mysql> source db.sql
```

可以看到，db.sql 创建了 db_apiserver 数据库和 tb_users 表，并默认添加了一条记录（用户名：admin，密码：admin）：

```sql
mysql> select * from tb_users \G;
*************************** 1. row ***************************
       id: 0
 username: admin
 password: $2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG
createdAt: 2018-05-28 00:25:33
updatedAt: 2018-05-28 00:25:33
deletedAt: NULL
1 row in set (0.00 sec)
```

## 在配置文件中添加数据库配置

API 启动需要连接数据库，所以需要在配置文件 conf/config.yaml 中配置数据库的 IP、端口、用户名、密码和数据库名信息。

[](./images/数据库配置.png)

# 8, 初始化 MySQL 数据库并建立连接

## 本节核心内容

- Go ORM 数量众多，我们介绍一个比较好的 ORM 包，并给出原因
- 介绍如何初始化数据库
- 介绍如何连接数据库

>源码在 demo04

apiserver 用的 ORM 是 GitHub 上 star 数最多的 gorm，相较于其他 ORM，它用起来更方便，更稳定，社区也更活跃。

## 初始化数据库

在 model/init.go 中添加数据初始化代码

因为一个 API 服务器可能需要同时访问多个数据库，为了对多个数据库进行初始化和连接管理，这里定义了一个叫 Database 的 struct：

```go
type Database struct {
    Self   *gorm.DB
    Docker *gorm.DB
}
```

Database 结构体有个 Init() 方法用来初始化连接：

```go
func (db *Database) Init() {
    DB = &Database {
        Self:   GetSelfDB(),
        Docker: GetDockerDB(),
    }
}
```

Init() 函数会调用 GetSelfDB() 和 GetDockerDB() 方法来同时创建两个 Database 的数据库对象。这两个 Get 方法最终都会调用 func openDB(username, password, addr, name string) *gorm.DB 方法来建立数据库连接，不同数据库实例传入不同的 username、password、addr 和名字信息，从而建立不同的数据库连接。openDB 函数为：

```go
func openDB(username, password, addr, name string) *gorm.DB {
    config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
        username,
        password,
        addr,
        name,
        true,
        //"Asia/Shanghai"),
        "Local")

    db, err := gorm.Open("mysql", config)
    if err != nil {
        log.Errorf(err, "Database connection failed. Database name: %s", name)
    }  

    // set for db connection
    setupDB(db)

    return db
}
```

可以看到，openDB() 最终调用 gorm.Open() 来建立一个数据库连接。

完整的 model/init.go 源码文件请参考 demo04/model/init.go。

### 主函数中增加数据库初始化入口

```go
package main

import (
    ...
    "apiserver/model"

    ...
)

...

func main() {
    ...

    // init db
    model.DB.Init()
    defer model.DB.Close()

    ...
}
```

通过 model.DB.Init() 来建立数据库连接，通过 defer model.DB.Close() 来关闭数据库连接。

# 9, 自定义业务错误信息

## 本节核心内容

- 如何自定义业务自己的错误信息
- 实际开发中是如何处理错误的
- 实际开发中常见的错误类型
- 通过引入新包 errno 来实现此功能，会展示该包的如下用法：
    - 如何新建 Err 类型的错误
    - 如何从 Err 类型的错误中获取 code 和 message

>源码路径：demo05

## 为什么要定制业务自己的错误码

在实际开发中引入错误码有如下好处：

- 可以非常方便地定位问题和定位代码行（看到错误码知道什么意思，grep 错误码可以定位到错误码所在行）
- 如果 API 对外开放，有个错误码会更专业些
- 错误码包含一定的信息，通过错误码可以判断出错误级别、错误模块和具体错误信息
- 在实际业务开发中，一个条错误信息需要包含两部分内容：直接展示给用户的 message 和用于开发人员 debug 的 error 。message 可能会直接展示给用户，error 是用于 debug 的错误信息，可能包含敏感/内部信息，不宜对外展示
- 业务开发过程中，可能需要判断错误是哪种类型以便做相应的逻辑处理，通过定制的错误码很容易做到这点，例如：

```go
    if err == errno.ErrBind {
        ...
    }
```

- Go 中的 HTTP 服务器开发都是引用 net/http 包，该包中只有 60 个错误码，基本都是跟 HTTP 请求相关的。在大型系统中，这些错误码完全不够用，而且跟业务没有任何关联，满足不了业务需求。

## 在 apiserver 中引入错误码

我们通过一个新包 errno 来做错误码的定制，详见 demo05/pkg/errno。

```
$ ls pkg/errno/
code.go  errno.go
```

errno 包由两个 Go 文件组成：code.go 和 errno.go。code.go 用来统一存自定义的错误码，code.go 的代码为：

```go
package errno

var (
    // Common errors
    OK                  = &Errno{Code: 0, Message: "OK"}
    InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
    ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

    // user errors
    ErrUserNotFound      = &Errno{Code: 20102, Message: "The user was not found."}
)
```

### 代码解析

在实际开发中，一个错误类型通常包含两部分：Code 部分，用来唯一标识一个错误；Message 部分，用来展示错误信息，这部分错误信息通常供前端直接展示。这两部分映射在 errno 包中即为 &Errno{Code: 0, Message: "OK"}。

### 错误码设计

目前错误码没有一个统一的设计标准，笔者研究了 BAT 和新浪开放平台对外公布的错误码设计，参考新浪开放平台 Error code 的设计，如下是设计说明：

错误返回值格式：

```json
{
  "code": 10002,
  "message": "Error occurred while binding the request body to the struct."
}
```

错误代码说明：

| 1 | 00 | 02 |
|---|----|----|
| 服务级错误（1 为系统级错误）| 服务模块代码 | 具体错误代码 |

- 服务级别错误：1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的
- 服务模块为两位数：一个大型系统的服务模块通常不超过两位数，如果超过，说明这个系统该拆分了
- 错误码为两位数：防止一个模块定制过多的错误码，后期不好维护
- code = 0 说明是正确返回，code > 0 说明是错误返回
- 错误通常包括系统级错误码和服务级错误码
- 建议代码中按服务模块将错误分类
- 错误码均为 >= 0 的数
- 在 apiserver 中 HTTP Code 固定为 http.StatusOK，错误码通过 code 来表示。

## 错误信息处理

通过 errno.go 来对自定义的错误进行处理，errno.go 的代码为：

```go
package errno

import "fmt"

type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

func (err *Err) Add(message string) error {
	err.Message += " " + message
	return err
}

func (err *Err) Addf(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)
	return err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

func IsErrUserNotFound(err error) bool {
	code, _ := DecodeErr(err)
	return code == ErrUserNotFound.Code
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalServerError.Code, err.Error()
}
```

### 代码解析

errno.go 源码文件中有两个核心函数 New() 和 DecodeErr()，一个用来新建定制的错误，一个用来解析定制的错误，稍后会介绍如何使用。

errno.go 同时也提供了 Add() 和 Addf() 函数，如果想对外展示更多的信息可以调用此函数，使用方法下面有介绍。

## 错误码实战

上面介绍了错误码的一些知识，这一部分讲开发中是如何使用 errno 包来处理错误信息的。为了演示，我们新增一个创建用户的 API：

1. router/router.go 中添加路由，详见 demo05/router/router.go：

[](./images/添加路由.png)

2. handler 目录下增加业务处理函数 handler/user/create.go，详见 demo05/handler/user/create.go。

## 编译并运行

## 测试验证

启动 apiserver：./apiserver

```
$ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user

{
  "code": 10002,
  "message": "Error occurred while binding the request body to the struct."
}
```

因为没有传入任何参数，所以返回 errno.ErrBind 错误。

```
$ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user -d'{"username":"admin"}'

{
  "code": 10001,
  "message": "password is empty"
}
```

因为没有传入 password，所以返回 fmt.Errorf("password is empty") 错误，该错误信息不是定制的错误类型，errno.DecodeErr(err) 解析时会解析为默认的 errno.InternalServerError 错误，所以返回结果中 code 为 10001，message 为 err.Error()。

```
$ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user -d'{"password":"admin"}'

{
  "code": 20102,
  "message": "The user was not found. This is add message."
}
```

因为没有传入 username，所以返回 errno.ErrUserNotFound 错误信息，并通过 Add() 函数在 message 信息后追加了 This is add message. 信息。

通过

```go
   if errno.IsErrUserNotFound(err) {
        log.Debug("err type is ErrUserNotFound")
    }
```

演示了如何通过定制错误方便地对比是不是某个错误，在该请求中，apiserver 会输出如下错误：

[](./images/输出错误.png)

可以看到在后台日志中会输出敏感信息 username can not found in db: xx.xx.xx.xx，但是返回给用户的 message （{"code":20102,"message":"The user was not found. This is add message."}）不包含这些敏感信息，可以供前端直接对外展示。

```
$ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user -d'{"username":"admin","password":"admin"}'

{
  "code": 0,
  "message": "OK"
}
```

如果 err = nil，则 errno.DecodeErr(err) 会返回成功的 code: 0 和 message: OK。

>如果 API 是对外的，错误信息数量有限，则制定错误码非常容易，强烈建议使用错误码。如果是内部系统，特别是庞大的系统，内部错误会非常多，这时候没必要为每一个错误制定错误码，而只需为常见的错误制定错误码，对于普通的错误，系统在处理时会统一作为 InternalServerError 处理。

# 10, 读取和返回 HTTP 请求

## 本节核心内容

- 如何读取 HTTP 请求数据
- 如何返回数据
- 如何定制业务的返回格式

本小节源码下载路径：demo06

## 读取和返回参数

在业务开发过程中，需要读取请求参数（消息体和 HTTP Header），经过业务处理后返回指定格式的消息。apiserver 也展示了如何进行参数的读取和返回，下面展示了如何读取和返回参数：

读取 HTTP 信息： 在 API 开发中需要读取的参数通常为：HTTP Header、路径参数、URL参数、消息体，读取这些参数可以直接使用 gin 框架自带的函数：

- Param()：返回 URL 的参数值，例如

```go
router.GET("/user/:id", func(c *gin.Context) {
    // a GET request to /user/john
    id := c.Param("id") // id == "john"
})
```

- Query()：读取 URL 中的地址参数，例如

```go
// GET /path?id=1234&name=Manu&value=
   c.Query("id") == "1234"
   c.Query("name") == "Manu"
   c.Query("value") == ""
   c.Query("wtf") == ""
```

- DefaultQuery()：类似 Query()，但是如果 key 不存在，会返回默认值，例如

```go
//GET /?name=Manu&lastname=
c.DefaultQuery("name", "unknown") == "Manu"
c.DefaultQuery("id", "none") == "none"
c.DefaultQuery("lastname", "none") == ""
```

- Bind()：检查 Content-Type 类型，将消息体作为指定的格式解析到 Go struct 变量中。apiserver 采用的媒体类型是 JSON，所以 Bind() 是按 JSON 格式解析的。

- GetHeader()：获取 HTTP 头。

返回HTTP消息： 因为要返回指定的格式，apiserver 封装了自己的返回函数，通过统一的返回函数 SendResponse 来格式化返回，小节后续部分有详细介绍。

## 增加返回函数

API 返回入口函数，供所有的服务模块返回时调用，所以这里将入口函数添加在 handler 目录下，handler/handler.go 的源码为：

```go
package handler

import (
	"net/http"

	"apiserver/pkg/errno"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
```

可以看到返回格式固定为：

```go
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
```

在返回结构体中，固定有 Code 和 Message 参数，这两个参数通过函数 DecodeErr() 解析 error 类型的变量而来（DecodeErr() 在上一节介绍过）。Data 域为 interface{} 类型，可以根据业务自己的需求来返回，可以是 map、int、string、struct、array 等 Go 语言变量类型。SendResponse() 函数通过 errno.DecodeErr(err) 来解析出 code 和 message，并填充在 Response 结构体中。

## 在业务处理函数中读取和返回数据

通过改写上一节 handler/user/create.go 源文件中的 Create() 函数，来演示如何读取和返回数据，改写后的源码为：

```go
package user

import (
	"fmt"

	. "apiserver/handler"
	"apiserver/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// Create creates a new user account.
func Create(c *gin.Context) {
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	admin2 := c.Param("username")
	log.Infof("URL username: %s", admin2)

	desc := c.Query("desc")
	log.Infof("URL key param desc: %s", desc)

	contentType := c.GetHeader("Content-Type")
	log.Infof("Header Content-Type: %s", contentType)

	log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
	if r.Username == "" {
		SendResponse(c, errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not found in db: xx.xx.xx.xx")), nil)
		return
	}

	if r.Password == "" {
		SendResponse(c, fmt.Errorf("password is empty"), nil)
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	SendResponse(c, nil, rsp)
}
```

这里也需要更新下路由，router/router.go（详见 demo06/router/router.go）：

[](./images/更新路由.png)

上例展示了如何通过 Bind()、Param()、Query() 和 GetHeader() 来获取相应的参数。

根据笔者的研发经验，建议：如果消息体有 JSON 参数需要传递，针对每一个 API 接口定义独立的 go struct 来接收，比如 CreateRequest 和 CreateResponse，并将这些结构体统一放在一个 Go 文件中，以方便后期维护和修改。这样做可以使代码结构更加规整和清晰，本例统一放在 handler/user/user.go 中，源码为：

```go
package user

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
}
```

## 编译并运行

## 测试

启动apiserver：./apiserver，发送 HTTP 请求：

```sh
$ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user/admin2?desc=test -d'{"username":"admin","password":"admin"}'

{
  "code": 0,
  "message": "OK",
  "data": {
    "username": "admin"
  }
}
```

查看 apiserver 日志：

[](./images/api日志.png)

可以看到成功读取了请求中的各类参数。并且 curl 命令返回的结果格式为指定的格式：

```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "username": "admin"
  }
}
```

# 11, 用户业务逻辑处理

## 本节核心内容

这一节是核心小节，讲解如何处理用户业务，这也是 API 的核心功能。本小节会讲解实际开发中需要的一些重要功能点。功能点包括：

- 各种场景的业务逻辑处理
    - 创建用户
    - 删除用户
    - 更新用户
    - 查询用户列表
    - 查询指定用户的信息
- 数据库的 CURD 操作

>本小节源码下载路径：demo07

## 配置路由信息

需要先在 router/router.go 文件中，配置路由信息：

```go
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
    ...
	// 用户路由设置
	u := g.Group("/v1/user")
	{
		u.POST("", user.Create)         // 创建用户
		u.DELETE("/:id", user.Delete)   // 删除用户 
		u.PUT("/:id", user.Update)      // 更新用户
		u.GET("", user.List)            // 用户列表
		u.GET("/:username", user.Get)   // 获取指定用户的详细信息
	}
    ...
	return g
}
```

>在 RESTful API 开发中，API 经常会变动，为了兼容老的 API，引入了版本的概念，比如上例中的 /v1/user，说明该 API 版本是 v1。
>
>很多 RESTful API 最佳实践文章中均建议使用版本控制，笔者这里也建议对 API 使用版本控制。

## 注册新的错误码

在 pkg/errno/code.go 文件中（详见 demo07/pkg/errno/code.go），新增如下错误码：

```go
var (
	// Common errors
        ...

	ErrValidation       = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase         = &Errno{Code: 20002, Message: "Database error."}
	ErrToken            = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}

	// user errors
	ErrEncrypt           = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrTokenInvalid      = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect = &Errno{Code: 20104, Message: "The password was incorrect."}
)
```

## 新增用户

更新 handler/user/create.go 中 Create() 的逻辑，更新后的内容见 demo07/handler/user/create.go。

创建用户逻辑：

- 从 HTTP 消息体获取参数（用户名和密码）
- 参数校验
- 加密密码
- 在数据库中添加数据记录
- 返回结果（这里是用户名）

从 HTTP 消息体解析参数，前面小节已经介绍了。

参数校验这里用的是 gopkg.in/go-playground/validator.v9 包（详见 go-playground/validator），实际开发过程中，该包可能不能满足校验需求，这时候可在程序中加入自己的校验逻辑，比如在 handler/user/creater.go 中添加校验函数 checkParam：

```go
package user

import (
    ...
)

// Create creates a new user account.
func Create(c *gin.Context) {
	log.Info("User Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	if err := r.checkParam(); err != nil {
		SendResponse(c, err, nil)
		return
	}
        ...
}

func (r *CreateRequest) checkParam() error {
	if r.Username == "" {
		return errno.New(errno.ErrValidation, nil).Add("username is empty.")
	}

	if r.Password == "" {
		return errno.New(errno.ErrValidation, nil).Add("password is empty.")
	}

	return nil
}
```

例子通过 Encrypt() 对密码进行加密：

```go
// Encrypt the user password.
func (u *UserModel) Encrypt() (err error) {
    u.Password, err = auth.Encrypt(u.Password)
    return
}
```

Encrypt() 函数引用 auth.Encrypt() 来进行密码加密，具体实现见 demo07/pkg/auth/auth.go。

最后例子通过 u.Create() 函数来向数据库中添加记录，ORM 用的是 gorm，gorm 详细用法请参考 GORM 指南。在 Create() 函数中引用的数据库实例是 DB.Self，该实例在 API 启动之前已经完成初始化。DB 是个全局变量，可以直接引用。

>在实际开发中，为了安全，数据库中是禁止保存密码的明文信息的，密码需要加密保存。
>
>我们将接收和处理相关的 Go 结构体统一放在 handler/user/user.go 文件中，这样可以使程序结构更清晰，功能更聚焦。当然每个人习惯不一样，大家根据自己的习惯放置即可。handler/user/user.go 对 UserInfo 结构体的处理，也出于相同的目的。

## 删除用户

删除用户代码详见 demo07/handler/user/delete.go。

删除时，首先根据 URL 路径 DELETE http://127.0.0.1/v1/user/1 解析出 id 的值 1，该 id 实际上就是数据库中的 id 索引，调用 model.DeleteUser() 函数删除，函数详见 demo07/model/user.go。

## 更新用户

更新用户代码详见 demo07/handler/user/update.go。

更新用户逻辑跟创建用户差不多，在更新完数据库字段后，需要指定 gorm model 中的 id 字段的值，因为 gorm 在更新时默认是按照 id 来匹配记录的。通过解析 PUT http://127.0.0.1/v1/user/1 来获取 id。

## 查询用户列表

查询用户列表代码详见 demo07/handler/user/list.go。

一般在 handler 中主要做解析参数、返回数据操作，简单的逻辑也可以在 handler 中做，像新增用户、删除用户、更新用户，代码量不大，所以也可以放在 handler 中。有些代码量很大的逻辑就不适合放在 handler 中，因为这样会导致 handler 逻辑不是很清晰，这时候实际处理的部分通常放在 service 包中。比如本例的 LisUser() 函数：

```go
package user
   
import (
    "apiserver/service"
    ...
)  
   
// List list the users in the database.
func List(c *gin.Context) {
    ...
    infos, count, err := service.ListUser(r.Username, r.Offset, r.Limit)
    if err != nil {
        SendResponse(c, err, nil)
        return
    }
    ...
}
```

查询一个 REST 资源列表，通常需要做分页，如果不做分页返回的列表过多，会导致 API 响应很慢，前端体验也不好。本例中的查询函数做了分页，收到的请求中传入的 offset 和 limit 参数，分别对应于 MySQL 的 offset 和 limit。

service.ListUser() 函数用来做具体的查询处理，代码详见 demo07/service/service.go。

在 ListUser() 函数中用了 sync 包来做并行查询，以使响应延时更小。在实际开发中，查询数据后，通常需要对数据做一些处理，比如 ListUser() 函数中会对每个用户记录返回一个 sayHello 字段。sayHello 只是简单输出了一个 Hello shortId 字符串，其中 shortId 是通过 util.GenShortId() 来生成的（GenShortId 实现详见 demo07/util/util.go）。像这类操作通常会增加 API 的响应延时，如果列表条目过多，列表中的每个记录都要做一些类似的逻辑处理，这会使得整个 API 延时很高，所以笔者在实际开发中通常会做并行处理。根据笔者经验，效果提升十分明显。

读者应该已经注意到了，在 ListUser() 实现中，有 sync.Mutex 和 IdMap 等部分代码，使用 sync.Mutex 是因为在并发处理中，更新同一个变量为了保证数据一致性，通常需要做锁处理。

使用 IdMap 是因为查询的列表通常需要按时间顺序进行排序，一般数据库查询后的列表已经排过序了，但是为了减少延时，程序中用了并发，这时候会打乱排序，所以通过 IdMap 来记录并发处理前的顺序，处理后再重新复位。

## 获取指定用户的详细信息

代码详见 demo07/handler/user/get.go。

获取指定用户信息时，首先根据 URL 路径 GET http://127.0.0.1/v1/user/admin 解析出 username 的值 admin，然后调用 model.GetUser() 函数查询该用户的数据库记录并返回，函数详见 demo07/model/user.go。

## 编译并运行

### 创建用户

```sh
$ curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user -d'{"username":"kong","password":"kong123"}'

{
  "code": 0,
  "message": "OK",
  "data": {
    "username": "kong"
  }
}
```

### 查询用户列表

```sh
$ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user -d'{"offset": 0, "limit": 20}'

{
  "code": 0,
  "message": "OK",
  "data": {
    "totalCount": 2,
    "userList": [
      {
        "id": 2,
        "username": "kong",
        "sayHello": "Hello qhXO5iIig",
        "password": "$2a$10$vE9jG71oyzstWVwB/QfU3u00Pxb.ye8hFIDvnyw60nHBv/xsJZoUO",
        "createdAt": "2018-06-02 14:47:54",
        "updatedAt": "2018-06-02 14:47:54"
      },
      {
        "id": 0,
        "username": "admin",
        "sayHello": "Hello qhXO5iSmgz",
        "password": "$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG",
        "createdAt": "2018-05-28 00:25:33",
        "updatedAt": "2018-05-28 00:25:33"
      }
    ]
  }
}
```

可以看到，新增了一个用户 kong，数据库 id 索引为 2。admin 用户是上一节中初始化数据库时初始化的。

>建议在 API 设计时，对资源列表进行分页。

### 获取用户详细信息

```sh
$ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user/kong

{
  "code": 0,
  "message": "OK",
  "data": {
    "username": "kong",
    "password": "$2a$10$vE9jG71oyzstWVwB/QfU3u00Pxb.ye8hFIDvnyw60nHBv/xsJZoUO"
  }
}
```

### 更新用户

在 查询用户列表 部分，会返回用户的数据库索引。例如，用户 kong 的数据库 id 索引是 2，所以这里调用如下 URL 更新 kong 用户：

```sh
$ curl -XPUT -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user/2 -d'{"username":"kong","password":"kongmodify"}'

{
  "code": 0,
  "message": "OK",
  "data": null
}
```

获取 kong 用户信息：

```sh
$ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user/kong

{
  "code": 0,
  "message": "OK",
  "data": {
    "username": "kong",
    "password": "$2a$10$E0kwtmtLZbwW/bDQ8qI8e.eHPqhQOW9tvjwpyo/p05f/f4Qvr3OmS"
  }
}
```

可以看到密码已经改变（旧密码为 $2a$10$vE9jG71oyzstWVwB/QfU3u00Pxb.ye8hFIDvnyw60nHBv/xsJZoUO）。

### 删除用户

在 查询用户列表 部分，会返回用户的数据库索引。例如，用户 kong 的数据库 id 索引是 2，所以这里调用如下 URL 删除 kong 用户：

```sh
$ curl -XDELETE -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user/2

{
  "code": 0,
  "message": "OK",
  "data": null
}
```

获取用户列表：

```sh
$ curl -XGET -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user -d'{"offset": 0, "limit": 20}'

{
  "code": 0,
  "message": "OK",
  "data": {
    "totalCount": 1,
    "userList": [
      {
        "id": 0,
        "username": "admin",
        "sayHello": "Hello EnqntiSig",
        "password": "$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG",
        "createdAt": "2018-05-28 00:25:33",
        "updatedAt": "2018-05-28 00:25:33"
      }
    ]
  }
}
```

可以看到用户 kong 未出现在用户列表中，说明他已被成功删除。