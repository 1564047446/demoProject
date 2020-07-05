// SetRand

/*
设置初始随机点模块
将设置的随机点，通过Proto进行打包处理，发送给服务器端
与proto相互关联
*/

package SetRand

//package main

import (
	"MyDemo/Proto"
	"MyDemo/Room/GetArea"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func test() {
	fmt.Println("This is SetRand.go")
}

//边界判断100 * 100 中心在坐标原点，范围平分1，2，3，4象限
//要求：尽量在中心附近生成左边
//坐标生成10个
const (
	MaxLen = 75
	MaxX   = 35
	MaxY   = 50
	MinX   = -50
	MinY   = -35
	//点生成个数
	PointNum = 20
)

//生成点的类
type Point struct {
	X []int
	Y float64
	Z []int
}

//初始化坐标点信息
func NewPoint(x int) *Point {
	return &Point{
		X: make([]int, x),
		//y轴坐标未一个恒定值
		Y: 0.8,
		Z: make([]int, x),
	}
}

/**********************************
*	以下模块未彻底完成，缺少边界判断	  *
***********************************/

var Num = []int{40, 30, 20}
var RangeLen = []int{70, 50, 30}
var RangeXY = []int{35, 25, 15}

//随机生成15个坐标点，y轴固定，只生成x，z就可以，以一个原点为中心进行生成
func (p *Point) GetPoint(OriX int, OriY int, x int) {
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < Num[x]; i++ {
		p.X[i] = rand.Intn(RangeLen[x]) - RangeXY[x]
		p.Z[i] = rand.Intn(RangeLen[x]) - RangeXY[x]
		//与石柱不冲突
		if !GetArea.InCircle(float32(p.X[i]), float32(p.Z[i])) {
			i--;
			continue
		}
		//第一次与圆心不冲突
		if (p.X[i] <= 5 && p.X[i] >= -5) && (p.Z[i] <= 5 && p.Z[i] >= -5) && x == 0 {
			i--;
			continue
		}
		//需要判断边界处理条件
	}
}

//生成随机坐标点
func InitPoint(x int) *Point {
	point := NewPoint(Num[x])
	point.GetPoint(0, 0, x)
	return point
}



//将生成的坐标点打包
func PackPoint(point *Point, x int) ([]byte, bool) {
	//未判断点做了改动
	dot := ","
	var messege_len = 0
	var messege string
	//将坐标点命名
	//name := "Marker_"
	//信息打包处理
	for i := 0; i < Num[x]; i++ {
		if i == Num[x] - 1 {
			messege = messege + strconv.Itoa(point.X[i])
			messege = (messege + dot)
			messege = (messege + strconv.FormatFloat(point.Y, 'E', -1, 32))
			messege = (messege + dot)
			messege = (messege + strconv.Itoa(point.Z[i]))
			messege_len += 5
			break
		}
		messege = (messege + strconv.Itoa(point.X[i]))
		messege = (messege + dot)
		messege = (messege + strconv.FormatFloat(point.Y, 'E', -1, 32))
		messege = (messege + dot)
		messege = (messege + strconv.Itoa(point.Z[i]))
		messege = (messege + dot)
		messege_len += 6
	}
	//文本长度 + 8表示userid和cmd所占用的字节数
	TextLen := len(messege) + 8
	//发送指令
	Cmd := 1099
	//用户ID，其实就相当于对两个客户进行一个广播处理\

	UserId := 0
	//测试输出
	fmt.Println("string:", messege)
	//将string强转byte数组
	mes := []byte(messege)
	//暂时未考虑到ok = false的情况
	data, ok := Proto.EnPackage(mes, TextLen, UserId, Cmd)
	if !ok {
		log.Println("Enpackage Wrong in SetRand")
		return data, false
	}
	return data, true
}

//测试模块专用，不需要时可以直接注释掉

