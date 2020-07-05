//Room

package Room

import (
	"strconv"
	"log"
	"MyDemo/User"
	"MyDemo/Proto"
	"fmt"
	"sync"
)

type Room struct {
	//房间名
	RoomId    int
	//房间容纳最大人数
	MaxPeople int
	//房间现在人数
	NowPoeple int
	//房间结束标记
	RoomDone chan bool
	//更新房间信息结束标记
	InfoDone chan bool
	//互斥锁
	RandPointDone chan bool
	mtx       sync.Mutex
	//黑暗角色
	RoomUserA *User.User
	//光明角色
	RoomUserB *User.User
}

func NewRoom(id int) *Room {
	return &Room{
		RoomId:    id,
		MaxPeople: 2,
		NowPoeple: 0,
		RoomDone: make(chan bool),
		InfoDone: make(chan bool),
		RoomUserA: User.NewUser(),
		RoomUserB: User.NewUser(),
	}
}

func test() {
	fmt.Println("This is Room.go")
}

//房间初始化
func InitRoom() []*Room {
	room := make([]*Room, 100)
	for i := 0; i < 100; i++ {
		room[i] = NewRoom(i)
	}
	return room
}

//添加玩家进房间
func (room *Room) AddPlayer(user *User.User) {
	room.NowPoeple++
	if room.NowPoeple == 1 {
		room.RoomUserA = user
		//room.RoomUserA.UserConn = conn
		room.RoomUserA.UserID = 1
	} else if room.NowPoeple == 2 {
		room.RoomUserB = user
		//room.RoomUserB.UserConn = conn
		room.RoomUserB.UserID = 2
	} else {
		log.Println("房间人数已满")
	}
	log.Println("现在房间人数:", room.NowPoeple)
}

//更新房间现在信息
func RoomInfo(room []*Room) []byte {
	cmd := 2002
	var messege string
	for i := 0; i < 6; i++ {
		if i == 5 {
			messege = messege + strconv.Itoa(room[i].RoomId) + ","
			messege = messege + strconv.Itoa(room[i].NowPoeple)
			break
		}
		messege = messege + strconv.Itoa(room[i].RoomId) + ","
		messege = messege + strconv.Itoa(room[i].NowPoeple) + ","
	}

	data, _ := Proto.EnPackage([]byte(messege), len(messege) + 8, 0, cmd)
	return data
}