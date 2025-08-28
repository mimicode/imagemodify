package imagemodify

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"strings"
	"time"
)

// modifyPNGMetadata 修改PNG图片的文本元数据
func (m *ImageModifier) modifyPNGMetadata(data []byte, metadata *ImageMetadata) ([]byte, error) {
	// PNG文件必须以PNG签名开头
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return nil, fmt.Errorf("不是有效的PNG文件")
	}

	// 移除现有的文本块
	cleanData := m.removeExistingTextChunks(data)

	// 准备要添加的文本块
	textChunks := m.createPNGTextChunks(metadata)

	// 在IEND块之前插入新的文本块
	return m.insertPNGTextChunks(cleanData, textChunks), nil
}

// createPNGTextChunks 创建PNG文本块
func (m *ImageModifier) createPNGTextChunks(metadata *ImageMetadata) [][]byte {
	var chunks [][]byte

	// 添加各种元数据作为tEXt块
	if metadata.Artist != "" {
		chunks = append(chunks, m.createPNGTextChunk("Author", metadata.Artist))
	}
	if metadata.Copyright != "" {
		chunks = append(chunks, m.createPNGTextChunk("Copyright", metadata.Copyright))
	}
	if metadata.Description != "" {
		chunks = append(chunks, m.createPNGTextChunk("Description", metadata.Description))
	}
	if metadata.DateTime != nil {
		// PNG标准的创建时间格式
		chunks = append(chunks, m.createPNGTextChunk("Creation Time", metadata.DateTime.Format("2006-01-02T15:04:05Z")))
	}
	if metadata.Location != "" {
		chunks = append(chunks, m.createPNGTextChunk("Location", metadata.Location))
	}
	if metadata.CameraMake != "" {
		chunks = append(chunks, m.createPNGTextChunk("Camera Make", metadata.CameraMake))
	}
	if metadata.CameraModel != "" {
		chunks = append(chunks, m.createPNGTextChunk("Camera Model", metadata.CameraModel))
	}
	if metadata.Software != "" {
		chunks = append(chunks, m.createPNGTextChunk("Software", metadata.Software))
	}
	if metadata.ImageWidth > 0 {
		chunks = append(chunks, m.createPNGTextChunk("Image Width", strconv.Itoa(metadata.ImageWidth)))
	}
	if metadata.ImageHeight > 0 {
		chunks = append(chunks, m.createPNGTextChunk("Image Height", strconv.Itoa(metadata.ImageHeight)))
	}

	return chunks
}

// createPNGTextChunk 创建单个PNG文本块
func (m *ImageModifier) createPNGTextChunk(keyword, text string) []byte {
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

	return chunk
}

// removeExistingTextChunks 移除现有的文本块
func (m *ImageModifier) removeExistingTextChunks(data []byte) []byte {
	result := make([]byte, 0, len(data))
	result = append(result, data[:8]...) // 保留PNG签名

	pos := 8 // 跳过PNG签名

	for pos < len(data)-8 {
		// 读取块长度
		if pos+4 > len(data) {
			break
		}
		chunkLen := binary.BigEndian.Uint32(data[pos : pos+4])

		// 读取块类型
		if pos+8 > len(data) {
			break
		}
		chunkType := data[pos+4 : pos+8]

		// 计算整个块的大小（长度+类型+数据+CRC）
		totalChunkSize := int(chunkLen) + 12
		chunkEnd := pos + totalChunkSize

		if chunkEnd > len(data) {
			// 块大小超出文件范围，复制剩余数据
			result = append(result, data[pos:]...)
			break
		}

		// 检查是否是文本块（tEXt, zTXt, iTXt）
		chunkTypeStr := string(chunkType)
		if chunkTypeStr == "tEXt" || chunkTypeStr == "zTXt" || chunkTypeStr == "iTXt" {
			// 跳过文本块
			pos = chunkEnd
			continue
		}

		// 保留非文本块
		result = append(result, data[pos:chunkEnd]...)
		pos = chunkEnd
	}

	return result
}

// insertPNGTextChunks 在PNG文件中插入文本块
func (m *ImageModifier) insertPNGTextChunks(data []byte, chunks [][]byte) []byte {
	// 查找IEND块的位置
	iendPos := m.findPNGIENDChunk(data)
	if iendPos == -1 {
		iendPos = len(data)
	}

	// 计算所有文本块的总大小
	totalChunkSize := 0
	for _, chunk := range chunks {
		totalChunkSize += len(chunk)
	}

	// 构造新的PNG数据
	result := make([]byte, 0, len(data)+totalChunkSize)
	result = append(result, data[:iendPos]...) // IEND之前的数据

	// 添加所有文本块
	for _, chunk := range chunks {
		result = append(result, chunk...)
	}

	result = append(result, data[iendPos:]...) // IEND块

	return result
}

// getPNGMetadata 获取PNG图片的文本元数据
func (m *ImageModifier) getPNGMetadata(imagePath string) (*ImageMetadata, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// PNG文件必须以PNG签名开头
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return nil, fmt.Errorf("不是有效的PNG文件")
	}

	metadata := &ImageMetadata{}
	pos := 8 // 跳过PNG签名

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

		// 检查是否是文本块
		if bytes.Equal(chunkType, []byte("tEXt")) {
			// 解析tEXt块
			if pos+int(chunkLen) <= len(data) {
				textData := data[pos : pos+int(chunkLen)]
				m.parsePNGTextChunk(textData, metadata)
			}
		}

		// 跳过数据和CRC
		pos += int(chunkLen) + 4

		// 检查是否是IEND块
		if bytes.Equal(chunkType, []byte("IEND")) {
			break
		}
	}

	return metadata, nil
}

// parsePNGTextChunk 解析PNG文本块
func (m *ImageModifier) parsePNGTextChunk(data []byte, metadata *ImageMetadata) {
	// tEXt块格式：关键字\0文本
	parts := bytes.SplitN(data, []byte{0}, 2)
	if len(parts) != 2 {
		return
	}

	keyword := string(parts[0])
	text := string(parts[1])

	// 根据关键字设置对应的元数据字段
	switch strings.ToLower(keyword) {
	case "author", "artist":
		metadata.Artist = text
	case "copyright":
		metadata.Copyright = text
	case "description":
		metadata.Description = text
	case "creation time", "datetime":
		if t, err := time.Parse("2006-01-02T15:04:05Z", text); err == nil {
			metadata.DateTime = &t
		} else if t, err := time.Parse("2006:01:02 15:04:05", text); err == nil {
			metadata.DateTime = &t
		}
	case "location":
		metadata.Location = text
	case "camera make", "make":
		metadata.CameraMake = text
	case "camera model", "model":
		metadata.CameraModel = text
	case "software":
		metadata.Software = text
	case "image width", "width":
		if w, err := strconv.Atoi(text); err == nil {
			metadata.ImageWidth = w
		}
	case "image height", "height":
		if h, err := strconv.Atoi(text); err == nil {
			metadata.ImageHeight = h
		}
	}
}
