package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func SubWalk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	ch <- t.Value
	SubWalk(t.Left, ch)
	SubWalk(t.Right, ch)
}

func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		close(ch)
		return
	}
	SubWalk(t, ch)
	close(ch)
}

func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int, 10), make(chan int, 10)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		x, okx := <-ch1
		y, oky := <-ch2           // 0
		if okx != oky || x != y { // ! ( okx ^ oky )
			return false
		} else if !okx {
			break
		}
	}
	return true
}

func AddNode(t1 *tree.Tree, v int) error {
	if t1 == nil {
		return fmt.Errorf("nil pointer exception")
	}
	if t1.Value >= v {
		if t1.Left == nil {
			t1.Left = &tree.Tree{Value: v, Left: nil, Right: nil}
			return nil
		}
		return AddNode(t1.Left, v)
	} else if t1.Value < v {
		if t1.Right == nil {
			t1.Right = &tree.Tree{Value: v, Left: nil, Right: nil}
			return nil
		}
		return AddNode(t1.Right, v)
	}
	return fmt.Errorf("some error exception")
}

func Print(t1 *tree.Tree) {
	if t1 == nil {
		return
	}
	go Print(t1.Left)
	go print(t1.Right)
}

func main() { // YOU CAN TEST BY USING go test -v !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	t1, t2 := tree.New(1), tree.New(1)
	ch := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch)
	for i := range ch {
		fmt.Print(i, " ")
	}
	fmt.Println()
	go Walk(t2, ch2)
	for i := range ch2 {
		fmt.Print(i, " ")
	}
	// AddNode(t1, 3)
	// go Walk(t1, ch)
	// for i := range ch {
	// 	fmt.Println(i)
	// }
	// fmt.Println(Same(t1, t2))
}
