/*
	接受客户端的连接请求
*/

package handle

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/yushaona/gmessage/server/job"

	"github.com/yushaona/gmessage/server/cache"

	"github.com/yushaona/gjson"
)

//Accept 等待客户端的连接
func Accept(listen net.Listener) {

	go cache.UserCache()
	//go HandleMsgQueue()
	for {
		conn, err := listen.Accept() //tcp三次握手完成,可以进行通信了
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go execute(conn)
	}

}

func HandleData(param *gjson.GJSON) (result gjson.GJSON, err error) {
	funcid := param.GetInt("funcid")
	if funcid == 0 {
		return result, fmt.Errorf("%s", "funcid is not exist")
	}
	return job.DoJob(funcid, param)

}

/*
消息数据包,主要包含该三个字段

id 也就是处理完成,发送给哪个用户ID,现在就是自己发给自己
funcid 执行哪个函数
data 数据


{
	"funcid":300,
	"userid":"462262902976417792",
	"session":"1",
	"macaddr":"111231321-"
	"type":0
}

*/
func HandleMsgQueue() {
	for {
		select {
		case msgData := <-cache.MsgQueue:
			//result := "[" + msgData.ID + "]===" + msgData.Data // 将请求的数据包发给id指定的用户
			var param gjson.GJSON
			param.Load(msgData.Data)
			//fmt.Println("result", result)
			var msgResult gjson.GJSON
			msgResponse, err := HandleData(&param) //
			fmt.Println(msgResponse.ToString())
			if err != nil {
				msgResult.SetInt("code", 0)
				msgResult.SetString("info", err.Error())
			} else {
				msgResult.SetInt("code", 1)
				msgResult.SetString("info", "ok")
				msgResult.SetObject("data", msgResponse)
			}
			//cache.GetConn(msgData.ID).Write([]byte(msgResult.ToString()))
		}
	}
}
func execute(c net.Conn) {

	var user *cache.UserInfo = nil
	var quitChannel = make(chan struct{})

	//var updateTime time.Time
	buf := make([]byte, 1024) // 单次接收的数据流
	for {
		n, err := c.Read(buf)
		if err != nil || err == io.EOF {
			fmt.Println("1111111111111111111")
			break
		}

		data := string(buf[:n])
		var main gjson.GJSON
		err = main.Load(data)

		if err != nil {
			fmt.Println(err)
			continue
		}

		/*
			{
				"userid":"462262902976417792",

				"funcid":300,
				"session":"1",
				"macaddr":"111231321-"

			}
		*/
		funcid := main.GetInt("funcid")

		switch funcid {
		case 1: //用户登录
			{
				userid := main.GetString("userid")
				if userid == "" {
					var obj gjson.GJSON
					obj.SetInt("code", 0)
					obj.SetString("info", "userid不能为空")
					c.Write([]byte(obj.ToString()))
					continue
				}
				if user == nil { // 一个socket只能初始化一次用户登录 -- 相当于当前的socket就和user绑定到了一起,有状态的TCP连接
					passwd := main.GetString("passwd")
					if passwd == "123456" {
						user = cache.CreateUser(userid, c)
						//updateTime = time.Now()
						cache.UserEnterChannel <- user
						var obj gjson.GJSON
						obj.SetInt("code", 1)
						obj.SetString("info", "ok")
						c.Write([]byte(obj.ToString()))
						go UserChannel(user)
						go UserVaild(c, quitChannel)
					} else {
						var obj gjson.GJSON
						obj.SetInt("code", 0)
						obj.SetString("info", "密码错误")
						c.Write([]byte(obj.ToString()))
					}
				}
			}
		default:
			{
				if user == nil {
					var obj gjson.GJSON
					obj.SetInt("code", 0)
					obj.SetString("info", "请登录")
					c.Write([]byte(obj.ToString()))
				} else {
					if funcid != 2 { // funcid=2 表示心跳包
						var msgResult gjson.GJSON
						msgResponse, err := HandleData(&main) //
						fmt.Println(msgResponse.ToString())
						if err != nil {
							msgResult.SetInt("code", 0)
							msgResult.SetString("info", err.Error())
						} else {
							msgResult.SetInt("code", 1)
							msgResult.SetString("info", "ok")
							msgResult.SetObject("data", msgResponse)
						}

						//表示处理后的结果,想要发给哪个用户 -- 使用一个公用的通道,接受数据,然后分发给对应的用户的通道
						cache.UserMessageChannel <- &cache.MsgData{ID: user.UserID, Data: msgResult.ToString()}
					}
					if funcid == 2 {
						fmt.Println("心跳包")
					}
					quitChannel <- struct{}{}
				}
			}
		}
	}

	//用户注销
	if user != nil {
		cache.UserLeaveChannel <- user
	}
}

//检查链接的有效性
func UserVaild(c net.Conn, quit <-chan struct{}) {
	d := 1 * time.Minute
	t := time.NewTimer(d) // 10s内客户端没有消息发送,会强制将用户下线

	for {
		select {
		case <-t.C:
			{
				c.Close()
			}
		case <-quit:
			{
				t.Reset(d)
			}
		}
	}
}

//异步发送数据
func UserChannel(user *cache.UserInfo) { //每个用户的数据通道
	for {
		select {
		case datastr := <-user.MsgChannel:
			user.Conn.Write([]byte(datastr))
		}
	}
}
