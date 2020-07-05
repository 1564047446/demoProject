//StringToPoint
//字符串转换成坐标点
//package main

/*
将从客户端得到的坐标信息进行解析
将坐标点封装进坐标点类
将坐标点类进行返回
与GetArea相关联
*/

package StringToPoint

import (
	"MyDemo/Room/GetArea"
	"fmt"
	"regexp"
	"strconv"
)

func test() {
	fmt.Println("This is StringToPoint.go")
}

//解析坐标字符串，返回一个只有数字的一个字符串切片
func SplitString(messege string) []string {
	reg := regexp.MustCompile(`[0-9]+`)
	getData := reg.FindAllString(messege, -1)
	return getData
}

//将解析的得到的坐标字符串切片转换成坐标结合
func GetPoint(text string) *GetArea.Point {
	str := SplitString(text)
	p := GetArea.NewPoint()
	tempX := 0
	for i, j := range str {
		if i%4 == 1 {
			x, _ := strconv.Atoi(j)
			tempX = x
		} else if i%4 == 3 {
			y, _ := strconv.Atoi(j)
			p.BuildPoint(float32(tempX), float32(y))
		}
	}
	return p
}

//单独测试模块专用，封装模块的时候记得一定注销，方便调试
/*func main() {
	text := "123,321,33,322,22,1,1,2,3,4,5,6,7,8,9"
	p := GetPoint(text)
	p.PrintPoint()
	GetArea.ParsePointInLine(p)
}*/
