package imagesha1

import (
	"bytes"
	"fmt"
	"image/jpeg"
)

// modifyJPEGSHA1 修改JPEG图片的SHA1值
// 通过在JPEG的Comment段中添加随机数据来改变SHA1，不影响图片显示
func (m *ImageSHA1Modifier) modifyJPEGSHA1(data []byte) ([]byte, error) {
	// 解码JPEG图片
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解码JPEG图片失败: %v", err)
	}

	// 生成随机字节作为注释
	randomComment := m.generateRandomBytes(16)

	// 创建缓冲区重新编码图片
	var buf bytes.Buffer

	// 使用默认质量重新编码，并在编码过程中插入随机数据
	// 这里我们通过修改编码选项来改变输出
	options := &jpeg.Options{
		Quality: 95, // 使用高质量以减少视觉差异
	}

	err = jpeg.Encode(&buf, img, options)
	if err != nil {
		return nil, fmt.Errorf("重新编码JPEG图片失败: %v", err)
	}

	// 在JPEG数据中插入随机注释段
	// JPEG格式允许插入注释段而不影响图片显示
	modifiedData := m.insertJPEGComment(buf.Bytes(), randomComment)

	return modifiedData, nil
}

// insertJPEGComment 在JPEG数据中插入注释段
func (m *ImageSHA1Modifier) insertJPEGComment(data []byte, comment []byte) []byte {
	// JPEG文件格式：
	// FF D8 (SOI) ... 各种段 ... FF D9 (EOI)
	// 注释段格式：FF FE [长度高字节] [长度低字节] [注释数据]

	if len(data) < 4 {
		return data
	}

	// 寻找插入位置（在SOI之后，第一个段之前）
	insertPos := 2 // 跳过SOI标记 (FF D8)

	// 构造注释段
	commentLength := len(comment) + 2 // 注释数据长度 + 长度字段本身
	commentSegment := make([]byte, 0, commentLength+2)
	commentSegment = append(commentSegment, 0xFF, 0xFE)                                       // 注释段标记
	commentSegment = append(commentSegment, byte(commentLength>>8), byte(commentLength&0xFF)) // 长度
	commentSegment = append(commentSegment, comment...)                                       // 注释数据

	// 构造新的JPEG数据
	result := make([]byte, 0, len(data)+len(commentSegment))
	result = append(result, data[:insertPos]...) // SOI
	result = append(result, commentSegment...)   // 注释段
	result = append(result, data[insertPos:]...) // 剩余数据

	return result
}
