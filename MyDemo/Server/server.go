//Server

package server

import (
	"strconv"
	"MyDemo/Proto"
	"MyDemo/Room/SetRand"
	"MyDemo/Room"
	"MyDemo/User"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var AllRoom []*Room.Room
var RoomNum = 0
var count = 1
//结束标记
var done = make(chan int)

func test() {
	fmt.Println("This is server.go")
}


//建立连接
func GetConnection2() {
	//建立连接
	netListen, err := net.Listen("tcp", "192.168.132.211:1236")
	CheckError(err)
	defer netListen.Close()
	Log("服务器版本2018/8/29 16:00")
	
	AllRoom = Room.InitRoom()
  	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		
		user := User.NewUser()
		user.UserConn = conn
		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnectionRead(conn, user)
		go UserTempChat(user)
		count++
		log.Println("count =", count)
	}
}

//实时发送房间信息
func SendRoomInfo(user *User.User) {
	temp := make([]byte, 1024)
	for {
		select {
			case <- user.RoomInfoDone: 
			{
				log.Println("房间信息停止更新")
				break
				//当房间人数达到2时，会自动结束房间信息更新
			}
		}
		data := Room.RoomInfo(AllRoom)
		data = append(data, temp...)
		
		user.UserConn.Write(data[:1024])
		time.Sleep(time.Second * 1)
	}
}


//第一次连接是发送给客户端的名称
func SetUser(conn net.Conn, id int, name string) {
	data := make([]byte, 0)
	data = append(data, Proto.IntToBytes(9)...)
	data = append(data, Proto.IntToBytes(id)...)
	data = append(data, Proto.IntToBytes(1001)...)
	data = append(data, []byte(name)...)
	
	temp := make([]byte, 1024)
	data = append(data, temp...)
	conn.Write(data[:1024])
}


/*
	将随机点进行打包发送，具体逻辑写完。
*/
func SendRandPoint2(room *Room.Room) {
	temp := make([]byte, 1024)
	for i:= 0; i < 3; i++ {
		log.Println("发送随机坐标点")
		
		select {
			case <- room.RandPointDone: {
				room.RandPointDone <- true
				fmt.Println("随机发送坐标点结束")
				break
			}
			default:
		}
		
		//判断是否为刚开局
		if len(room.RoomUserA.Data) != 0 && len(room.RoomUserB.Data) != 0 {
			data := User.UpdateStatus(room.RoomUserA, room.RoomUserB)
			//发送剩余点信息
			User.SendClear(room.RoomUserA, room.RoomUserB, data)
			//清除所有玩家数据
			User.ClearPointInfo(room.RoomUserA, room.RoomUserB)
		}
		time.Sleep(time.Second * 3)
		RandPoint := SetRand.InitPoint(i)

		messege, flag := SetRand.PackPoint(RandPoint, i)
		if !flag {
			log.Println("Send RandPoint Wrong")
		}
		
		messege = append(messege, temp...)
		
		//time.Sleep(time.Millisecond * 500)
		//新添加进协程进行处理
		room.RoomUserA.InitPointMap(RandPoint, i)
		room.RoomUserB.InitPointMap(RandPoint, i)
		log.Println("开始发送")
		room.RoomUserA.UserConn.Write(messege[:1024])
		log.Println("A发送成功")
		room.RoomUserB.UserConn.Write(messege[:1024])
		log.Println("B发送成功")
		time.Sleep(time.Second * 57)
	}
}



//客户端传输消息的获取
func handleConnectionRead(conn net.Conn, Player *User.User) {

	//接收解包
	log.Println("Start Reading!")
	defer conn.Close()
	//定义在外面会有数据不清空的现象，需要做一个清除的处理
	tempmes := make([]byte, 0)
	mes := make([]byte, 1024)
	for {
		buffer := mes
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			AllRoom[Player.RoomId].RandPointDone <- true
			AllRoom[Player.RoomId] = Room.NewRoom(Player.RoomId)
			log.Println("Send Finished!")
			
			count--
			RoomNum--
			conn.Close()
			break
		}

		//将获取到的消息写入玩家消息通道
		tempmes = Proto.DePackage2(append(tempmes, buffer[:n]...), Player.MessegeWrite)
		
	}
}

//结束房间
func RoomDone(room *Room.Room) {
	time.Sleep(time.Second * 500)
	room.RandPointDone <- true
	//room.RoomUserA.UserConn = nil
	//room.RoomUserB.UserConn = nil
	room = Room.NewRoom(room.RoomUserA.RoomId)
}

//PlanB 启动B计划

func handleConnectionWrite2(conn net.Conn, room *Room.Room) {
	temp := make([]byte, 1024)
	if conn == room.RoomUserA.UserConn {
		for {
			select {
			//对玩家1的消息读取并广播
			case messege, ok := <-room.RoomUserA.MessegeWrite:
				{
					if !ok {
						log.Println("End")
						break
					}
					//解析数据
					data, flag, cmd, textlen := Proto.DePackage(messege)
					data = append(data, temp...)
					switch cmd {
					case 1002:
						{
							log.Println("碰球逻辑开始")
							if flag {
								pointname := string(messege[12:textlen + 4])
								point := User.UseNameFindPoint(pointname, room.RoomUserA.Data)
								log.Println("球的名字:", pointname)
								//新加协程处理
								room.RoomUserA.UpdatePoint(point, room.RoomUserB)
							}
						}
					case 1005:
						{
							if flag {
								room.RoomUserB.UserConn.Write(data[:1024])
								//log.Println("messege1 :", data)
							} else {
								log.Println("UserA 解析指令出错，不进行任何处理")
							}
						}
					case 3008 :
						{
								if flag {
								room.RoomUserB.UserConn.Write(data[:1024])
								//log.Println("messege1 :", data)
							} else {
								log.Println("UserA 解析指令出错，不进行任何处理")
							}
						}
					case 4008 :
						{
							room.RoomUserB.UserConn.Write(data[:1024])
						}
					}
					//log.Println("User1 Finished")
				}
			}
		}
	} else if room.RoomUserB.UserConn == conn {
		for {
			select {
			case messege, ok := <-room.RoomUserB.MessegeWrite:
				{
					if !ok {
						log.Println("End")
						break
					}
					//解析数据
					data, flag, cmd, textlen := Proto.DePackage(messege)
					data = append(data, temp...)
					switch cmd {
					case 1002:
						{
							log.Println("碰球逻辑开始")
							if flag {
								pointname := string(messege[12:textlen + 4])
								point := User.UseNameFindPoint(pointname, room.RoomUserB.Data)
								log.Println("球的名字:", pointname)
								//更新状态会时刻向客户端发送更新后的信息
								//新加协程处理
								go room.RoomUserB.UpdatePoint(point, room.RoomUserA)
							}
						}
					case 1005:
						{
							if flag {
								room.RoomUserA.UserConn.Write(data[:1024])
								//log.Println("messege1 :", data)
							} else {
								log.Println("UserB 解析指令出错，不进行任何处理")
							}
						}
						case 3008 :
						{
							if flag {
								room.RoomUserA.UserConn.Write(data[:1024])
								//log.Println("messege1 :", data)
							} else {
								log.Println("UserA 解析指令出错，不进行任何处理")
							}
						}
						case 4008 :
						{
							room.RoomUserA.UserConn.Write(data[:1024])
						}
						
					}
				}
			}
		}
	}
}

//创建加入房间
func CreInsRoom(id, cmd int) []byte {
	mes := make([]byte, 0)
	userid := Proto.IntToBytes(id)
	usercmd := Proto.IntToBytes(cmd)
	str := "Fail"
	data := []byte(str)
	Len := Proto.IntToBytes(len(data) + 8)
	mes = append(mes, Len...)
	mes = append(mes, userid...)
	mes = append(mes, usercmd...)
	mes = append(mes, data...)
	return mes
}



//用户临时会话
func UserTempChat(user *User.User) {
	log.Println("开启临时会话")
	for {
		select {
			//临时会话结束
			case ok := <- user.TempChat: {
				user.TempChat <- true
				log.Println("临时会话结束", ok)
				break
			}
			//读取临时消息
			case messege := <-user.MessegeWrite: {
				data, _, cmd, texLen := Proto.DePackage(messege)
				switch cmd {
					//新建房间
					case 7001: {
						
						id := string(data[12:texLen + 4])
						id2, _ := strconv.Atoi(id)
						log.Println("id:", id2)
						if AllRoom[id2].NowPoeple == 0 {
							user.RoomId = id2
							//添加玩家到该房间
							AllRoom[id2].AddPlayer(user)
							//结束发送房间状态信息
							messege := CreInsRoom(0, 7001)
							data := make([]byte, 1024)
							messege = append(messege, data...)
							user.UserConn.Write(messege[:1024])
							//加入房间对战
							go handleConnectionWrite2(user.UserConn, AllRoom[user.RoomId])
							user.TempChat <- true
							user.RoomInfoDone <- true
						} else {
							log.Println("创建失败，房间已经存在")
							messege := CreInsRoom(0, 7002)
							data := make([]byte, 1024)
							messege = append(messege, data...)
							user.UserConn.Write(messege[:1024])
						}
					}
					//加入房间
					case 8001: {
						
						id := string(data[12:texLen + 4])
						id2, _ := strconv.Atoi(id)
						log.Println("id:", id2)
						if AllRoom[id2].NowPoeple == 1 {
							user.RoomId = id2
							//添加玩家到该房间
							AllRoom[id2].AddPlayer(user)
							
							messege := CreInsRoom(0, 8001)
							data := make([]byte, 1024)
							messege = append(messege, data...)
							user.UserConn.Write(messege[:1024])
							//开启对战通知
							temp := StartBattle()
							temp = append(temp, data...)
							AllRoom[user.RoomId].RoomUserA.UserConn.Write(temp[:1024])
							AllRoom[user.RoomId].RoomUserB.UserConn.Write(temp[:1024])
							//分配角色
							SetUser(AllRoom[user.RoomId].RoomUserA.UserConn, count, "B")
							SetUser(AllRoom[user.RoomId].RoomUserB.UserConn, count, "A")
							//开启房间对战
							go handleConnectionWrite2(user.UserConn, AllRoom[user.RoomId])
							go StartRoom(AllRoom[user.RoomId])
							//结束临时会话
							user.TempChat <- true
						} else {
							log.Println("加入失败，房间已满")
							messege := CreInsRoom(0, 8002)
							data := make([]byte, 1024)
							messege = append(messege, data...)
							user.UserConn.Write(messege[:1024])
						}
					}	
				}
			}
		}
	}
}

//战斗开始
func StartBattle() []byte {
	id := 0
	cmd := 8009
	mes := make([]byte, 0)
	userid := Proto.IntToBytes(id)
	usercmd := Proto.IntToBytes(cmd)
	str := "ok"
	data := []byte(str)
	Len := Proto.IntToBytes(len(data) + 8)
	mes = append(mes, Len...)
	mes = append(mes, userid...)
	mes = append(mes, usercmd...)
	mes = append(mes, data...)
	return mes
}

//开启房间对战
func StartRoom(room *Room.Room) {
	cmd := 8010
	str := "ok"
	mes := make([]byte, 0)
	data := make([]byte, 1024)
	mes = append(mes, Proto.IntToBytes(10)...)
	mes = append(mes, Proto.IntToBytes(2)...)
	mes = append(mes, Proto.IntToBytes(cmd)...)
	mes = append(mes, []byte(str)...)
	mes = append(mes, data...)
	time.Sleep(1 * time.Second)
	room.RoomUserA.UserConn.Write(mes[:1024])
	room.RoomUserB.UserConn.Write(mes[:1024])
	go SendRandPoint2(room)
	go RoomDone(room)
}


//出错检查
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//简单的输出信息
func Log(v ...interface{}) {
	fmt.Println(v...)
}
