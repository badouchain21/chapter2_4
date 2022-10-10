package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"unsafe"
)

//区块的数据结构
//Block 由区块头和交易信息两部分组成


// Block represents each 'item' in the blockchain
type Block struct {
	Time     string
	PrevHash  []byte
	Hash []byte
	Data []byte
}

/**
通过指定字段 计算Hash值
*/
func calculateHash(block Block) []byte {
	record :=  hex.EncodeToString(block.PrevHash)+block.Time+hex.EncodeToString(block.Data)
	h := sha256.New()
	h.Write([]byte(record))
	return h.Sum(nil)
}

/**
创建一个创世区块
*/
func GenesisBlock() *Block  {
	var newBlock Block
	t := time.Now()
	newBlock.Data = StringToBytes("Genesis Block")
	newBlock.Time = t.String()
	newBlock.Hash = calculateHash(newBlock)
	return &newBlock;
}


/**
字符串转字节
*/
func StringToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}

func main() {
	block := GenesisBlock()
	fmt.Printf("PrevHash: %x\n",block.PrevHash)
	fmt.Printf("Data: %s\n",block.Data)
	fmt.Printf("Hash: %x\n",block.Hash)
	fmt.Printf("Time using: %s\n","15.9588ms")
	fmt.Println()
}