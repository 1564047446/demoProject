//Proto

/*
消息协议模块
主要将消息进行解析，分别执行不同操作
将得到的信息进行压缩处理
最后打包发送给客户端
*/

//1001 分发用户名
//1005 玩家移动处理
//1099 发送随机坐标点
//1006 状态点的获取
//1003 发送新的状态点

package Proto

//package main

import (
	//"MyDemo/User"
	"bytes"
	"encoding/binary"
	"log"
)

func test() {
	log.Println("This is Proto.go ")

}

//所有涉及到Cmd的命令都要注意一下，不论是获取原值还是更新返回
const (
	//字节总长度占4个字节
	TextLength = 4
	//用户ID占四个字节
	UserId = 4
	//操作指令占4个字节
	CommandLength = 4
)

//压包
func EnPackage(message []byte, Len int, Id int, Cmd int) ([]byte, bool) {
	data := make([]byte, 0)
	var textlen int
	textlen = Len
	//将文本长度添加进byte数组   前四个字节
	data = append(data, IntToBytes(textlen)...)

	//将用户ID添加进byte数组  第5-8个字节
	data = append(data, IntToBytes(Id)...)

	//将新的命令返回
	//将CMD添加进byte数组 第9-12个字节
	data = append(data, IntToBytes(Cmd)...)

	//最后将消息添加进数组中返回
	data = append(data, message...)

	//将发送信息完整的字节流输出，方便后续DEBUG
	//log.Println("这里是压包:", data)
	return data, true
}

/**********************************
*	以下模块彻底完成，				  *
***********************************/

//解包
func DePackage2(message []byte, read chan []byte) ([]byte) {
	length := len(message)
	var i int
	for i = 0; i < length; i++ {
		if length < i + 8 {
			log.Println("长度不够")
			break
		}
		
		if BytesToInt(message[i:i + 4]) < 999 {
			//包长度获取
			meslen := BytesToInt(message[i:i + 4])
			if length < i + meslen + 4 {
				break
			}
			data := message[i : i + meslen + 4]
			read <- data
			i += meslen + 3
		}
	}
	if i == length {
		return make([]byte, 0)
	}
	
	return message[i:]

}



//解包
func DePackage(message []byte) ([]byte, bool, int, int) {

	//获取
	TextLen := BytesToInt(message[:4])
	//玩家ID获取
	Id := BytesToInt(message[4:8])
	//指令获取
	Cmd := BytesToInt(message[8:12])
	//解析命令
	switch Cmd {
	//正常移动
	case 1005:
		{
			//log.Println("User is moving......")
			//将移动的消息进行打包处理
			move, ok := EnPackage(message[12:], TextLen, Id, Cmd)
			//将移动信息返回并发送客户端
			if ok {
				//log.Println("move...", move)
				return move, true, Cmd, TextLen
			} else {
				log.Println("Wrong moving Cmd!")
				return move, false, Cmd, TextLen
			}
		}
	//标记站点
	//新加入改动 解析点坐标时，返回长度+4防止出现粘包
	case 1002:
		{
			log.Println("User Get Point")
			return message[:TextLen + 4], true, Cmd, TextLen
		}
	//
	case 3008:
		{
			log.Println("技能释放")
			//将移动的消息进行打包处理
			move, ok := EnPackage(message[12:], TextLen, Id, Cmd)
			//将移动信息返回并发送客户端
			if ok {
				log.Println("放技能..", move)
				return move, true, Cmd, TextLen
			} else {
				log.Println("技能释放失败!")
				return move, false, Cmd, TextLen
			}
		}
	case 7001:
	{
		log.Println("创建房间")
		id, _ := EnPackage(message[12:], TextLen, Id, Cmd)
		return id, true, Cmd, TextLen
	}
	case 8001: 
	{
		log.Println("加入房间")
		id, _ := EnPackage(message[12:], TextLen, Id, Cmd)
		return id, true, Cmd, TextLen	
	}
	case 4008:
	{
		log.Println("技能：黑暗之墙")
		id, _ := EnPackage(message[12:], TextLen, Id, Cmd)
		return id, true, Cmd, TextLen	
	}
	//解析到错误指令的处理
	default:
		{
			log.Println("Something wrong in Parse package")
			data := make([]byte, 100)
			return data, false, Cmd, TextLen

		}
	}
	//等到逻辑完善后需要删除
	log.Println("没有解析到正确的Command")
	data := make([]byte, 0)
	return data, false, Cmd, TextLen
}

//整形转换成字节，采用小端机型进行接收和发送
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形，采用小端机型进行接收和发送
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int(x)
}

//测试模块专用，不需要时可以直接注释掉
