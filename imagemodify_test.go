package imagemodify

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestJPEG 创建测试用的JPEG图片
func createTestJPEG(path string) error {
	// 创建一个简单的测试图片
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// 填充一些颜色
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if (x+y)%20 < 10 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // 红色
			} else {
				img.Set(x, y, color.RGBA{0, 255, 0, 255}) // 绿色
			}
		}
	}

	// 保存为JPEG
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
}

// createTestPNG 创建测试用的PNG图片
func createTestPNG(path string) error {
	// 创建一个简单的测试图片
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// 填充一些颜色
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if (x+y)%20 < 10 {
				img.Set(x, y, color.RGBA{0, 0, 255, 255}) // 蓝色
			} else {
				img.Set(x, y, color.RGBA{255, 255, 0, 255}) // 黄色
			}
		}
	}

	// 保存为PNG
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// TestNewImageModifier 测试创建修改器
func TestNewImageModifier(t *testing.T) {
	modifier := NewImageModifier()
	if modifier == nil {
		t.Fatal("NewImageModifier() 返回了 nil")
	}
}

// TestModifyJPEGSHA1 测试JPEG格式SHA1修改
func TestModifyJPEGSHA1(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testJPEG := filepath.Join(tempDir, "test.jpg")

	// 创建测试JPEG图片
	if err := createTestJPEG(testJPEG); err != nil {
		t.Fatalf("创建测试JPEG失败: %v", err)
	}

	modifier := NewImageModifier()

	// 获取原始SHA1
	originalSHA1, err := modifier.GetImageSHA1(testJPEG)
	if err != nil {
		t.Fatalf("获取原始SHA1失败: %v", err)
	}

	// 修改SHA1
	newSHA1, err := modifier.ModifyImageSHA1(testJPEG)
	if err != nil {
		t.Fatalf("修改JPEG SHA1失败: %v", err)
	}

	// 验证SHA1发生了变化
	if originalSHA1 == newSHA1 {
		t.Error("SHA1值未发生变化")
	}

	// 验证文件仍然是有效的JPEG
	data, err := os.ReadFile(testJPEG)
	if err != nil {
		t.Fatalf("读取修改后的文件失败: %v", err)
	}

	_, err = jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		t.Errorf("修改后的文件不是有效的JPEG: %v", err)
	}
}

// TestModifyPNGSHA1 测试PNG格式SHA1修改
func TestModifyPNGSHA1(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testPNG := filepath.Join(tempDir, "test.png")

	// 创建测试PNG图片
	if err := createTestPNG(testPNG); err != nil {
		t.Fatalf("创建测试PNG失败: %v", err)
	}

	modifier := NewImageModifier()

	// 获取原始SHA1
	originalSHA1, err := modifier.GetImageSHA1(testPNG)
	if err != nil {
		t.Fatalf("获取原始SHA1失败: %v", err)
	}

	// 修改SHA1
	newSHA1, err := modifier.ModifyImageSHA1(testPNG)
	if err != nil {
		t.Fatalf("修改PNG SHA1失败: %v", err)
	}

	// 验证SHA1发生了变化
	if originalSHA1 == newSHA1 {
		t.Error("SHA1值未发生变化")
	}

	// 验证文件仍然是有效的PNG
	data, err := os.ReadFile(testPNG)
	if err != nil {
		t.Fatalf("读取修改后的文件失败: %v", err)
	}

	_, err = png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Errorf("修改后的文件不是有效的PNG: %v", err)
	}
}

// TestMultipleModifications 测试多次修改产生不同的SHA1
func TestMultipleModifications(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testJPEG := filepath.Join(tempDir, "test_multiple.jpg")

	// 创建测试JPEG图片
	if err := createTestJPEG(testJPEG); err != nil {
		t.Fatalf("创建测试JPEG失败: %v", err)
	}

	modifier := NewImageModifier()

	// 记录多次修改的SHA1值
	var sha1Values []string

	for i := 0; i < 3; i++ {
		newSHA1, err := modifier.ModifyImageSHA1(testJPEG)
		if err != nil {
			t.Fatalf("第%d次修改失败: %v", i+1, err)
		}
		sha1Values = append(sha1Values, newSHA1)
	}

	// 验证每次的SHA1都不相同
	for i := 0; i < len(sha1Values); i++ {
		for j := i + 1; j < len(sha1Values); j++ {
			if sha1Values[i] == sha1Values[j] {
				t.Errorf("第%d次和第%d次修改产生了相同的SHA1: %s", i+1, j+1, sha1Values[i])
			}
		}
	}
}

// TestUnsupportedFormat 测试不支持的格式
func TestUnsupportedFormat(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// 创建一个文本文件
	err := os.WriteFile(testFile, []byte("这不是图片文件"), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	modifier := NewImageModifier()

	// 尝试修改不支持的格式
	_, err = modifier.ModifyImageSHA1(testFile)
	if err == nil {
		t.Error("应该返回不支持格式的错误")
	}
}

// TestNonExistentFile 测试不存在的文件
func TestNonExistentFile(t *testing.T) {
	modifier := NewImageModifier()

	// 尝试修改不存在的文件
	_, err := modifier.ModifyImageSHA1("non_existent_file.jpg")
	if err == nil {
		t.Error("应该返回文件不存在的错误")
	}
}

// TestModifyJPEGPixelSHA1 测试JPEG格式像素微调SHA1修改
func TestModifyJPEGPixelSHA1(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testJPEG := filepath.Join(tempDir, "test_pixel.jpg")

	// 创建测试JPEG图片
	if err := createTestJPEG(testJPEG); err != nil {
		t.Fatalf("创建测试JPEG失败: %v", err)
	}

	modifier := NewImageModifier()

	// 获取原始SHA1
	originalSHA1, err := modifier.GetImageSHA1(testJPEG)
	if err != nil {
		t.Fatalf("获取原始SHA1失败: %v", err)
	}

	// 通过像素微调修改SHA1
	newSHA1, err := modifier.ModifyImageSHA1ByPixel(testJPEG)
	if err != nil {
		t.Fatalf("像素微调修改JPEG SHA1失败: %v", err)
	}

	// 验证SHA1发生了变化
	if originalSHA1 == newSHA1 {
		t.Error("SHA1值未发生变化")
	}

	// 验证文件仍然是有效的JPEG
	data, err := os.ReadFile(testJPEG)
	if err != nil {
		t.Fatalf("读取修改后的文件失败: %v", err)
	}

	_, err = jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		t.Errorf("修改后的文件不是有效的JPEG: %v", err)
	}
}

// TestModifyPNGPixelSHA1 测试PNG格式像素微调SHA1修改
func TestModifyPNGPixelSHA1(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testPNG := filepath.Join(tempDir, "test_pixel.png")

	// 创建测试PNG图片
	if err := createTestPNG(testPNG); err != nil {
		t.Fatalf("创建测试PNG失败: %v", err)
	}

	modifier := NewImageModifier()

	// 获取原始SHA1
	originalSHA1, err := modifier.GetImageSHA1(testPNG)
	if err != nil {
		t.Fatalf("获取原始SHA1失败: %v", err)
	}

	// 通过像素微调修改SHA1
	newSHA1, err := modifier.ModifyImageSHA1ByPixel(testPNG)
	if err != nil {
		t.Fatalf("像素微调修改PNG SHA1失败: %v", err)
	}

	// 验证SHA1发生了变化
	if originalSHA1 == newSHA1 {
		t.Error("SHA1值未发生变化")
	}

	// 验证文件仍然是有效的PNG
	data, err := os.ReadFile(testPNG)
	if err != nil {
		t.Fatalf("读取修改后的文件失败: %v", err)
	}

	_, err = png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Errorf("修改后的文件不是有效的PNG: %v", err)
	}
}

// TestMultiplePixelModifications 测试多次像素微调产生不同的SHA1
func TestMultiplePixelModifications(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testJPEG := filepath.Join(tempDir, "test_multiple_pixel.jpg")

	// 创建测试JPEG图片
	if err := createTestJPEG(testJPEG); err != nil {
		t.Fatalf("创建测试JPEG失败: %v", err)
	}

	modifier := NewImageModifier()

	// 记录多次修改的SHA1值
	var sha1Values []string

	for i := 0; i < 3; i++ {
		newSHA1, err := modifier.ModifyImageSHA1ByPixel(testJPEG)
		if err != nil {
			t.Fatalf("第%d次像素微调失败: %v", i+1, err)
		}
		sha1Values = append(sha1Values, newSHA1)
	}

	// 验证每次的SHA1都不相同
	for i := 0; i < len(sha1Values); i++ {
		for j := i + 1; j < len(sha1Values); j++ {
			if sha1Values[i] == sha1Values[j] {
				t.Errorf("第%d次和第%d次像素微调产生了相同的SHA1: %s", i+1, j+1, sha1Values[i])
			}
		}
	}
}
