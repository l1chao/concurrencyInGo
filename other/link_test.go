package main_test

import "testing"

type LinkNode struct {
	Val  int
	Next *LinkNode
}

func DeleteX(l *LinkNode, x int) {
	p := l // l是头节点
	for p.Next != nil {
		if p.Next.Val == x {
			p.Next = p.Next.Next
		} else {
			p = p.Next
		}
	}
}

func Test1(t *testing.T) {
	//有头节点的删除x

}
