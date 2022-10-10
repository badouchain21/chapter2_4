package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
)

/**
区块的数据结构
Block 由区块头和交易信息两部分组成
首先从 “区块” 谈起。在区块链中，真正存储有效信息的是区块（block）。而在比特币中，真正有价值的信息就是交易（transaction）。实际上，交易信息是所有加密货币的价值所在。除此以外，区块还包含了一些技术实现的相关信息，比如版本，当前时间戳和前一个区块的哈希。
不过，我们要实现的是一个简化版的区块链，而不是一个像比特币技术规范所描述那样成熟完备的区块链。所以在我们目前的实现中，区块仅包含了部分关键信息，它的数据结构如下：
 */
//区块的数据结构
type Block struct {
	Timestamp     int64   //当前时间戳
	Data          []byte  //区块实际存储的信息
	PrevBlockHash []byte  //前一个块的哈希
	Hash          []byte  //当前块的哈希
	Nonce         int   //在对工作量证明进行验证时用到
}


/**
在我们的简化版区块中，还有一个 Hash 字段，那么，要如何计算哈希呢？哈希计算，是区块链一个非常重要的部分。正是由于它，才保证了区块链的安全。计算一个哈希，是在计算上非常困难的一个操作。即使在高速电脑上，也要耗费很多时间 (这就是为什么人们会购买 GPU，FPGA，ASIC 来挖比特币) 。这是一个架构上有意为之的设计，它故意使得加入新的区块十分困难，继而保证区块一旦被加入以后，就很难再进行修改。在接下来的内容中，我们将会讨论和实现这个机制。
目前，我们仅取了 Block 结构的部分字段（Timestamp, Data 和 PrevBlockHash），并将它们相互拼接起来，然后在拼接后的结果上计算一个 SHA-256，然后就得到了哈希。把这个功能用以下的SetHash函数来实现。
 */
//设置当前块哈希
// Hash=sha256(PrevBlockHash+Data+Timestamp)
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

//用于生成新块
//当前块的哈希会基于传入的参数Data 和PrevBlockHash计算得到
//func NewBlock(data string, prevBlockHash []byte) *Block {
//	block := &Block{
//		Timestamp:		time.Now().Unix(),
//		PrevBlockHash:	prevBlockHash,
//		Hash:			[]byte{},
//		Data:			[]byte(data) }
//
//	block.SetHash()
//	return block
//}

/**
有了区块，下面让我们来实现区块链。本质上，区块链就是一个有着特定结构的数据库，是一个有序，每一个块都连接到前一个块的链表。也就是说，区块按照插入的顺序进行存储，每个块都与前一个块相连。这样的结构，能够让我们快速地获取链上的最新块，并且高效地通过哈希来检索一个块。
在 Golang 中，可以通过一个 array 和 map 来实现这个结构：array 存储有序的哈希（Golang 中 array 是有序的），map 存储 hash -> block 对(Golang 中, map 是无序的)。 但是在基本的原型阶段，我们只用到了 array，因为现在还不需要通过哈希来获取块。
 */
//第一个区块链————这是一个Block 指针数组
type Blockchain struct {
	blocks []*Block
}

//现在，让我们能够给它添加一个区块：
//添加区块
//data就是交易
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

/**
为了加入一个新的块，我们必须要有一个已有的块，但是，初始状态下，我们的链是空的，一个块都没有！所以，在任何一个区块链中，都必须至少有一个块。这个块，也就是链中的第一个块，通常叫做创世块（genesis block）
 */
//创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block1", []byte{})
}

//创建一个有创世块的区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

//=========================================工作量证明（Proof-of-Work）================================
/**
我们构造了一个非常简单的数据结构 – 区块，它也是整个区块链数据库的核心。目前所完成的区块链原型，已经可以通过链式关系把区块相互关联起来：每个块都与前一个块相关联。
但是，当前实现的区块链有一个巨大的缺陷：向链中加入区块太容易，也太廉价了。而区块链和比特币的其中一个核心就是，要想加入新的区块，必须先完成一些非常困难的工作。在本文，我们将会弥补这个缺陷。
区块链的一个关键点就是，一个人必须经过一系列困难的工作，才能将数据放入到区块链中。正是由于这种困难的工作，才保证了区块链的安全和一致。此外，完成这个工作的人，也会获得相应奖励（这也就是通过挖矿获得币）。
这个机制与生活现象非常类似：一个人必须通过努力工作，才能够获得回报或者奖励，用以支撑他们的生活。在区块链中，是通过网络中的参与者（矿工）不断的工作来支撑起了整个网络。矿工不断地向区块链中加入新块，然后获得相应的奖励。在这种机制的作用下，新生成的区块能够被安全地加入到区块链中，它维护了整个区块链数据库的稳定性。值得注意的是，完成了这个工作的人必须要证明这一点，即他必须要证明他的确完成了这些工作。
整个 “努力工作并进行证明” 的机制，就叫做工作量证明（proof-of-work）。要想完成工作非常地不容易，因为这需要大量的计算能力：即便是高性能计算机，也无法在短时间内快速完成。另外，这个工作的困难度会随着时间不断增长，以保持每 10 分钟出 1 个新块的速度。在比特币中，这个工作就是找到一个块的哈希，同时这个哈希满足了一些必要条件。这个哈希，也就充当了证明的角色。因此，寻求证明（寻找有效哈希），就是矿工实际要做的事情。
 */

/**
哈希计算
获得指定数据的一个哈希值的过程，就叫做哈希计算。一个哈希，就是对所计算数据的一个唯一表示。对于一个哈希函数，输入任意大小的数据，它会输出一个固定大小的哈希值。下面是哈希的几个关键特性：
1.无法从一个哈希值恢复原始数据。也就是说，哈希并不是加密。
2.对于特定的数据，只能有一个哈希，并且这个哈希是唯一的。
3.即使是仅仅改变输入数据中的一个字节，也会导致输出一个完全不同的哈希
在区块链中，哈希被用于保证一个块的一致性。哈希算法的输入数据包含了前一个块的哈希，因此使得不太可能（或者，至少很困难）去修改链中的一个块：因为如果一个人想要修改前面一个块的哈希，那么他必须要重新计算这个块以及后面所有块的哈希。

 */

/**
Hashcash
比特币使用 Hashcash ，一个最初用来防止垃圾邮件的工作量证明算法。它可以被分解为以下步骤：

1.取一些公开的数据（比如，如果是 email 的话，它可以是接收者的邮件地址；在比特币中，它是区块头）
2.给这个公开数据添加一个计数器。计数器默认从 0 开始
3.将 data(数据) 和 counter(计数器) 组合到一起，获得一个哈希
4.检查哈希是否符合一定的条件：
1.如果符合条件，结束
2.如果不符合，增加计数器，重复步骤 3-4
ca07ca 是计数器的 16 进制值，十进制的话是 13240266.
*/

//定义挖矿的难度值 ，以下表示哈希的前24位必须是0
const targetBits=24

/**
在比特币中，当一个块被挖出来以后，“target bits” 代表了区块头里存储的难度，也就是开头有多少个 0。这里的 24 指的是算出来的哈希前 24 位必须是 0，如果用 16 进制表示，就是前 6 位必须是 0，这一点从最后的输出可以看出来。目前我们并不会实现一个动态调整目标的算法，所以将难度定义为一个全局的常量即可。
24 其实是一个可以任意取的数字，其目的只是为了有一个目标（target）而已，这个目标占据不到 256 位的内存空间。同时，我们想要有足够的差异性，但是又不至于大的过分，因为差异性越大，就越难找到一个合适的哈希
 */
//每个块的工作量都必须要证明，所以有个指向Block的指针
//target是目标，我们最终要找的哈希必须要小于目标
type ProofOfWork struct {
	block  *Block
	target *big.Int
}


//target等于1左移256-targetBits 位？
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}
//工作量证明需要用到的数据有：PrevBlockHash, Data, Timestamp, targetBits, nonce(计数器，密码学术语)
func (pow *ProofOfWork) prepareData(nonce int) []byte {   //这个方法用来准备数据，也可以用来验证工作量
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}
//将一个 int64 转化为一个字节数组（byte array）
func IntToHex(num int64) []byte {
	buff:=new(bytes.Buffer)
	err:=binary.Write(buff, binary.BigEndian, num)
	if err !=nil{
		log.Panic(err)
	}

	return buff.Bytes()

}

var (maxNonce = math.MaxInt64)  //对循环进行限制


//Pow算法的核心就是寻找有效哈希
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int //hashInt是hash的整形表示
	var hash [32]byte
	nonce := 0  //计数器

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {   //防止溢出的“无限”循环
		data := pow.prepareData(nonce)   //准备数据
		hash = sha256.Sum256(data)        //对数据进行哈希计算
		hashInt.SetBytes(hash[:])         //将将哈希转换成一个大整数

		if hashInt.Cmp(pow.target) == -1 {   //将大整数与目标进行比较
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()  //调用计算哈希的方法

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
//验证工作量，只要哈希小于目标就是有效工作量
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}



//测试
func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))  //FormatBool 将布尔值转换为字符串 "true" 或 "false"
		fmt.Println()
	}
}
/**
可以看到每个哈希都是 3 个字节的 0 开始，并且获得这些哈希需要花费一些时间，这次我们产生三个块花费了一分多钟，比没有工作量证明之前慢了很多（也就是成本高了很多）。

我们离真正的区块链又进了一步：现在需要经过一些困难的工作才能加入新的块，因此挖矿就有可能了。但是，它仍然缺少一些至关重要的特性：区块链数据库并不是持久化的，没有钱包，地址，交易，也没有共识机制。不过，所有的这些，我们都会在接下来的文章中实现，现在，愉快地挖矿吧！
 */


