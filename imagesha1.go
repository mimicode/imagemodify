package imagesha1

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ImageSHA1Modifier 图片SHA1修改器
type ImageSHA1Modifier struct {
}

// NewImageSHA1Modifier 创建新的图片SHA1修改器
func NewImageSHA1Modifier() *ImageSHA1Modifier {
	return &ImageSHA1Modifier{}
}

// ModifyImageSHA1 修改图片文件的SHA1值
// imagePath: 图片文件路径
// 返回: 修改后的SHA1值和错误信息
func (m *ImageSHA1Modifier) ModifyImageSHA1(imagePath string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("图片文件不存在: %s", imagePath)
	}

	// 读取原始文件
	originalData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %v", err)
	}

	// 计算原始SHA1
	originalSHA1 := fmt.Sprintf("%x", sha1.Sum(originalData))

	// 根据文件扩展名确定图片格式
	ext := strings.ToLower(filepath.Ext(imagePath))
	var modifiedData []byte

	switch ext {
	case ".jpg", ".jpeg":
		modifiedData, err = m.modifyJPEGSHA1(originalData)
	case ".png":
		modifiedData, err = m.modifyPNGSHA1(originalData)
	default:
		return "", fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return "", fmt.Errorf("修改图片SHA1失败: %v", err)
	}

	// 验证修改后的数据与原始数据不同
	newSHA1 := fmt.Sprintf("%x", sha1.Sum(modifiedData))
	if newSHA1 == originalSHA1 {
		return "", fmt.Errorf("SHA1修改失败，值未发生变化")
	}

	// 写回文件
	err = os.WriteFile(imagePath, modifiedData, 0644)
	if err != nil {
		return "", fmt.Errorf("写入修改后的图片失败: %v", err)
	}

	return newSHA1, nil
}

// GetImageSHA1 获取图片文件的SHA1值
func (m *ImageSHA1Modifier) GetImageSHA1(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %v", err)
	}
	return fmt.Sprintf("%x", sha1.Sum(data)), nil
}

// generateRandomBytes 生成随机字节
func (m *ImageSHA1Modifier) generateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return bytes
}
