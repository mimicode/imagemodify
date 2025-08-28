package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mimicode/imagemodify"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("使用方法: go run example/pixel_modify_example.go <图片文件路径> <模式>")
		fmt.Println()
		fmt.Println("模式:")
		fmt.Println("  random   - 随机数据修改SHA1（原有方式）")
		fmt.Println("  pixel    - 像素微调修改SHA1（新方式）")
		fmt.Println("  metadata - 元数据修改SHA1")
		fmt.Println()
		fmt.Println("支持的格式: JPEG (.jpg, .jpeg), PNG (.png)")
		os.Exit(1)
	}

	imagePath := os.Args[1]
	mode := os.Args[2]

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", imagePath)
	}

	// 创建图片修改器
	modifier := imagemodify.NewImageModifier()

	// 获取原始SHA1值
	originalSHA1, err := modifier.GetImageSHA1(imagePath)
	if err != nil {
		log.Fatalf("获取原始SHA1失败: %v", err)
	}

	fmt.Printf("📋 图片信息: %s\n", filepath.Base(imagePath))
	fmt.Printf("🔍 原始SHA1: %s\n", originalSHA1)

	var newSHA1 string
	var methodName string

	switch mode {
	case "random":
		methodName = "随机数据修改"
		newSHA1, err = modifier.ModifyImageSHA1(imagePath)
	case "pixel":
		methodName = "像素微调修改"
		newSHA1, err = modifier.ModifyImageSHA1ByPixel(imagePath)
	case "metadata":
		methodName = "元数据修改"
		// 创建简单的元数据
		metadata := &imagemodify.ImageMetadata{
			Artist:      "测试摄影师",
			Copyright:   "© 2024 测试版权",
			Description: "像素微调测试图片",
			Software:    "ImageModify Pixel Tweaker v1.0",
		}
		newSHA1, err = modifier.ModifyImageMetadata(imagePath, metadata)
	default:
		log.Fatalf("未知模式: %s", mode)
	}

	if err != nil {
		log.Fatalf("%s失败: %v", methodName, err)
	}

	fmt.Printf("🔧 修改方式: %s\n", methodName)
	fmt.Printf("✨ 新的SHA1: %s\n", newSHA1)

	// 验证SHA1确实发生了变化
	if originalSHA1 == newSHA1 {
		fmt.Println("⚠️  警告: SHA1值未发生变化")
	} else {
		fmt.Println("✅ SHA1值已成功修改")

		// 计算SHA1差异位数
		diffCount := 0
		for i := 0; i < len(originalSHA1) && i < len(newSHA1); i++ {
			if originalSHA1[i] != newSHA1[i] {
				diffCount++
			}
		}
		fmt.Printf("📊 SHA1差异字符数: %d/%d\n", diffCount, len(originalSHA1))
	}

	// 获取文件信息
	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
	} else {
		fmt.Printf("💾 文件大小: %d 字节\n", fileInfo.Size())
		fmt.Printf("📁 文件格式: %s\n", filepath.Ext(imagePath))
	}

	fmt.Println("\n🎯 测试说明:")
	fmt.Println("- random模式: 在图片中插入随机注释/文本块")
	fmt.Println("- pixel模式: 微调边缘像素的亮度值（±2级别）")
	fmt.Println("- metadata模式: 修改图片的元数据信息")
	fmt.Println("\n所有方式都不会影响图片的视觉显示效果！")
}
