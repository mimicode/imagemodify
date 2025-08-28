package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mimicode/imagemodify"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法:")
		fmt.Println("  go run example/metadata_example.go <图片文件路径> [操作类型]")
		fmt.Println("")
		fmt.Println("操作类型:")
		fmt.Println("  show     - 显示当前元数据 (默认)")
		fmt.Println("  modify   - 修改元数据示例")
		fmt.Println("  custom   - 自定义元数据修改")
		fmt.Println("")
		fmt.Println("支持的格式: JPEG (.jpg, .jpeg), PNG (.png)")
		os.Exit(1)
	}

	imagePath := os.Args[1]
	operation := "show"
	if len(os.Args) > 2 {
		operation = os.Args[2]
	}

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", imagePath)
	}

	// 创建SHA1修改器
	modifier := imagemodify.NewImageSHA1Modifier()

	switch operation {
	case "show":
		showMetadata(modifier, imagePath)
	case "modify":
		modifyMetadataExample(modifier, imagePath)
	case "custom":
		customMetadataModify(modifier, imagePath)
	default:
		fmt.Printf("未知操作: %s\n", operation)
		os.Exit(1)
	}
}

// showMetadata 显示图片的当前元数据
func showMetadata(modifier *imagemodify.ImageSHA1Modifier, imagePath string) {
	fmt.Printf("📋 图片元数据信息: %s\n", imagePath)
	fmt.Println("" + strings.Repeat("=", 50))

	// 获取当前SHA1
	currentSHA1, err := modifier.GetImageSHA1(imagePath)
	if err != nil {
		log.Printf("获取SHA1失败: %v", err)
	} else {
		fmt.Printf("🔍 当前SHA1: %s\n\n", currentSHA1)
	}

	// 获取元数据
	metadata, err := modifier.GetImageMetadata(imagePath)
	if err != nil {
		log.Fatalf("获取元数据失败: %v", err)
	}

	// 显示元数据
	fmt.Printf("👤 作者/艺术家: %s\n", getDisplayValue(metadata.Artist))
	fmt.Printf("©️  版权信息: %s\n", getDisplayValue(metadata.Copyright))
	fmt.Printf("📝 图片描述: %s\n", getDisplayValue(metadata.Description))

	if metadata.DateTime != nil {
		fmt.Printf("📅 拍摄时间: %s\n", metadata.DateTime.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("📅 拍摄时间: %s\n", getDisplayValue(""))
	}

	fmt.Printf("📍 拍摄地点: %s\n", getDisplayValue(metadata.Location))
	fmt.Printf("📷 相机制造商: %s\n", getDisplayValue(metadata.CameraMake))
	fmt.Printf("📷 相机型号: %s\n", getDisplayValue(metadata.CameraModel))
	fmt.Printf("💻 处理软件: %s\n", getDisplayValue(metadata.Software))

	if metadata.ImageWidth > 0 && metadata.ImageHeight > 0 {
		fmt.Printf("📐 记录尺寸: %dx%d 像素\n", metadata.ImageWidth, metadata.ImageHeight)
	}
}

// modifyMetadataExample 修改元数据示例
func modifyMetadataExample(modifier *imagemodify.ImageSHA1Modifier, imagePath string) {
	fmt.Printf("🔧 修改图片元数据示例: %s\n", imagePath)
	fmt.Println("" + strings.Repeat("=", 50))

	// 获取原始SHA1
	originalSHA1, err := modifier.GetImageSHA1(imagePath)
	if err != nil {
		log.Fatalf("获取原始SHA1失败: %v", err)
	}
	fmt.Printf("原始SHA1: %s\n", originalSHA1)

	// 创建示例元数据
	now := time.Now()
	metadata := &imagemodify.ImageMetadata{
		Artist:      "张三摄影师",
		Copyright:   "© 2024 张三摄影工作室",
		Description: "这是一张通过ImageSHA1库修改的示例图片",
		DateTime:    &now,
		Location:    "北京市朝阳区",
		CameraMake:  "Canon",
		CameraModel: "EOS R5",
		Software:    "ImageSHA1 Metadata Modifier v1.0",
		ImageWidth:  1920,
		ImageHeight: 1080,
	}

	// 修改元数据
	newSHA1, err := modifier.ModifyImageMetadata(imagePath, metadata)
	if err != nil {
		log.Fatalf("修改元数据失败: %v", err)
	}

	fmt.Printf("新的SHA1: %s\n", newSHA1)

	if originalSHA1 == newSHA1 {
		fmt.Println("⚠️  警告: SHA1值未发生变化")
	} else {
		fmt.Println("✅ SHA1值已成功修改")
	}

	fmt.Println("\n📋 已设置的元数据:")
	fmt.Printf("  👤 作者: %s\n", metadata.Artist)
	fmt.Printf("  ©️  版权: %s\n", metadata.Copyright)
	fmt.Printf("  📝 描述: %s\n", metadata.Description)
	fmt.Printf("  📅 时间: %s\n", metadata.DateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  📍 地点: %s\n", metadata.Location)
	fmt.Printf("  📷 相机: %s %s\n", metadata.CameraMake, metadata.CameraModel)
	fmt.Printf("  💻 软件: %s\n", metadata.Software)
}

// customMetadataModify 自定义元数据修改
func customMetadataModify(modifier *imagemodify.ImageSHA1Modifier, imagePath string) {
	fmt.Printf("🎯 自定义元数据修改: %s\n", imagePath)
	fmt.Println("" + strings.Repeat("=", 50))

	// 获取原始SHA1
	originalSHA1, err := modifier.GetImageSHA1(imagePath)
	if err != nil {
		log.Fatalf("获取原始SHA1失败: %v", err)
	}
	fmt.Printf("原始SHA1: %s\n\n", originalSHA1)

	// 这里演示几种不同的修改方式
	examples := []struct {
		name     string
		metadata *imagemodify.ImageMetadata
	}{
		{
			name: "修改作者信息",
			metadata: &imagemodify.ImageMetadata{
				Artist: "李四摄影师",
			},
		},
		{
			name: "修改版权信息",
			metadata: &imagemodify.ImageMetadata{
				Copyright: "© 2024 某某版权所有",
			},
		},
		{
			name: "修改拍摄时间",
			metadata: &imagemodify.ImageMetadata{
				DateTime: func() *time.Time { t := time.Date(2024, 6, 15, 14, 30, 0, 0, time.Local); return &t }(),
			},
		},
		{
			name: "修改地点信息",
			metadata: &imagemodify.ImageMetadata{
				Location: "上海市浦东新区",
			},
		},
	}

	// 依次应用每个修改示例
	for i, example := range examples {
		fmt.Printf("%d. %s\n", i+1, example.name)

		newSHA1, err := modifier.ModifyImageMetadata(imagePath, example.metadata)
		if err != nil {
			log.Printf("  ❌ 修改失败: %v", err)
			continue
		}

		fmt.Printf("  新SHA1: %s\n", newSHA1)

		if i == 0 {
			if originalSHA1 == newSHA1 {
				fmt.Println("  ⚠️  警告: SHA1值未发生变化")
			} else {
				fmt.Println("  ✅ SHA1值已成功修改")
			}
		} else {
			fmt.Println("  ✅ SHA1值已更新")
		}

		originalSHA1 = newSHA1 // 更新用于下次比较
		fmt.Println()
	}

	fmt.Println("🎉 所有修改示例完成！")
}

// getDisplayValue 获取显示值，空值显示为 "(未设置)"
func getDisplayValue(value string) string {
	if value == "" {
		return "(未设置)"
	}
	return value
}
