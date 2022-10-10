package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"unsafe"
)
/**
使用Go语言实现一条区块链
 */

// Block represents each 'item' in the blockchain
type Block struct {
	Time     string
	PrevHash  []byte
	Hash []byte
	Data []byte
}
/**
区块链的基本构成单位是区块，区块又分为区块头和区块体两部分.我们实现的简单区块链则是按照上述分类,用以下字段构成的。
字身 Tme出放时间聊，也就是区块创建的时间，数据类型为int64。
字段 PrevHash首一个地的 Hash 值，即父哈希，数据类型为[]byte。字段Hash 出首地的Hash值,教据类型为[byte。
字段Data,区地右出出守际右放信息也就是交易，数据类型为[byte。
其中,Timestamp. PrevHash、Hash属于区块头，Data则周个nn-o天 定区中的“交易（Transaction)”字段，目前暂时不会涉及太复杂的结遒，I uaitaI一串字符信息即可.字Hach表示出前反块的Hash值，是由工作里证明开么”开付到的，是区块链安全性的基石,将在后续实验环节中进行介绍。
 */
// Blockchain is a series of validated Blocks
type Blockchain struct {
	blocks []*Block
}

/**
创建新的区块
 */
func (bc *Blockchain) NewBlock(data string)  {
	var newBlock Block
	t := time.Now()
	preBlock := bc.blocks[len(bc.blocks)-1];
	newBlock.Time = t.String()
	newBlock.PrevHash = preBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	newBlock.Data = StringToBytes(data)
	bc.blocks = append(bc.blocks,&newBlock)
}

/**
字符串转字节
 */
func StringToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
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
返回一条区块链 携带创世区块
 */
func NewBlockchain() *Blockchain  {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func main() {
	//初始化区块链
	bc := NewBlockchain()
	//添加区块
	bc.NewBlock("Send 1 BTC to Ivan")
	bc.NewBlock("Send 2 more BTC to Ivan")
	//循环打印
	for _,block := range bc.blocks {
		fmt.Printf("PrevHash: %x\n",block.PrevHash)
		fmt.Printf("Data: %s\n",block.Data)
		fmt.Printf("Hash: %x\n",block.Hash)
		fmt.Println()
	}
}
