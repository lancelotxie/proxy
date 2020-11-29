# Proxy
Proxy是一款干净的正向代理工具,支持HTTP/HTTPS协议,内存占用低,CPU占用低,智能分流,内网直连

## 安装
- Proxy使用Golang开发,安装前请自行安装Golang开发环境
- Proxy运行平台:Linux
- 编译:
  ```sh
  git clone https://github.com/lancelotXie/proxy.git
  ```
  1. 编译服务端:
        ```sh
        cd ./proxy/proxy.server
        ```
        ```sh
        go build -o server.o ./main.go
        ```
    2. 编译客户端:
        ```sh
        cd ./proxy/proxy.client
        ```
        ```sh
        go build -o client.o ./main.go
        ```
    3. 编译配置端:
        ```sh
        cd ./proxy/proxy.controller
        ```
        ```sh
        go build -o controller.o ./main.go
        ```
- 配置监听地址:
    1. 配置服务端:
        ```sh
        cd ./proxy/proxy.server
        ```
        ```sh
        ./server.o -ctrl-port 8086
        ```
        ```sh
        # 另外开一个终端
        cd ./proxy/proxy.controller &&
        ./controller.o -ctrl-port 8086 set listen.addr 0.0.0.0:17600 &&
        ./controller.o -ctrl-port 8086 save
        ```
    2. 配置客户端,假设服务端IP为123.123.123.123:
        ```sh
        cd ./proxy/proxy.client
        ```
        ```sh
        ./client.o
        ```
        ```sh
        # 另外开一个终端
        cd ./proxy/proxy.controller &&
        ./controller.o set remote.addr 123.123.123.123:17600 &&
        ./controller.o save
        ```
- 运行:
    1. 运行服务端:
        ```sh
        cd ./proxy/proxy.server && ./server.o
        ```
    2. 运行客户端:
        ```sh
        cd ./proxy/proxy.client && ./client.o
        ```
- 设置系统代理为127.0.0.1:9527
- 安装完毕