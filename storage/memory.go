package storage

import "ledis/common"

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
	common.Logger.Info("store length = ", length)
	for _, block := range FreeList {
		common.Logger.Info("block = ", block.Offset, block.Length)
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

func GetStore(offset int, length int) []byte {
	return Content[offset:offset+length]
}
