package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run verify_image.go <图片文件路径>")
		os.Exit(1)
	}

	imagePath := os.Args[1]

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Printf("❌ 文件不存在: %s\n", imagePath)
		os.Exit(1)
	}

	// 打开文件
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Printf("❌ 无法打开文件: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// 根据扩展名解码图片
	ext := strings.ToLower(filepath.Ext(imagePath))
	var img image.Image

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			fmt.Printf("❌ JPEG解码失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ JPEG图片解码成功")

	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			fmt.Printf("❌ PNG解码失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ PNG图片解码成功")

	default:
		fmt.Printf("❌ 不支持的图片格式: %s\n", ext)
		os.Exit(1)
	}

	// 获取图片信息
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("📐 图片尺寸: %d x %d 像素\n", width, height)
	fmt.Printf("📁 文件格式: %s\n", ext)

	// 获取文件大小
	fileInfo, err := os.Stat(imagePath)
	if err == nil {
		fmt.Printf("💾 文件大小: %d 字节\n", fileInfo.Size())
	}

	// 验证像素数据（检查前几个像素点）
	fmt.Println("🎨 像素数据验证:")
	for y := 0; y < 3 && y < height; y++ {
		for x := 0; x < 3 && x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// RGBA值是16位的，转换为8位
			r8, g8, b8, a8 := r>>8, g>>8, b>>8, a>>8
			fmt.Printf("   像素[%d,%d]: RGBA(%d,%d,%d,%d)\n", x, y, r8, g8, b8, a8)
		}
	}

	fmt.Println("✅ 图片验证完成，图片文件完好无损！")
}
