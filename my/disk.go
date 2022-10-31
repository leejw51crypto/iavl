package main

import (
	"fmt"
	"os"

	"github.com/cosmos/iavl"
	//dbm "github.com/tendermint/tm-db"
	dbm "github.com/cosmos/cosmos-db"
)

var USE_FORCE_STORAGE = false  // false: switch fastnode, legacy, true: force below
var USE_LEGACY_STORAGE = false //true: not use fast ndoe, false: use fast node
var TX_COUNT int64 = 1

type Disk struct {
	mydb             *dbm.GoLevelDB
	mytree           *iavl.MutableTree
	useLegacyStorage bool
}

func (disk *Disk) writeBlock(blockid int64) {
	fmt.Printf("--------------------------\nwrite block %d\n", blockid)
	var key = []byte("alice")
	data := fmt.Sprintf("abc%d", blockid)
	fmt.Printf("begin Set key %s  data %s\n", string(key), data)
	disk.mytree.Set(key, []byte(data))
	fmt.Println("end Set")
}

func (disk *Disk) writeInt64(key string, v int64) error {
	error := disk.mydb.Set([]byte(key), []byte(fmt.Sprintf("%d", v)))
	return error
}
func (disk *Disk) readInt64(key string) (int64, error) {
	v, e := disk.mydb.Get([]byte(key))
	if e != nil {
		return 0, e
	}
	var ret int64 = 0
	// convert v to string
	s := string(v)
	// scanf string to int64
	_, e = fmt.Sscanf(s, "%d", &ret)
	if e != nil {
		return 0, e
	}

	return ret, nil
}
func (disk *Disk) initialize() {

	dbpath, _ := os.Getwd()
	dbpath = dbpath + "/data"
	disk.mydb, _ = dbm.NewGoLevelDBWithOpts("test", dbpath, nil)

	uselegacy, err := disk.readInt64("uselegacystorage")
	if err != nil {
		uselegacy = 1
	}
	// print uselegacy, err
	fmt.Printf("uselegacy %d err %v\n", uselegacy, err)
	// switch
	if uselegacy == 1 {
		disk.useLegacyStorage = true
		disk.writeInt64("uselegacystorage", 0)
	} else {
		disk.useLegacyStorage = false
		disk.writeInt64("uselegacystorage", 1)
	}
	fmt.Printf("useLegacyStorage %v\n", disk.useLegacyStorage)

	// force disk.useLegacyStorage
	if USE_FORCE_STORAGE {
		fmt.Printf("force useLegacyStorage useLegacyStorage %v\n", USE_LEGACY_STORAGE)
		disk.useLegacyStorage = USE_LEGACY_STORAGE
	}

	disk.mytree, _ = iavl.NewMutableTree(disk.mydb, 128, disk.useLegacyStorage)
	latestversion, _ := disk.mytree.Load()
	// get latest hash
	root, _ := disk.mytree.Hash()
	fmt.Printf("latest version = %v  hash = %x\n", latestversion, root)

	// from db
	myversion, err := disk.readInt64("myversion")
	myhash, err := disk.mydb.Get([]byte("myhash"))
	// print myversion, myhash
	fmt.Printf("myversion %d myhash %x\n", myversion, myhash)

	fmt.Printf("initialized\n")
	fmt.Printf("-----------------------------------------\n")
	if latestversion > 0 {
		// if myversion != latestversion, panic
		if myversion != latestversion {
			panic("myversion != latestversion")
		}
		// if myhash != root, panic
		// print myhash, root
		fmt.Printf("myhash %x root %x\n", myhash, root)
		// print myversion , latestversion
		fmt.Printf("myversion %d latestversion %d\n", myversion, latestversion)
		if string(myhash) != string(root) {
			panic("myhash != root")
		}
	}
}
func (disk *Disk) write() {
	// for loop with 10 times
	var s int64 = 0
	var count int64 = TX_COUNT
	s, _ = disk.readInt64("height")
	fmt.Printf("height %d\n", s)
	for i := s; i < s+count; i++ {
		//time.Sleep(time.Second)
		disk.writeBlock(int64(i))
		fmt.Println("save version..")
		hash, version, err := disk.mytree.SaveVersion()
		fmt.Printf("V:%d H:%x E:%v\n", version, hash, err)
		err = disk.writeInt64("myversion", version)
		if err != nil {
			panic(err)
		}
		err = disk.mydb.Set([]byte("myhash"), hash)
		if err != nil {
			panic(err)
		}
		err = disk.writeInt64("height", int64(i+1))
		if err != nil {
			panic(err)
		}

	}
}
func (disk *Disk) read() {

	var show_all = true
	if show_all {
		versions := disk.mytree.AvailableVersions()
		// print versions
		fmt.Printf("versions %v\n", versions)
		for i, version := range versions {
			value, err := disk.mytree.GetVersioned([]byte("alice"), int64(version))
			fmt.Printf("version:%d value %s   e:%v\n", i, string(value), err)
		}
	}

}
func (disk *Disk) process() {
	fmt.Println("disk process")
	disk.initialize()
	disk.write()
	//disk.read()
	disk.mydb.Close()
}
func RunProcess() {
	disk := Disk{}
	disk.process()
}
