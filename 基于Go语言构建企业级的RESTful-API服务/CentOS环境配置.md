1. 下载安装包

```sh
$ wget https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz
```

2. 设置安装目录

```sh
$ export GO_INSTALL_DIR=$HOME
```

3. 解压 Go 安装包

```sh
$ tar -xvzf go1.10.2.linux-amd64.tar.gz -C $GO_INSTALL_DIR
```

4. 设置环境变量

```sh
$ export GO_INSTALL_DIR=$HOME
$ export GOROOT=$GO_INSTALL_DIR/go
$ export GOPATH=$HOME/mygo
$ export PATH=$GOPATH/bin:$PATH:$GO_INSTALL_DIR/go/bin
```

如果不想每次登录系统都设置一次环境变量，可以将上面 4 行追加到 $HOME/.bashrc 文件中。

5. 执行 go version 检查 Go 是否成功安装

```sh
$ go version
go version go1.10.2 linux/amd64
```

看到 go version 命令输出 go 版本号 go1.10.2 linux/amd64，说明 go 命令安装成功。

6. 创建 $GOPATH/src 目录

$GOPATH/src 是 Go 源码存放的目录，所以在正式开始编码前要先确保 $GOPATH/src 目录存在，执行命令：

```sh
$ mkdir -p $GOPATH/src
```

7. 安装 git

```sh
$ yum install git
```

8. clone 项目

将vendor中的文件夹拷贝到 src 目录下面。

9. 安装 MySQL

```sh
$ rpm -q mariadb-server
```

检查是否安装了 MySQL 。

如果提示 package mariadb-server is not installed 则说明没有安装 MySQL，需要手动安装。如果出现 mariadb-server-xxx.xxx.xx.el7.x86_64 则说明已经安装。

安装 MySQL 的步骤为：

1. 安装 MySQL 和 MySQL 客户端

```sh
$ sudo yum -y install mariadb  mariadb-server
```

2. 启动 MySQL

```sh
$ sudo systemctl start mariadb
```

3. 设置开机启动

```sh
$ sudo systemctl enable mariadb
```

4. 设置初始密码

```sh
$ sudo mysqladmin -u root password root
```