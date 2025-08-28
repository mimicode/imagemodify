package imagemodify

import (
	"encoding/json"
	"fmt"
	"os"
)

// modifyJPEGMetadata 修改JPEG图片的元数据（通过注释段）
func (m *ImageModifier) modifyJPEGMetadata(data []byte, metadata *ImageMetadata) ([]byte, error) {
	// 移除现有的注释段
	cleanData := m.removeJPEGComments(data)

	// 将元数据序列化为JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("序列化元数据失败: %v", err)
	}

	// 在JPEG中插入包含元数据的注释段
	return m.insertJPEGComment(cleanData, metadataJSON), nil
}

// removeJPEGComments 移除现有的注释段
func (m *ImageModifier) removeJPEGComments(data []byte) []byte {
	if len(data) < 4 {
		return data
	}

	pos := 2 // 跳过SOI
	result := make([]byte, 0, len(data))
	result = append(result, data[:2]...) // 保留SOI

	for pos < len(data)-4 {
		// 检查段标记
		if data[pos] != 0xFF {
			// 不是段标记，复制剩余数据
			result = append(result, data[pos:]...)
			break
		}

		marker := data[pos+1]
		pos += 2

		// 检查是否是注释段 (0xFE)
		if marker == 0xFE {
			// 读取段长度并跳过注释段
			if pos+2 > len(data) {
				break
			}
			segmentLen := int(data[pos])<<8 | int(data[pos+1])
			pos += segmentLen
			continue
		}

		// 不是注释段，保留这个段
		result = append(result, 0xFF, marker)

		// 检查是否有长度字段的段
		if marker >= 0xE0 && marker <= 0xEF || marker == 0xFE {
			// 有长度字段
			if pos+2 > len(data) {
				break
			}
			segmentLen := int(data[pos])<<8 | int(data[pos+1])
			endPos := pos + segmentLen
			if endPos > len(data) {
				endPos = len(data)
			}
			result = append(result, data[pos:endPos]...)
			pos = endPos
		} else {
			// 没有长度字段的段，复制剩余数据
			result = append(result, data[pos:]...)
			break
		}
	}

	return result
}

// getJPEGMetadata 获取JPEG图片的元数据（从注释段）
func (m *ImageModifier) getJPEGMetadata(imagePath string) (*ImageMetadata, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 查找注释段
	comment := m.extractJPEGComment(data)
	if comment == nil {
		// 没有注释段，返回空元数据
		return &ImageMetadata{}, nil
	}

	// 尝试解析JSON格式的元数据
	var metadata ImageMetadata
	err = json.Unmarshal(comment, &metadata)
	if err != nil {
		// 不是JSON格式的注释，返回空元数据
		return &ImageMetadata{}, nil
	}

	return &metadata, nil
}

// extractJPEGComment 提取JPEG注释段内容
func (m *ImageModifier) extractJPEGComment(data []byte) []byte {
	if len(data) < 4 {
		return nil
	}

	pos := 2 // 跳过SOI

	for pos < len(data)-4 {
		// 检查段标记
		if data[pos] != 0xFF {
			break
		}

		marker := data[pos+1]
		pos += 2

		// 检查是否是注释段 (0xFE)
		if marker == 0xFE {
			// 读取段长度
			if pos+2 > len(data) {
				break
			}
			segmentLen := int(data[pos])<<8 | int(data[pos+1])
			pos += 2

			// 提取注释数据
			commentLen := segmentLen - 2 // 减去长度字段本身
			if pos+commentLen <= len(data) && commentLen > 0 {
				return data[pos : pos+commentLen]
			}
			return nil
		}

		// 跳过其他段
		if marker >= 0xE0 && marker <= 0xEF || marker == 0xFE {
			// 有长度字段
			if pos+2 > len(data) {
				break
			}
			segmentLen := int(data[pos])<<8 | int(data[pos+1])
			pos += segmentLen
		} else {
			// 没有长度字段的段
			break
		}
	}

	return nil
}
