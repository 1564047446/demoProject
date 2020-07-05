//PackStatus.go

//玩家状态整合和打包
//package packStus
package main

import (
	"MyDemo/User"
	//"log"
)

type PointStatus struct {
	Name   string
	Status int
}

func UpdateStatus(userA *User.User, userB *User.User) []PointStatus {
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

func FindPoint(user *User.User, name string) int {
	var temp int
	for key, value := range user.PointMap {
		if key.Name == name {
			temp = value
			break
		}
	}
	return temp
}
