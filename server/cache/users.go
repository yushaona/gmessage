/*
	连接缓存 --  维护用户ID和连接句柄
*/
package cache

import (
	"fmt"
	"net"
	"time"
)

var (
	UserEnterChannel   = make(chan *UserInfo)
	UserLeaveChannel   = make(chan *UserInfo)
	UserMessageChannel = make(chan *MsgData, 1024)

	usersMap = make(map[string]*UserInfo)
)

type UserInfo struct {
	UserID     string
	Conn       net.Conn
	SessionID  string
	MsgChannel chan string
}

func CreateUser(userid string, c net.Conn) (user *UserInfo) {
	return &UserInfo{
		UserID:     userid,
		Conn:       c,
		SessionID:  time.Now().Format("2006-01-02-15:04:05"),
		MsgChannel: make(chan string, 10),
	}
}

//UserCache 缓存用户信息
func UserCache() {
	for {
		select {
		case user := <-UserEnterChannel:
			//强制踢出前一个
			if u, isok := usersMap[user.UserID]; isok {
				u.Conn.Close()
			}
			usersMap[user.UserID] = user
		case user := <-UserLeaveChannel:
			fmt.Println(user.UserID, "退出")
			user.Conn.Close()
			close(user.MsgChannel)

			if u, isok := usersMap[user.UserID]; isok {
				if user.SessionID == u.SessionID {
					delete(usersMap, user.UserID)
					fmt.Println("enter")
				}
			}

		case msgdata := <-UserMessageChannel: //msgdata 接受者id + 数据包data
			{
				if user, isok := usersMap[msgdata.ID]; isok {
					user.MsgChannel <- msgdata.Data
				}
			}
		}
	}
}

// //Keeplive 更新时间
// func Keeplive(id string) {
// 	if t, ok := connCache[id]; ok {
// 		t.lasttime = time.Now()
// 	}
// }

// func GetConn(id string) net.Conn {
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	if t, ok := connCache[id]; ok {
// 		return t.conn
// 	}
// 	return nil
// }

// func DelConn(id string) {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	connCache[id].conn.Close()
// 	close(connCache[id].datachannel)
// 	delete(connCache, id)
// }

// func SetConn(id string, c net.Conn) {
// 	lock.Lock()
// 	if t, ok := connCache[id]; ok {
// 		if t.conn != c {
// 			t.conn = c
// 		}
// 		t.lasttime = time.Now()
// 	} else {
// 		temp := &connData{
// 			id:       id,
// 			conn:     c,
// 			lasttime: time.Now(),
// 		}
// 		connCache[id] = temp
// 	}
// 	lock.Unlock()
// 	return
// }

// func GetNum() (result int) {
// 	return len(connCache)
// }

// // 存储用户ID和用户的连接句柄
// var (
// 	connCache map[string]*connData // 连接缓冲池
// 	once      sync.Once
// 	lock      sync.RWMutex
// )

// func init() {
// 	connCache = make(map[string]*connData)
// }

// // 定时检测 连接的有效性 -- 无效的连接清理出缓存
// func CheckVaild() {
// 	once.Do(func() {
// 		fmt.Println("CheckVaild start")
// 		ticker := time.NewTicker(time.Second * 10)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				lock.RLock()
// 				for k, v := range connCache {

// 					if time.Now().Sub(v.lasttime) > time.Second*10 {
// 						//30s连接时间没有更新,需要从缓存中清理掉
// 						delete(connCache, k)
// 						v.conn.Close()
// 					}
// 				}
// 				lock.RUnlock()
// 			}

// 			fmt.Printf("当前客户端数量:%d \n", GetNum())
// 		}
// 	})
// }
