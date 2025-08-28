package imagesha1

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageMetadata 图片元数据结构
type ImageMetadata struct {
	// 基本信息
	Artist      string // 作者/创作者
	Copyright   string // 版权信息
	Description string // 图片描述

	// 拍摄信息
	DateTime    *time.Time // 拍摄时间
	Location    string     // 拍摄地点
	CameraMake  string     // 相机制造商
	CameraModel string     // 相机型号

	// 技术参数
	Software    string // 处理软件
	ImageWidth  int    // 图片宽度 (仅PNG文本块使用)
	ImageHeight int    // 图片高度 (仅PNG文本块使用)
}

// MetadataModifier 元数据修改器接口
type MetadataModifier interface {
	// ModifyImageMetadata 修改图片元数据
	ModifyImageMetadata(imagePath string, metadata *ImageMetadata) (string, error)

	// GetImageMetadata 获取图片元数据
	GetImageMetadata(imagePath string) (*ImageMetadata, error)
}

// 扩展现有的ImageSHA1Modifier以支持元数据修改
// ModifyImageMetadata 通过修改元数据来改变图片SHA1
func (m *ImageSHA1Modifier) ModifyImageMetadata(imagePath string, metadata *ImageMetadata) (string, error) {
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
		modifiedData, err = m.modifyJPEGMetadata(originalData, metadata)
	case ".png":
		modifiedData, err = m.modifyPNGMetadata(originalData, metadata)
	default:
		return "", fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return "", fmt.Errorf("修改图片元数据失败: %v", err)
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

// GetImageMetadata 获取图片的元数据信息
func (m *ImageSHA1Modifier) GetImageMetadata(imagePath string) (*ImageMetadata, error) {
	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("图片文件不存在: %s", imagePath)
	}

	// 根据文件扩展名确定图片格式
	ext := strings.ToLower(filepath.Ext(imagePath))

	switch ext {
	case ".jpg", ".jpeg":
		return m.getJPEGMetadata(imagePath)
	case ".png":
		return m.getPNGMetadata(imagePath)
	default:
		return nil, fmt.Errorf("不支持的图片格式: %s", ext)
	}
}
