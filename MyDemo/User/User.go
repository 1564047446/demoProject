//User

package User

//package main

import (
	"time"
	"sync"
	"MyDemo/Proto"
	"MyDemo/Room/GetArea"
	"MyDemo/Room/SetRand"
	"MyDemo/User/GetSkill"
	"container/list"
	"fmt"
	"log"
	"net"
	"strconv"
)

type UserPoint struct {
	X      int
	Y      int
	Z      int
	Name   string
	Status int
}

func NewUserPoint(x int, y int, z int, Name string) UserPoint {
	return UserPoint{
		X:      x,
		Y:      y,
		Z:      z,
		Name:   Name,
		Status: 0,
	}
}

//代替被消亡的点的名字
const (
	Instead = "XXX"
)

type User struct {
	//用户ID
	UserID int
	//用户连接状态
	UserConn net.Conn
	//用户求得最终面积
	UserArea float32
	//用户读取信息保存
	MessegeRead []byte
	//光球点坐标信息
	PointMap map[UserPoint]int
	//玩家站点个数
	PointNum int
	//上一次经过的点
	LastPoint UserPoint
	//现在的点
	CurrentPoint UserPoint
	//光球坐标链表
	PointList *list.List
	//玩家行走路径
	PointPath *list.List
	//消息管道，用来协程间的通信
	MessegeWrite chan []byte
	//用户结束临时会话标记
	TempChat chan bool
	//用户结束房间信息状态更新
	RoomInfoDone chan bool
	Data         []UserPoint
	//房间号
	RoomId int
	//互斥锁
	MyLock sync.Mutex
	//玩家当前技能
	Skills map[string]int
}

func NewUser() *User {
	return &User{
		UserID:       0,
		UserConn:     nil,
		UserArea:     0.0,
		PointMap:     make(map[UserPoint]int),
		MessegeRead:  make([]byte, 1024),
		Data:         make([]UserPoint, 0),
		PointList:    list.New(),
		PointPath:    list.New(),
		PointNum:     0,
		LastPoint:    NewUserPoint(0, 0, 0, Instead),
		CurrentPoint: NewUserPoint(0, 0, 0, Instead),
		MessegeWrite: make(chan []byte),
		TempChat: make(chan bool),
		RoomInfoDone: make(chan bool),
		RoomId : 0,
		Skills: make(map[string]int),
	}
}

//打印玩家路径
func (u *User) PrintPath() {
	log.Println("打印玩家", u.UserID, "的路径")
	for p := u.PointPath.Front(); p != nil; p = p.Next() {
		log.Println(p.Value)
	}
}

//打印玩家获得点
func (u *User) PrintList() {
	log.Println("打印玩家", u.UserID, "的点")
	for p := u.PointList.Front(); p != nil; p = p.Next() {
		log.Println(p.Value)
	}
}

//打印玩家map信息
func (u *User) PrintMap() {
	log.Println("打印玩家", u.UserID, "MAp")
	for key, value := range u.PointMap {
		log.Println(key, value)
	}
}

//初始化玩家坐标点信息为0
func (user *User) InitPointMap(point *SetRand.Point, x int) {
	user.MyLock.Lock()
	user.PointNum = 0
	for i := 0; i < SetRand.Num[x]; i++ {
		tempPoint := NewUserPoint(point.X[i], int(point.Y), point.Z[i], "Marker_"+strconv.Itoa(i))
		user.Data = append(user.Data, tempPoint)
		user.PointMap[tempPoint] = 0
	}
	user.MyLock.Unlock()
}

//通过坐标名字找到坐标
func UseNameFindPoint(name string, point []UserPoint) UserPoint {
	var temp UserPoint
	for i := 0; i < len(point); i++ {
		if point[i].Name == name {
			temp = point[i]
			break
		}
	}
	return temp
}

func ClearPointInfo(a *User, b *User) {
	
	
	//清除A列表
	for point := a.PointList.Front(); point != nil; {
		next := point.Next()
		a.PointList.Remove(point)
		point = next
	}
	//清除A路径
	for point := a.PointPath.Front(); point != nil; {
		next := point.Next()
		a.PointPath.Remove(point)
		point = next
	}
	//清除A的Map
	for key, _ := range a.PointMap {
		delete(a.PointMap, key)
	}

		//清除B列表
	for point := b.PointList.Front(); point != nil; {
		next := point.Next()
		b.PointList.Remove(point)
		point = next
	}
	//清除B路径
	for point := b.PointPath.Front(); point != nil; {
		next := point.Next()
		b.PointPath.Remove(point)
		point = next
	}
	//清除B的Map
	for key, _ := range b.PointMap {
		delete(b.PointMap, key)
	}
	a.Data = make([]UserPoint, 0)
	b.Data = make([]UserPoint, 0)
	a.LastPoint.Name = Instead
	b.LastPoint.Name = Instead
	//a，b解锁
	
}

//消亡状态更新并获得环形面积
//并且发送生成坐标点顺序
func DestroyStatus(user *User, userB *User) float32 {
	//新建坐标点集合
	temp := 0
	point := GetArea.NewPoint()
	for key, value := range user.PointMap {
		if value == 2 {
			user.PointMap[key] = 3
			log.Println(key.Name)
			temp++
			//建立点集合
			point.BuildPoint(float32(key.X), float32(key.Z), key.Name)
		}
	}
	
	point.PrintPoint()
	area, newpoint := GetArea.ParsePointInLine(point)
	//加上互斥锁
	
	user.MyLock.Lock()
	user.UserArea += area * (1 + (float32(temp) - 3))
	user.MyLock.Unlock()
	data, ok := GetArea.PackPoint(newpoint, 1200)
	//发送成环坐标
	if ok {
		datat := make([]byte, 1024)
		data = append(data, datat...)
		user.UserConn.Write(data[:1024])
		userB.UserConn.Write(data[:1024])
	}
	data2 := GetArea.PackArea(int(user.UserArea), int(userB.UserArea))
	data3 := GetArea.PackArea(int(userB.UserArea), int(user.UserArea))
	datat1 := make([]byte, 1024)
	data2 = append(data2, datat1...)
	data3 = append(data3, datat1...)
	user.UserConn.Write(data2[:1024])
	userB.UserConn.Write(data3[:1024])
	return area
}

//实时更新分数
func GetMark(user *User, userB *User) {
	for {
		t := make([]byte, 1024)
		data2 := GetArea.PackArea(int(user.UserArea), int(userB.UserArea))
		data3 := GetArea.PackArea(int(userB.UserArea), int(user.UserArea))
		data2 = append(data2, t...)
		data3 = append(data3, t...)
		user.UserConn.Write(data2[:1024])
		userB.UserConn.Write(data3[:1024])
		time.Sleep(time.Second * 1)
	}
}

//更新玩家球的坐标信息
func (user *User) UpdatePoint(NowPoint UserPoint, userB *User) {
	//上锁
	//user.MyLock.Lock()
	if user.LastPoint.Name == NowPoint.Name {
		log.Println("玩家在同一个点碰撞不进行处理")
		return 
	}
	log.Println("更新光球状态")
	//第一次走点的特殊处理
	if user.LastPoint.Name == Instead {
		log.Println("第一次走点", NowPoint.Name)
		for key, _ := range user.PointMap {
			if key.Name == NowPoint.Name {
				user.PointMap[key] = 1
				break
			}
		}
	}
	name := make([]string, 2)
	name[0], name[1] = user.LastPoint.Name, NowPoint.Name
	flag1 := false
	for p := user.PointList.Front(); p != nil; p = p.Next() {
		if p.Value == NowPoint.Name {
			flag1 = true
			break
		}
	}
	//添加进点集合
	if !flag1 {
		user.PointList.PushBack(NowPoint.Name)
	}
	flag1 = false
	//插入路径
	tempP := user.PointPath.Back()
	if tempP == nil {
		user.PointPath.PushBack(name)
	} else {
		teampval := tempP.Value.([]string)
		if teampval[0] != NowPoint.Name {
			user.PointPath.PushBack(name)
		}
	}
	//道具的获取判断
	tempP = user.PointPath.Back()
	tempval := tempP.Value.([]string)
	pn := GetArea.NewPoint()
	if tempval[0] != Instead {
		for val, _ := range user.PointMap {
			//找到走过的两个点并判断是否生成道具
			if val.Name == tempval[0] || val.Name == tempval[1] {
				pn.BuildPoint(float32(val.X), float32(val.Z), val.Name)
			}
		}
		//判断是否有道具生成
		temp, stoneName := GetArea.CircleLine(pn)
		log.Println("石柱名称", stoneName)
		for i := 0; i < temp; i++ {
			skill := GetSkill.RandSkill()
			isHere := false
			//如果有发送道具名称
			if !isHere {
				data := GetArea.PackSkill(skill)
				data2 := make([]byte, 1024)
				data = append(data, data2...)
				user.UserConn.Write(data[:1024])
				user.Skills[skill] = 1
			} 
			
		}
		//道具获取表现
		if temp != 0 {
			lenStone := len(stoneName)
			userid := 1
			cmd := 1201
			q := make([]byte, 1024)
			data := make([]byte, 0)
			data = append(data, Proto.IntToBytes(lenStone + 8)...)
			data = append(data, Proto.IntToBytes(userid)...)
			data = append(data, Proto.IntToBytes(cmd)...)
			data = append(data, []byte(stoneName)...)
			user.UserConn.Write(append(data, q...))
			userB.UserConn.Write(append(data, q...))
		}
	}
	length := user.PointList.Len()
	//更新点状态
	if length >= 2 {
		for p := user.PointList.Front(); p != nil; p = p.Next() {
			tempname := p.Value.(string)
			cnt := 0
			for q := user.PointPath.Front(); q != nil; q = q.Next() {
				n := q.Value.([]string)
				if n[0] == Instead {
					continue
				}
				if n[0] == tempname || n[1] == tempname {
					cnt++
				}
			}
			if cnt > 2 {
				cnt = 2
			}
			for key, _ := range user.PointMap {
				if key.Name == tempname {
					user.PointMap[key] = cnt
				}
			}
		}
	}
	//更新上一次走过的节点
	user.LastPoint = NowPoint
	iscir := false
	for key, value := range user.PointMap {
		//开启成环判断
		if key.Name == name[1] && value == 2 {
			iscir = IsCircle(user)
			break
		}
	}
	if iscir {
		log.Println("Is a Circle")
		//更新状态并且发送坐标点,求得玩家获取的面积
		DestroyStatus(user, userB)
		
		//更新消亡点状态
		ChangeStatus(user, userB)
		if user.UserID % 2 == 1 {
			//发送消亡点之后的状态更新，cmd位置需要有所改变
			data3 := UpdateStatus(user, userB)
			SendStatus(user, data3, 1)
			SendStatus(userB, data3, 1)
			
		} else {
			data3 := UpdateStatus(userB, user)
			SendStatus(user, data3, 1)
			SendStatus(userB, data3, 1)
		}

	} else {
		log.Println("Not a Circle")
		if user.UserID % 2 == 1 {
			//发送消亡点之后的状态更新，cmd位置需要有所改变
			data3 := UpdateStatus(user, userB)
			SendStatus(user, data3, 1)
			SendStatus(userB, data3, 1)
			
		} else {
			data3 := UpdateStatus(userB, user)
			SendStatus(user, data3, 1)
			SendStatus(userB, data3, 1)
		}
	}
	
	//打印当前操作玩家路径
	for p := user.PointPath.Front(); p != nil; p= p.Next() {
		log.Println("路径:", p.Value)
	}

}

func PointMoment(user *User, userB *User) {
	for {
		if user.UserID % 2 == 1 {
			data3 := UpdateStatus(user, userB)
			SendStatus(user, data3, 1)
			SendStatus(userB, data3, 1)
		} else {
			data3 := UpdateStatus(user, userB)
			SendStatus(user, data3, 1)
			SendStatus(userB, data3, 1)
		}
		time.Sleep(time.Second * 1)
	}
}

//玩家成环后，改变两个玩家的状态
func ChangeStatus(a *User, b *User) {
	//获取消亡点
	tempname := make([]string, 0)
	for key, value := range a.PointMap {
		if value == 3 {
			tempname = append(tempname, key.Name)
		}
	}
	log.Println("消亡点名称", tempname)
	//根据消亡点删除B玩家map中的消亡点
	for i := 0; i < len(tempname); i++ {
		for key, _ := range b.PointMap {
			if key.Name == tempname[i] {
				_, ok := b.PointMap[key]

				if ok {
					delete(b.PointMap, key)
				}
				break
			}
		}
	}
	//根据消亡点更新B的List
	for i := 0; i < len(tempname); i++ {
		for p := b.PointList.Front(); p != nil; {
			next := p.Next()
			if tempname[i] == p.Value {
				b.PointList.Remove(p)
			}
			p = next
		}
	}
	//根据小网点更新B的Path
	for i := 0; i < len(tempname); i++ {
		for p := b.PointPath.Front(); p != nil; {
			next := p.Next()
			val := p.Value.([]string)
			if tempname[i] == val[0] || tempname[i] == val[1] {
				b.PointPath.Remove(p)
			}
			p = next
		}
	}
	//初始化B的Map
	for key, _ := range b.PointMap {
		b.PointMap[key] = 0
	}
	//根据新的Path和List更新B的Map
	for p := b.PointList.Front(); p != nil; p = p.Next() {
		tempname := p.Value.(string)
		cnt := 0
		for q := b.PointPath.Front(); q != nil; q = q.Next() {
			n := q.Value.([]string)
			if n[0] == Instead {
				continue
			}
			if n[0] == tempname || n[1] == tempname {
				cnt++
			}
		}
		if cnt > 2 {
			cnt = 2
		}
		for key, _ := range b.PointMap {
			if key.Name == tempname {
				b.PointMap[key] = cnt
			}
		}
	}

	//初始化A的所有信息
	for key, _ := range a.PointMap {
		a.PointMap[key] = 0
	}
	for key, _ := range a.PointMap {
		for i := 0; i < len(tempname); i++ {
			if key.Name == tempname[i] {
				delete(a.PointMap, key)
				break
			}
		}
	}
	for p := a.PointList.Front(); p != nil; {
		next := p.Next()
		a.PointList.Remove(p)
		p = next
	}
	for p := a.PointPath.Front(); p != nil; {
		next := p.Next()
		a.PointPath.Remove(p)
		p = next
	}
	a.LastPoint = NewUserPoint(0, 0, 0, Instead)

	//判断消亡的点是不是B最后走过的路径
	ok := false
	for i := 0; i < len(tempname); i++ {
		if b.LastPoint.Name == tempname[i] {
			ok = true
		}
	}
	if ok {
		p := b.PointPath.Back()
		if p == nil {
			b.LastPoint = NewUserPoint(0, 0, 0, Instead)
		} else {
			tm := p.Value.([]string)
			b.LastPoint = NewUserPoint(0, 0, 0, tm[1])
		}
	}
	
}

//发送状态数据
func SendStatus(a *User, p []PointStatus, flag int) {

		//发送数据
		var messege string
		dot := ","
		for i, _ := range p {
			if i == 0 {
				messege = messege + p[i].Name
				messege = messege + dot
				messege = messege + strconv.Itoa(p[i].Status)

			} else {
				messege = messege + dot
				messege = messege + p[i].Name
				messege = messege + dot
				messege = messege + strconv.Itoa(p[i].Status)
			}
		}
		//文本长度 + 8表示userid和cmd所占用的字节数
		TextLen := len(messege) + 8
		//发送指令
		Cmd := 1003
		if flag == 2 {
			Cmd += 1
		}
		//用户ID，其实就相当于对两个客户进行一个广播处理\

		UserId := 0
		//测试输出
		fmt.Println("string:", messege)
		//将string强转byte数组
		mes := []byte(messege)
		//暂时未考虑到ok = false的情况
		data, _ := Proto.EnPackage(mes, TextLen, UserId, Cmd)
		//发送数据
		//新添加协程处理
		data2 := make([]byte, 1024)
		data = append(data, data2...)
		a.UserConn.Write(data[:1024])
		//log.Println(p)
	
}

//发送清除数据
func SendClear(a *User, b *User, p []PointStatus) {

	if a.UserConn != nil && b.UserConn != nil {
		//发送数据
		var messege string
		dot := ","
		for i, _ := range p {
			if i == 0 {
				messege = messege + p[i].Name
				messege = messege + dot
				messege = messege + strconv.Itoa(p[i].Status)

			} else {
				messege = messege + dot
				messege = messege + p[i].Name
				messege = messege + dot
				messege = messege + strconv.Itoa(p[i].Status)
			}
		}
		//文本长度 + 8表示userid和cmd所占用的字节数
		TextLen := len(messege) + 8
		//发送指令
		Cmd := 1100
		//用户ID，其实就相当于对两个客户进行一个广播处理\

		UserId := 0
		//测试输出
		//fmt.Println("string:", messege)
		//将string强转byte数组
		mes := []byte(messege)
		//暂时未考虑到ok = false的情况
		data, _ := Proto.EnPackage(mes, TextLen, UserId, Cmd)
		//发送数据
		//新添加协程处理
		data2 := make([]byte, 1024)
		data = append(data, data2...)
		a.UserConn.Write(data[:1024])
		b.UserConn.Write(data[:1024])
		//log.Println(p)
	} else {
		if a.UserConn == nil && b.UserConn == nil {
			log.Println("There is no player here, Can't Send messege")
		} else if a.UserConn == nil {
			log.Println("PlayerA not here, Can't Send messege")
		} else {
			log.Println("PlayerA not here, Can't Send messege")
		}

	}
}

//光球坐标是否成环
func IsCircle(U *User) bool {
	cnt := 0
	for _, num := range U.PointMap {
		if num >= 2 {
			cnt++
		}
	}
	if cnt >= 3 {
		return true
	} else {
		return false
	}
}

//简单测试函数，无关紧要
func test(user *User) {
	fmt.Println("This is User import")
	//a := NewUserPoint(1, 2, 3, "aaa")

}

type PointStatus struct {
	Name   string
	Status int
}

func UpdateStatus(userA *User, userB *User) []PointStatus {
	pointset := make([]PointStatus, 0)
	//枚举10种状态
	// 0 0 0, 0 1 1, 0 2 2
	// 1 0 3, 1 1 4, 1 2 5
	// 2 0 6, 2 1 7, 2 2 8
	//其余状态为9
	for keyA, valueA := range userA.PointMap {
		switch valueA {
		case 0:
			{
				temp := FindPoint(userB, keyA.Name)
				switch temp {
				case 0:
					{
						point := PointStatus{Name: keyA.Name, Status: 0}
						pointset = append(pointset, point)
					}
				case 1:
					{
						point := PointStatus{Name: keyA.Name, Status: 1}
						pointset = append(pointset, point)
					}
				case 2:
					{
						point := PointStatus{Name: keyA.Name, Status: 2}
						pointset = append(pointset, point)
					}
				case 3:
					{
						point := PointStatus{Name: keyA.Name, Status: 9}
						pointset = append(pointset, point)
					}
				}
			}
		case 1:
			{
				flag := FindPoint(userB, keyA.Name)
				switch flag {
				case 0:
					{
						point := PointStatus{Name: keyA.Name, Status: 3}
						pointset = append(pointset, point)
					}
				case 1:
					{
						point := PointStatus{Name: keyA.Name, Status: 4}
						pointset = append(pointset, point)
					}
				case 2:
					{
						point := PointStatus{Name: keyA.Name, Status: 5}
						pointset = append(pointset, point)
					}
				case 3:
					{
						point := PointStatus{Name: keyA.Name, Status: 9}
						pointset = append(pointset, point)
					}
				}
			}
		case 2:
			{
				flag := FindPoint(userB, keyA.Name)
				switch flag {
				case 0:
					{
						point := PointStatus{Name: keyA.Name, Status: 6}
						pointset = append(pointset, point)
					}
				case 1:
					{
						point := PointStatus{Name: keyA.Name, Status: 7}
						pointset = append(pointset, point)
					}
				case 2:
					{
						point := PointStatus{Name: keyA.Name, Status: 8}
						pointset = append(pointset, point)
					}
				case 3:
					{
						point := PointStatus{Name: keyA.Name, Status: 9}
						pointset = append(pointset, point)
					}
				}
			}
		case 3:
			{
				flag := FindPoint(userB, keyA.Name)
				switch flag {
				case 0:
					{
						point := PointStatus{Name: keyA.Name, Status: 9}
						pointset = append(pointset, point)
					}
				case 1:
					{
						point := PointStatus{Name: keyA.Name, Status: 9}
						pointset = append(pointset, point)
					}
				case 2:
					{
						point := PointStatus{Name: keyA.Name, Status: 9}
						pointset = append(pointset, point)
					}
				case 3:
					{
						point := PointStatus{Name: keyA.Name, Status: 9}
						pointset = append(pointset, point)
					}
				}
			}
		}
	}
	return pointset
}

func FindPoint(user *User, name string) int {
	var temp int
	for key, value := range user.PointMap {
		if key.Name == name {
			temp = value
			break
		}
	}
	return temp
}
