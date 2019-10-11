package htable

import "ledis/utils"

// 内存存储

// Block 存储块
type Block struct {
	Offset int // 存储空间 Content 中的 offset
	Length int // 存储空间 Content 中的 length
	Next   *Block
}

var FreeLength int    // 剩余空间大小
var Content []byte    // 存储
var FreeList []*Block // 空闲存储块
var UsedList []*Block // 已使用的存储块

func init() {
	FreeLength = 4000000
	Content = make([]byte, FreeLength)
	initBlock := &Block{
		Offset: 0,
		Length: FreeLength,
	}
	FreeList = append(FreeList, initBlock)
}

func Store(content []byte) *Block {
	length := len(content)
	for _, block := range FreeList {
		if block.Length >= length {
			for i:=block.Offset;i<block.Offset+length;i++ {
				Content[i] = content[i-block.Offset]
			}
			newBlock := &Block{
				Offset: block.Offset,
				Length: length,
			}
			UsedList = append(UsedList, newBlock)
			block.Offset = block.Offset + length
			block.Length -= length
			return newBlock
		}
	}
	return nil
}

func (block *Block) GetContent() []byte {
	return Content[block.Offset:block.Offset+block.Length]
}

// KeyToInt 将 Key 转换为数组下标
func KeyToInt(key interface{}) int {
	return utils.HashInt(key) % HashTableLength
}