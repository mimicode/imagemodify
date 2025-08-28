package imagemodify

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// ImageModifier 图片修改器
type ImageModifier struct {
}

// NewImageModifier 创建新的图片修改器
func NewImageModifier() *ImageModifier {
	return &ImageModifier{}
}

// ModifyImageSHA1 修改图片文件的SHA1值（随机数据模式）
// imagePath: 图片文件路径
// 返回: 修改后的SHA1值和错误信息
func (m *ImageModifier) ModifyImageSHA1(imagePath string) (string, error) {
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
		modifiedData = m.insertJPEGComment(originalData, m.generateRandomBytes(16))
	case ".png":
		modifiedData = m.insertPNGTextChunk(originalData, "Random", string(m.generateRandomBytes(32)))
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
func (m *ImageModifier) GetImageSHA1(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %v", err)
	}
	return fmt.Sprintf("%x", sha1.Sum(data)), nil
}

// ModifyImageSHA1ByPixel 通过微调边缘像素亮度来修改图片SHA1值
// imagePath: 图片文件路径
// 返回: 修改后的SHA1值和错误信息
func (m *ImageModifier) ModifyImageSHA1ByPixel(imagePath string) (string, error) {
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
		modifiedData, err = m.modifyJPEGPixel(originalData)
	case ".png":
		modifiedData, err = m.modifyPNGPixel(originalData)
	default:
		return "", fmt.Errorf("不支持的图片格式: %s", ext)
	}

	if err != nil {
		return "", fmt.Errorf("像素微调失败: %v", err)
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

// generateRandomBytes 生成随机字节
func (m *ImageModifier) generateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return bytes
}

// modifyJPEGPixel 通过微调像素修改JPEG图片
func (m *ImageModifier) modifyJPEGPixel(data []byte) ([]byte, error) {
	// 解码JPEG图片
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解码JPEG图片失败: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建一个可编辑的图片副本
	newImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	// 随机选择边缘像素进行微调
	edgePixels := m.getEdgePixels(width, height)
	if len(edgePixels) == 0 {
		return nil, fmt.Errorf("找不到边缘像素")
	}

	// 随机选择一个边缘像素
	randomBytes := m.generateRandomBytes(4)
	pixelIndex := int(randomBytes[0]) % len(edgePixels)
	selectedPixel := edgePixels[pixelIndex]

	// 获取原像素颜色
	origColor := newImg.RGBAAt(selectedPixel.X, selectedPixel.Y)

	// 微调亮度（对RGB值进行微小调整）
	adjustment := int8(randomBytes[1]%5) - 2 // -2 到 +2 的微调
	newR := m.clampUint8(int(origColor.R) + int(adjustment))
	newG := m.clampUint8(int(origColor.G) + int(adjustment))
	newB := m.clampUint8(int(origColor.B) + int(adjustment))

	// 设置新颜色
	newImg.SetRGBA(selectedPixel.X, selectedPixel.Y, color.RGBA{
		R: newR,
		G: newG,
		B: newB,
		A: origColor.A,
	})

	// 重新编码为JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 95})
	if err != nil {
		return nil, fmt.Errorf("重新编码JPEG失败: %v", err)
	}

	return buf.Bytes(), nil
}

// modifyPNGPixel 通过微调像素修改PNG图片
func (m *ImageModifier) modifyPNGPixel(data []byte) ([]byte, error) {
	// 解码PNG图片
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解码PNG图片失败: %v", err)
	}

	// 获取图片边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建一个可编辑的图片副本
	newImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	// 随机选择边缘像素进行微调
	edgePixels := m.getEdgePixels(width, height)
	if len(edgePixels) == 0 {
		return nil, fmt.Errorf("找不到边缘像素")
	}

	// 随机选择一个边缘像素
	randomBytes := m.generateRandomBytes(4)
	pixelIndex := int(randomBytes[0]) % len(edgePixels)
	selectedPixel := edgePixels[pixelIndex]

	// 获取原像素颜色
	origColor := newImg.RGBAAt(selectedPixel.X, selectedPixel.Y)

	// 微调亮度（对RGB值进行微小调整）
	adjustment := int8(randomBytes[1]%5) - 2 // -2 到 +2 的微调
	newR := m.clampUint8(int(origColor.R) + int(adjustment))
	newG := m.clampUint8(int(origColor.G) + int(adjustment))
	newB := m.clampUint8(int(origColor.B) + int(adjustment))

	// 设置新颜色
	newImg.SetRGBA(selectedPixel.X, selectedPixel.Y, color.RGBA{
		R: newR,
		G: newG,
		B: newB,
		A: origColor.A,
	})

	// 重新编码为PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, newImg)
	if err != nil {
		return nil, fmt.Errorf("重新编码PNG失败: %v", err)
	}

	return buf.Bytes(), nil
}

// PixelCoord 像素坐标
type PixelCoord struct {
	X, Y int
}

// getEdgePixels 获取边缘像素坐标列表
func (m *ImageModifier) getEdgePixels(width, height int) []PixelCoord {
	var edgePixels []PixelCoord

	// 上下边界
	for x := 0; x < width; x++ {
		edgePixels = append(edgePixels, PixelCoord{X: x, Y: 0})          // 上边界
		edgePixels = append(edgePixels, PixelCoord{X: x, Y: height - 1}) // 下边界
	}

	// 左右边界（避免重复角落像素）
	for y := 1; y < height-1; y++ {
		edgePixels = append(edgePixels, PixelCoord{X: 0, Y: y})         // 左边界
		edgePixels = append(edgePixels, PixelCoord{X: width - 1, Y: y}) // 右边界
	}

	return edgePixels
}

// clampUint8 将数值限制在 0-255 范围内
func (m *ImageModifier) clampUint8(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}
