package main

import (
	"fmt"

	"github.com/cosmos/iavl"
)

func show_root_hash(tree *iavl.MutableTree) {
	root, _ := tree.Hash()
	fmt.Printf("version %d roothash %x\n", tree.Version(), root)
}

func save(tree *iavl.MutableTree) {
	hash, version, error := tree.SaveVersion()
	// print hash, version, err
	//fmt.Printf("hash %x version %d \n", hash, version)
	// print hash, version, err
	fmt.Printf("hash %x version %d error %v\n", hash, version, error)
}

func main() {
	fmt.Println("###########################################")
	fmt.Println("###########################################")
	fmt.Println("###########################################")
	RunProcess()
}
