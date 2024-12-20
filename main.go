package main

import (
	"fmt"
)

// binary search tree
type Node struct {
	data  int
	left  *Node
	right *Node
}

func (n *Node) insert(data int) {
	if n == nil {
		return
	}

	if data < n.data {
		if n.left == nil {
			n.left = &Node{data: data}

		} else {
			n.left.insert(data)
		}

	} else {
		if n.right == nil {
			n.right = &Node{data: data}
		} else {
			n.right.insert(data)
		}
	}
}

func (n *Node) search(data int) (bool, int) {
	if n == nil {
		return false, 0
	}
	if data == n.data {
		return true, n.data
	}
	if data < n.data {
		return n.left.search(data)
	}
	return n.right.search(data)
}

func (n *Node) PrintInOrder() {
	if n == nil {
		return
	}

	n.left.PrintInOrder()
	fmt.Print(n.data, " ")
	n.right.PrintInOrder()
}

func (n *Node) PrintPreOrder() {
	if n == nil {
		return
	}
	fmt.Print(n.data, " ")
	n.left.PrintPreOrder()
	n.right.PrintPreOrder()
}

//! Step-1 make the tree structure
//! step-2 make sha1 hashing algorithm
//! step-3 implement the tree into the git algorithm thingy
//! step-4 networking
//! step-5 add the git algorithm to the networking
//! step-6 make it so that i can form adapters and hooks on every commit
//! profit ?

//~ never done go btw : p and i got an interview in it

func main() {

	result, err := Encrypt("./example-repo/main.c")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)

}
