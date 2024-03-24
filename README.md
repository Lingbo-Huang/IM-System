
# IM-System

是一个即时通信系统，用于将golang学习的结果和基础的net包里conn的使用都串起来，实现的功能：服务器响应，客户端集成，用户上线，广播，私聊，在线用户查询，修改用户名，超时强制踢出功能。


## 作者

- [@Lingbo-Huang](https://www.github.com/octokatherine)


## 截图

![上线](https://github.com/Lingbo-Huang/Images/blob/main/img/00db6a3c112e083ae37a3e792aa9cba.png)

![私聊](https://github.com/Lingbo-Huang/Images/blob/main/img/51d1ce29e0396ae79b206cc09aa9948.png)

![聊天](https://github.com/Lingbo-Huang/Images/blob/main/img/c4895783d2b42c3ecfeb2a87d0c9e46.png)

![退出](https://github.com/Lingbo-Huang/Images/blob/main/img/6971b6f585cdc66bd77e31d39f80a9e.png)

![超时强制踢出](https://github.com/Lingbo-Huang/Images/blob/main/img/6f22b1c8d41516d942f158d89b63534.png)



## 运行

要部署这个项目，请运行

```bash
  build -o server server.go main.go user.go
  build -o client client.go
```
启动server
```
./server
```
启动客户端
```
./client
或./client -ip 127.0.0.1 -port 8888
或监听 nc 127.0.0.1 8888
```

