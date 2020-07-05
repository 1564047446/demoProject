// GetArea.go
//package main

/*此模块完成该游戏的核心玩法，并附带将点进行打包方法PackPoint，\
返回字节数组和bool类型，将排序好的方法进行打包，\
用来发送客户端\
包括对点进行有序的调整\
对点所围成的图形进行面积求和\
*/
package GetArea

import (
	"strconv"
	"MyDemo/Proto"
	//"MyDemo/User/GetSkill"
	"fmt"
	"log"
	"sort"
	//"strconv"
)

//坐标点信息和个数
type Point struct {
	x    []float32
	y    []float32
	PNum int
	Name []string
}

//最左侧坐标点
type MinPoint struct {
	x float32
	y float32
	Name string
}

//最右侧坐标点
type MaxPoint struct {
	x float32
	y float32
	Name string
}

//直线坐标左侧的点的集合
type LeftPoint struct {
	x float32
	y float32
	Name string
	//暂时用不上
	//PNum int
}

//直线坐标右侧的点的集合
type RightPoint struct {
	x float32
	y float32
	Name string
	//暂时用不上
	//PNum int
}

//分割线类
type Line struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
	A  float32
	B  float32
	C  float32
}

//最左坐标点构造函数
func NewMinPoint(x float32, y float32, name string) *MinPoint {
	return &MinPoint{
		x: x,
		y: y,
		Name: name,
	}
}

//最右坐标点构造函数
func NewMaxPoint(x float32, y float32, name string) *MaxPoint {
	return &MaxPoint{
		x: x,
		y: y,
		Name: name,
	}
}

//坐标点信息和个数
func NewPoint() *Point {
	return &Point{
		x:    make([]float32, 100),
		y:    make([]float32, 100),
		Name: make([]string, 100),
		PNum: 0,
	}
}

//分割线类构造函数
func NewLine(x1 float32, y1 float32, x2 float32, y2 float32) *Line {
	return &Line{
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
		A:  float32(y2 - y1),
		B:  float32(x1 - x2),
		C:  float32(float32(y1)*float32(x2-x1) - float32(x1)*float32(y2-y1)),
	}
}

//计算环形面积公式
func CalculateArea(p *Point) float32 {
	p.x[p.PNum], p.y[p.PNum] = p.x[0], p.y[0]
	var sum float32 = 0

	//环形面积公式一定要有一定的顺序，顺时针或者逆时针，返回值可能为负数，集的一定要去绝对值。
	for i := 0; i < p.PNum; i++ {
		sum += (p.x[i]*p.y[i+1] - p.x[i+1]*p.y[i])
	}
	sum /= 2

	//判断返回面积的正负
	if sum < 0 {
		return -sum
	}
	return sum
}

//判断点在直线左右
func (l *Line) ParseLine(x float32, y float32) int {
	xx := float32(x)
	yy := float32(y)
	if l.A*xx+l.B*yy+l.C == 0 { //在直线右边
		return 0
	} else if l.A*xx+l.B*yy+l.C < 0 { //在直线左边
		return -1
	} else {
		return 1
	}
}

//构造点坐标
func (p *Point) BuildPoint(x float32, y float32, name string) {
	p.x[p.PNum], p.y[p.PNum] = x, y
	p.Name[p.PNum] = name
	p.PNum++
}

//遍历坐标点
func (p *Point) PrintPoint() {
	for i := 0; i < p.PNum; i++ {
		fmt.Println("x =", p.x[i], "y =", p.y[i])
	}
}

//从坐标集合中找出分割坐标点，并建立直线
func SetParseLine(p *Point) *Line {
	var Maxx float32 = 0
	var Maxy float32 = 0
	var Minx float32 = 0
	var Miny float32 = 0
	var MaxName string
	var MinName string
	//查找方式发生改变
	//找到左下角和右上角的坐标作为分割线从而避免使用sort排序的时候可能会出现乱序的情况
	for i := 0; i < p.PNum; i++ {
		//找出X坐标最小的坐标进行标记
		if p.x[i] < Minx {
			Minx = p.x[i]
			Miny = p.y[i]
			MinName = p.Name[i]
		} else if p.x[i] == Minx { //找到左下角的坐标点，
			//特征：X最小，Y最小
			if p.y[i] < Miny {
				Minx = p.x[i]
				Miny = p.y[i]
				MinName = p.Name[i]
			}
		}
		//找出X坐标最大的坐标进行标记
		if p.x[i] > Maxx {
			Maxx = p.x[i]
			Maxy = p.y[i]
			MaxName = p.Name[i]
		} else if p.x[i] == Maxx {
			if p.y[i] > Maxy {
				Maxx = p.x[i]
				Maxy = p.y[i]
				MaxName = p.Name[i]
			}
		}
	}
	//生成最左侧坐标点
	p1 := NewMinPoint(Minx, Miny, MinName)
	//生成最右侧坐标点
	p2 := NewMaxPoint(Maxx, Maxy, MaxName)
	//建立一条分割线
	l := NewLine(p1.x, p1.y, p2.x, p2.y)
	return l
}

//根据直线将待求坐标点进行划分,根据划分后的点进行计算面积，求得的面积返回
func ParsePointInLine(p *Point) (float32, *Point){
	//var left []LeftPoint
	//var right []RightPoint
	left := make([]LeftPoint, 0)
	right := make([]RightPoint, 0)
	line := SetParseLine(p)
	cntLeft := 0
	cntRight := 0

	//开始进行以直线为分割线进行坐标点的分割
	for i := 0; i < p.PNum; i++ {
		//在直线的左边
		//fmt.Println(p.x[i], p.y[i])
		if line.ParseLine(p.x[i], p.y[i]) == -1 {
			temp := LeftPoint{p.x[i], p.y[i], p.Name[i]}
			left = append(left, temp)
			cntLeft++
		} else {
			//在直线的右方
			temp := RightPoint{p.x[i], p.y[i], p.Name[i]}
			right = append(right, temp)
			cntRight++
		}
	}

	//对分割线左侧坐标进行排序
	sort.Sort(LeftSortPoint(left))
	//对分割线右侧坐标进行排序
	sort.Sort(RightSortPoint(right))

	//调试的时候需要用到，暂时没有什么其他的BUG，暂时不删掉，留一下，防止以后出问题，方便随时查看
	/*
		//将左右坐标区间进行排序
		fmt.Println("left:", left)
		sort.Sort(LeftSortPoint(left))
		fmt.Println("after sort:", left)
		//求出分割后点的总个数
		//fmt.Println(LeftSortPoint(left).Len())
		fmt.Println("right:", right)
		sort.Sort(RightSortPoint(right))
		fmt.Println("after sort:", right)

	*/

	//分割线左侧点的个数，可能为一条线段，所以值可能出现为 0，需要特别注意一下
	leftLen := LeftSortPoint(left).Len()
	//分割线右侧点的个数，同leftLen需要判断图形为一条线段的情况
	rightLen := RightSortPoint(right).Len()

	//构造一个新的点的集合，方便后续计算面积
	pointSet := NewPoint()
	for i := 0; i < leftLen; i++ {
		pointSet.BuildPoint(left[i].x, left[i].y, left[i].Name)
	}
	for i := 0; i < rightLen; i++ {
		pointSet.BuildPoint(right[i].x, right[i].y, right[i].Name)
	}

	//当这个图形为线段时的特殊处理
	if leftLen == 0 {
		pointSet.BuildPoint(right[0].x, right[0].y, right[0].Name)
	} else {
		pointSet.BuildPoint(left[0].x, left[0].y, left[0].Name)
	}
	//最后求得环形面积
	sum := CalculateArea(pointSet)
	fmt.Println("最后求得的面积为:", sum)
	return sum, pointSet
}

//将生成的坐标点打包
//这里的z暂时设定为整形，以后需要改成float32
func PackPoint(point *Point, cmd int) ([]byte, bool) {
	dot := ","
	var messege string
	//信息打包处理
	for i := 0; i < point.PNum; i++ {
		if i == point.PNum - 1 {
			messege = messege + point.Name[i]
			break
		}
		messege = messege + point.Name[i] + dot

	}
	//文本长度 + 8表示userid和cmd所占用的字节数
	TextLen := len(messege) + 8
	//发送指令
	Cmd := cmd
	//用户ID，其实就相当于对两个客户进行一个广播处理
	UserId := 0
	//测试输出
	//fmt.Println("string:", messege)
	//将string强转byte数组
	mes := []byte(messege)
	//暂时未考虑到ok = false的情况
	data, ok := Proto.EnPackage(mes, TextLen, UserId, Cmd)
	if !ok {
		//...............
		log.Println("Enpackage Wrong in SetRand")
		return data, false
	}
	return data, true
}

/*

对左侧需要排序的点进行sort接口的一个封装

*/
type LeftSortPoint []LeftPoint

func (p LeftSortPoint) Len() int {
	return len(p)
}

func (s LeftSortPoint) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p LeftSortPoint) Less(i, j int) bool {
	if p[i].x == p[j].x {
		return p[i].y < p[j].y
	}
	return p[i].x < p[j].x
}

/*
左侧点sort接口封装完毕
*/

/*

对右侧需要排序的点进行sort接口进行封装

*/
type RightSortPoint []RightPoint

func (p RightSortPoint) Len() int {
	return len(p)
}

func (s RightSortPoint) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p RightSortPoint) Less(i, j int) bool {
	if p[i].x == p[j].x {
		return p[i].y > p[j].y
	}
	return p[i].x > p[j].x
}

/*

右侧待排序点封装完毕

*/

func test() {
	fmt.Println("This is GetArea.go")
}

type TempPoint struct {
	x float32
	z float32
}

func NewTemp(x, z float32) *TempPoint {
	return &TempPoint {
		x : x,
		z : z,
	}
}


type Circle struct {
	x float32
	y float32
	z float32
	r float32
	Name string
}

var Circles []Circle = []Circle{{0,1,16, 2, "stone (12)"}, {0,1,27, 2, "stone (1)"}, {-19,1,19, 2, "stone (2)"}, {0,1,37.5, 2, "stone (3)"}, {-11.5,1,11.5, 2, "stone (4)"}, {-11,1,-11, 2, "stone (5)"}, {26.5,1,0, 2, "stone (6)"}, {26.5,1,26.5, 2, "stone (7)"}, {26.5,1,-26.5, 2, "stone (8)"}, {-16,1,0, 2, "stone (9)"}, {19,1,19, 2, "stone (10)"}, {-38,1,0, 2, "stone (11)"}}

func InCircle(x, z float32) bool {
	for i := 0; i < len(Circles); i++ {
		if (Circles[i].x - 2 <= x && Circles[i].x + 2 >= x) && (Circles[i].z - 2 <= z && Circles[i].z + 2 >= z) {
			return false
		}
	}
	return true
}

func SamePoint(x1, z1, x2, z2 float32) bool {
	if (x1 == x2 && z1 == z2) {
		return true
	} else {
		return false
	}
}

//判断线段是否与圆形相交
func CircleLine(point *Point) (int, string) {
	cnt := 0
	temp1 := NewTemp(point.x[0], point.y[0])
	temp2 := NewTemp(point.x[1], point.y[1])
	var name string
	for i := 0; i < len(Circles); i++ {		
		var A, B, C, dist1, dist2, angle1, angle2 float32
			
		if temp1.x == temp2.x {
			A, B, C = 1, 0, -temp1.x
		} else if temp1.z == temp2.z {
			A, B, C = 0, 1, -temp1.z
		} else {
			A = temp1.z - temp2.z
			B = temp2.x  - temp1.x
			C = temp1.x * temp2.z - temp1.z * temp2.x
		}
		
		dist1 = A * Circles[i].x + B * Circles[i].z + C
		dist1 *= dist1
		dist2 = (A * A + B * B) * Circles[i].r * Circles[i].r
		
		if dist1 > dist2 {
			//点到直线的距离大于半径 不相交
			continue
		}
		angle1 = (Circles[i].x - temp1.x) * (temp2.x - temp1.x) + (Circles[i].z - temp1.z) * (temp2.z - temp1.z)
		angle2 = (Circles[i].x - temp2.x) * (temp1.x - temp2.x) + (Circles[i].z - temp2.z) * (temp1.z - temp2.z)
		if angle1 > 0 && angle2 > 0 {
			if cnt == 0 {
				name = name + Circles[i].Name
			} else {
				name = name + "," + Circles[i].Name
			}
			cnt++
		}
		
	}
	return cnt, name
}

//打包技能
func PackSkill(str string) []byte {
	log.Println("获得技能")
	//str := GetSkill.RandSkill()
	cmd := 2008
	id := 1
	mes := []byte(str)
	TextLen := len(mes) + 8
	data, _ := Proto.EnPackage(mes, TextLen, id, cmd)
	return data
}

func PackArea(areaA, areaB int) []byte {
	mes := strconv.Itoa(areaA)
	mes = mes + "," + strconv.Itoa(areaB)
	cmd := 5008
	id := 1
	data := []byte(mes)
	texlen := len(data) + 8
	messege, _ := Proto.EnPackage(data, texlen, id, cmd)
	return messege
}

//分别测试了三角形，正方形，和横竖两条线段，完全正确，面积求得符合预期要求

