//获取技能
//package main

package GetSkill

import (
	"math/rand"	
)

var Skills = [100]string{"神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","神行","迷魂","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","瞬移","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","巨像","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","迷魂","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒","神怒"}
//var Skills = [100]string{"捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁","捆仙锁"}
func RandSkill() string {
	x := rand.Intn(100)
	return Skills[x]
}

