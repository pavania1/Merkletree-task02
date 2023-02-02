package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type MerkleTree struct {
	Root *Merkletreenode
}

type Merkletreenode struct {
	Left  *Merkletreenode
	Right *Merkletreenode
	Data  []byte
}

func buildMerkletreenode(left *Merkletreenode, right *Merkletreenode, Data []byte) *Merkletreenode {

	newnode := Merkletreenode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(Data)
		newnode.Data = hash[0:]
	} else {
		prevHashVal := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashVal)
		newnode.Data = hash[0:]

	}

	newnode.Left = left
	newnode.Right = right

	return &newnode

}

func buildMerkleTree(values [][]byte) *MerkleTree {
	var newnodes []Merkletreenode

	for _, v := range values {
		newnode := buildMerkletreenode(nil, nil, v)
		newnodes = append(newnodes, *newnode)
	}

	for len(newnodes) > 1 {
		//check whether the length of the tree is even or not
		if len(newnodes)%2 == 1 {
			newnodes = append(newnodes, newnodes[len(newnodes)-1])
		}

		//here we are changing every val in value into merkle node
		var parentHashes []Merkletreenode
		for i := 0; i < len(newnodes); i += 2 {
			newnode := buildMerkletreenode(&newnodes[i], &newnodes[i+1], nil)
			parentHashes = append(parentHashes, *newnode)
		}
		newnodes = parentHashes
	}

	tree := MerkleTree{&newnodes[0]}

	return &tree
}

func Addnode(values []string) []string {
	var Newnode string
	fmt.Println("Enter the node you want to add to the merkle tree")
	fmt.Scanln(&Newnode)
	values = append(values, Newnode)
	return values

}

func DeleteNode(Values []string) []string {
	var deletednode string
	flag := 0
	fmt.Println("enter the node you want to delete")
	fmt.Scanln(&deletednode)
	for i := 0; i < len(Values); i++ {
		if Values[i] == deletednode {
			Values = append(Values[:i], Values[i+1:]...)
			flag = 1
		}
	}
	if flag == 0 {
		fmt.Println("transaction not found")
	}
	return Values
}

func (root *Merkletreenode) verify(data string) bool {
	var hash []byte
	bytedata := []byte(data)
	hash32 := sha256.Sum256(bytedata)
	hash = append(hash, hash32[0:]...)
	return VerifyNode(root, hash)
}

func VerifyNode(root *Merkletreenode, target []byte) bool {
	if root == nil {
		return false
	}

	if bytes.Equal(root.Data, target) {
		return true
	}
	var left, right bool
	if root.Left != nil {
		left = VerifyNode(root.Left, target)
	}
	if root.Right != nil {
		right = VerifyNode(root.Right, target)
	}
	return left || right
}

func main() {
	values := []string{"Pavani", "is", "doing", "internship", "in", "vitwit", "company"}
	//values := []string{}
	updatedvalues := Addnode(values)
	var rootHash *MerkleTree
	DeleteNode(updatedvalues)
	data := [][]byte{}

	if len(updatedvalues) == 0 {
		fmt.Println("NO transactions")
	} else {
		for k, v := range updatedvalues {
			_ = k
			a := []byte(v)
			data = append(data, a)
		}
		rootHash = buildMerkleTree(data)
		// fmt.Println("Root hash", hex.EncodeToString(rootHash.Root.Value))
	}
	// fmt.Println("roothash", rootHash.Root)
	c := rootHash.Root
	fmt.Println("node is present/not: ", c.verify("vitwit"))

}
