package imagesha1

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// modifyPNGSHA1 修改PNG图片的SHA1值
// 通过在PNG文件中添加自定义文本块来改变SHA1，不影响图片显示
func (m *ImageSHA1Modifier) modifyPNGSHA1(data []byte) ([]byte, error) {
	// PNG文件必须以PNG签名开头
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return nil, fmt.Errorf("不是有效的PNG文件")
	}

	// 生成随机文本作为自定义块
	randomText := m.generateRandomBytes(32)

	// 在PNG中插入自定义文本块
	modifiedData := m.insertPNGTextChunk(data, "Random", string(randomText))

	return modifiedData, nil
}

// insertPNGTextChunk 在PNG文件中插入文本块
func (m *ImageSHA1Modifier) insertPNGTextChunk(data []byte, keyword, text string) []byte {
	// PNG块结构：
	// [4字节长度][4字节类型][数据][4字节CRC]

	// 查找IEND块的位置（PNG文件的最后一个块）
	iendPos := m.findPNGIENDChunk(data)
	if iendPos == -1 {
		// 如果找不到IEND块，直接在文件末尾添加
		iendPos = len(data)
	}

	// 构造tEXt块数据
	textData := []byte(keyword)
	textData = append(textData, 0) // 分隔符
	textData = append(textData, []byte(text)...)

	// 构造完整的tEXt块
	chunk := make([]byte, 0, 12+len(textData))

	// 长度（4字节，大端序）
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(textData)))
	chunk = append(chunk, lengthBytes...)

	// 类型（4字节）
	chunkType := []byte("tEXt")
	chunk = append(chunk, chunkType...)

	// 数据
	chunk = append(chunk, textData...)

	// CRC（4字节）
	crcData := append(chunkType, textData...)
	crc := crc32.ChecksumIEEE(crcData)
	crcBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(crcBytes, crc)
	chunk = append(chunk, crcBytes...)

	// 构造新的PNG数据
	result := make([]byte, 0, len(data)+len(chunk))
	result = append(result, data[:iendPos]...) // IEND之前的数据
	result = append(result, chunk...)          // 新的文本块
	result = append(result, data[iendPos:]...) // IEND块

	return result
}

// findPNGIENDChunk 查找PNG文件中IEND块的位置
func (m *ImageSHA1Modifier) findPNGIENDChunk(data []byte) int {
	// PNG签名是8字节
	pos := 8

	for pos < len(data)-8 {
		// 读取块长度
		if pos+4 > len(data) {
			break
		}
		chunkLen := binary.BigEndian.Uint32(data[pos : pos+4])
		pos += 4

		// 读取块类型
		if pos+4 > len(data) {
			break
		}
		chunkType := data[pos : pos+4]
		pos += 4

		// 检查是否是IEND块
		if bytes.Equal(chunkType, []byte("IEND")) {
			return pos - 8 // 返回块开始位置
		}

		// 跳过数据和CRC
		pos += int(chunkLen) + 4
	}

	return -1 // 未找到IEND块
}
