# ImageModify - Go图片SHA1修改库

这是一个Go语言编写的图片SHA1值修改库，可以在不改变图片显示内容、宽高和格式的情况下修改图片文件的SHA1值。

## 功能特性

- ✅ 支持JPEG (.jpg, .jpeg) 和 PNG (.png) 格式
- ✅ 直接在原图上修改，不改变图片尺寸和格式
- ✅ 不影响图片内容显示
- ✅ 每次执行都会生成不同的SHA1值
- ✅ 支持三种修改方式：
  - **随机数据修改**：快速改变SHA1值
  - **像素微调修改**：通过微调边缘像素亮度改变SHA1（最精妙的方式）
  - **元数据修改**：通过修改有意义的元数据改变SHA1值
- ✅ 丰富的元数据支持：作者、版权、拍摄时间、地点、相机信息等
- ✅ 简单易用的API接口

## 工作原理

### JPEG格式
- **随机数据模式**：通过在JPEG文件中插入注释段（Comment Segment）来改变文件内容。
- **像素微调模式**：通过微调边缘像素的RGB值（±2级别）来改变图像数据。
- **元数据模式**：通过修改EXIF元数据信息来改变文件内容。

### PNG格式
- **随机数据模式**：通过在PNG文件中插入文本块（tEXt chunk）来改变文件内容。
- **像素微调模式**：通过微调边缘像素的RGB值（±2级别）来改变图像数据。
- **元数据模式**：通过在PNG文件中插入元数据文本块来改变文件内容。

所有方式都不会影响图片的显示效果和视觉质量。

## 安装使用

### 作为库使用

#### 1. 随机数据修改（快速模式）

```go
package main

import (
    "fmt"
    "github.com/mimicode/imagemodify"
)

func main() {
    // 创建图片修改器
    modifier := imagemodify.NewImageModifier()
    
    // 修改图片SHA1（使用随机数据）
    newSHA1, err := modifier.ModifyImageSHA1("path/to/your/image.jpg")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("新的SHA1值: %s\n", newSHA1)
}
```

#### 2. 像素微调修改（精妙模式）

```go
package main

import (
    "fmt"
    "github.com/mimicode/imagemodify"
)

func main() {
    // 创建图片修改器
    modifier := imagemodify.NewImageModifier()
    
    // 通过像素微调修改SHA1（最精妙的方式）
    newSHA1, err := modifier.ModifyImageSHA1ByPixel("path/to/your/image.jpg")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("新的SHA1值: %s\n", newSHA1)
}
```

#### 3. 元数据修改（智能模式）

```go
package main

import (
    "fmt"
    "time"
    "github.com/mimicode/imagemodify"
)

func main() {
    // 创建图片修改器
    modifier := imagemodify.NewImageModifier()
    
    // 设置元数据
    now := time.Now()
    metadata := &imagemodify.ImageMetadata{
        Artist:      "张三摄影师",
        Copyright:   "© 2024 张三摄影工作室",
        Description: "美丽的天空风景图",
        DateTime:    &now,
        Location:    "北京市朝阳区",
        CameraMake:  "Canon",
        CameraModel: "EOS R5",
        Software:    "Adobe Lightroom",
    }
    
    // 通过修改元数据来改变SHA1
    newSHA1, err := modifier.ModifyImageMetadata("path/to/your/image.jpg", metadata)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("新的SHA1值: %s\n", newSHA1)
    
    // 获取修改后的元数据
    retrievedMetadata, err := modifier.GetImageMetadata("path/to/your/image.jpg")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("作者: %s\n", retrievedMetadata.Artist)
    fmt.Printf("版权: %s\n", retrievedMetadata.Copyright)
}
```

### 运行示例程序

#### 随机数据修改示例
```bash
# 编译并运行示例
go run example/main.go path/to/your/image.jpg
```

#### 多模式统一示例（推荐）
```bash
# 随机数据修改模式
go run example/pixel_modify_example.go path/to/your/image.jpg random

# 像素微调修改模式（最精妙）
go run example/pixel_modify_example.go path/to/your/image.jpg pixel

# 元数据修改模式
go run example/pixel_modify_example.go path/to/your/image.jpg metadata
```

#### 元数据修改示例
```bash
# 显示当前元数据
go run example/metadata_example.go path/to/your/image.jpg show

# 修改元数据示例
go run example/metadata_example.go path/to/your/image.jpg modify

# 自定义元数据修改
go run example/metadata_example.go path/to/your/image.jpg custom
```

#### 验证图片完整性
```bash
# 验证修改后的图片能否正常显示
go run example/verify_image.go path/to/your/image.jpg
```

## API 文档

### ImageModifier

主要的图片SHA1修改器结构体。

#### 方法

##### `NewImageModifier() *ImageModifier`
创建新的图片SHA1修改器实例。

##### `ModifyImageSHA1(imagePath string) (string, error)`
使用随机数据修改指定路径图片的SHA1值。

**参数:**
- `imagePath`: 图片文件的完整路径

**返回值:**
- `string`: 修改后的SHA1值（十六进制字符串）
- `error`: 错误信息，如果操作成功则为nil

##### `ModifyImageSHA1ByPixel(imagePath string) (string, error)`
通过微调边缘像素亮度修改指定路径图片的SHA1值（最精妙的方式）。

**参数:**
- `imagePath`: 图片文件的完整路径

**返回值:**
- `string`: 修改后的SHA1值（十六进制字符串）
- `error`: 错误信息，如枟操作成功则为nil

##### `ModifyImageMetadata(imagePath string, metadata *ImageMetadata) (string, error)`
通过修改元数据来改变图片的SHA1值。

**参数:**
- `imagePath`: 图片文件的完整路径
- `metadata`: 要设置的元数据结构体

**返回值:**
- `string`: 修改后的SHA1值（十六进制字符串）
- `error`: 错误信息，如果操作成功则为nil

##### `GetImageSHA1(imagePath string) (string, error)`
获取指定图片文件的当前SHA1值。

**参数:**
- `imagePath`: 图片文件的完整路径

**返回值:**
- `string`: 当前的SHA1值（十六进制字符串）
- `error`: 错误信息，如果操作成功则为nil

##### `GetImageMetadata(imagePath string) (*ImageMetadata, error)`
获取指定图片文件的元数据信息。

**参数:**
- `imagePath`: 图片文件的完整路径

**返回值:**
- `*ImageMetadata`: 图片的元数据结构体
- `error`: 错误信息，如果操作成功则为nil

### ImageMetadata

图片元数据结构体，包含各种图片相关信息。

```go
type ImageMetadata struct {
    // 基本信息
    Artist      string    // 作者/创作者
    Copyright   string    // 版权信息
    Description string    // 图片描述
    
    // 拍摄信息
    DateTime     *time.Time // 拍摄时间
    Location     string     // 拍摄地点
    CameraMake   string     // 相机制造商
    CameraModel  string     // 相机型号
    
    // 技术参数
    Software     string     // 处理软件
    ImageWidth   int        // 图片宽度 (仅PNG文本块使用)
    ImageHeight  int        // 图片高度 (仅PNG文本块使用)
}
```

**使用说明:**
- 只需要设置您需要修改的字段，空字段将被忽略
- `DateTime` 字段使用指针，可以设置为 `nil` 表示不修改
- 不同格式的图片对元数据的支持略有不同

## 支持的格式

| 格式 | 扩展名 | 修改方式 |
|------|--------|----------|
| JPEG | .jpg, .jpeg | 插入注释段 |
| PNG  | .png | 插入文本块 |

## 注意事项

1. **文件备份**: 建议在修改重要图片前先进行备份
2. **文件权限**: 确保程序对目标文件有读写权限
3. **格式支持**: 目前仅支持JPEG和PNG格式
4. **文件完整性**: 修改后的文件保持原有的图片格式和显示效果

## 错误处理

库会返回以下类型的错误：

- 文件不存在
- 文件读取/写入权限问题
- 不支持的图片格式
- 图片解码失败
- SHA1修改失败

## 示例输出

### 像素微调模式（推荐）
```
📋 图片信息: example.jpg
🔍 原始SHA1: a1b2c3d4e5f6789012345678901234567890abcd
🔧 修改方式: 像素微调修改
✨ 新的SHA1: 1234567890abcdef1234567890abcdef12345678
✅ SHA1值已成功修改
📈 SHA1差异字符数: 38/40
💾 文件大小: 123456 字节
📁 文件格式: .jpg

🎯 所有方式都不会影响图片的视觉显示效果！
```

### 传统模式
```
原始SHA1: a1b2c3d4e5f6789012345678901234567890abcd
新的SHA1: 1234567890abcdef1234567890abcdef12345678
✓ SHA1值已成功修改
文件大小: 123456 字节
文件格式: .jpg
```

## 许可证

MIT License