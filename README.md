# gmessage
```
1.客户端client TCP三次握手连接服务端server.
连接成功后,客户端首先需要发送funcid=1的数据包,绑定当前的socket和userid,使当前的连接有状态(一对一的关系)..这样服务端也可以通过userid找到对应的客户端socket,然后给客户端推送消息.

2.客户端正常的发送数据包给服务端,通过自定的接口函数处理完成后,通过公用消息通道UserMessageChannel再分发给不同的实际userid

3.一个用户ID只能在一个客户端中使用,如果在另一个客户端中使用,代码中会将上个用户强制登出

4.客户端定时发送心跳包给服务端,超过1min socket上没有数据包,会认为客户端已经挂起,会强制注销连接,为了节省服务器的开销

5.客户端发送数据包的时候,tcp会自动将多次的请求数据包(小包)合并发送到服务器,服务器无法识别合并到一起的数据包,所以需要增加拆包和封包的功能,避免TCP粘包..
其他的解决方案:比如固定数据包大小,或者每个数据包增加一个特殊字符,作为分割符号
```


代码拆解: 通过通信的方式进行共享内存
UserEnterChannel UserLeaveChannel 两个通道用于用户的信息的增加和删除
UserMessageChannel 公用消息分发通道,通过更改 MsgData中的ID,可以将消息,发给不同的userid,默认是发给自己
quitChannel 检测连接是否活动,不活动,强制关闭连接

 每个用户注册的时候,会创建一个user.MsgChannel,当需要给用户发消息的时,只需要user.MsgChannel中写入数据.func UserChannel(user *cache.UserInfo) 负责检测用户的消息通道...


 mapfunc = make(map[int]CommonInterface) 用于存放 实际的功能和funcid的映射关系...实际的业务代码,只需要 实现CommonInterface 接口,然后,添加到map[int]CommonInterface 中,即可