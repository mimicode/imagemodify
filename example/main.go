package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mimicode/imagemodify"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run example/main.go <图片文件路径>")
		fmt.Println("支持的格式: JPEG (.jpg, .jpeg), PNG (.png)")
		os.Exit(1)
	}

	imagePath := os.Args[1]

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", imagePath)
	}

	// 创建SHA1修改器
	modifier := imagemodify.NewImageModifier()

	// 获取原始SHA1值
	originalSHA1, err := modifier.GetImageSHA1(imagePath)
	if err != nil {
		log.Fatalf("获取原始SHA1失败: %v", err)
	}
	fmt.Printf("原始SHA1: %s\n", originalSHA1)

	// 修改图片SHA1
	newSHA1, err := modifier.ModifyImageSHA1(imagePath)
	if err != nil {
		log.Fatalf("修改SHA1失败: %v", err)
	}
	fmt.Printf("新的SHA1: %s\n", newSHA1)

	// 验证SHA1确实发生了变化
	if originalSHA1 == newSHA1 {
		fmt.Println("警告: SHA1值未发生变化")
	} else {
		fmt.Println("✓ SHA1值已成功修改")
	}

	// 获取文件信息
	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
	} else {
		fmt.Printf("文件大小: %d 字节\n", fileInfo.Size())
		fmt.Printf("文件格式: %s\n", filepath.Ext(imagePath))
	}

	fmt.Println("\n可以再次运行此程序来验证每次都会生成不同的SHA1值")
}
