package main

import (
	"golang.org/x/tour/tree"
	"testing"
)

func TestWalk(t *testing.T) {
	t1 := tree.New(1)
	ch := make(chan int)

	got := make([]int, 0)
	except := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	go Walk(t1, ch)
	for i := range ch {
		got = append(got, i)
	}

	if len(got) != len(except) {
		t.Errorf("[x] no correct, lens are different")
	}
	for i := 0; i < len(got); i++ {
		if got[i] != except[i] {
			t.Errorf("[x] no correct, values are different")
		}
	}
}

type sameTest struct {
	t1, t2   *tree.Tree
	excepted bool
}

var sameTests = []sameTest{
	sameTest{tree.New(1), tree.New(1), true},
	sameTest{tree.New(1), tree.New(2), false},
	sameTest{tree.New(3), tree.New(3), true},
}

func TestSame(t *testing.T) {
	for i, test := range sameTests {
		got := Same(test.t1, test.t2)
		if got != test.excepted {
			t.Errorf("[x] test#%d => excepted:%v, got:%v", i, test.excepted, got)
		}
	}
}

// type addTest struct {
// 	toAdd, excepted []int
// }

// var addTests = []addTest{
// 	addTest{[]int{1,2,1,4,1,6,1,8,1,10}, []int{1,1,1,1,1,2,4,6,8,10}},
// }


// func TestAdd(t *testing.T) {

// }