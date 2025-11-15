package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StorageService 儲存服務
type StorageService struct {
	basePath string
}

// NewStorageService 建立新的儲存服務
func NewStorageService(basePath string) *StorageService {
	return &StorageService{
		basePath: basePath,
	}
}

// SaveAPIData 儲存 API 接收的資料
func (s *StorageService) SaveAPIData(data string) (string, error) {
	return s.saveData("API_print", data)
}

// SaveMQTTData 儲存 MQTT 接收的資料
func (s *StorageService) SaveMQTTData(data string) (string, error) {
	return s.saveData("MQTT_print", data)
}

// saveData 儲存資料到指定類型的資料夾
func (s *StorageService) saveData(dataType string, data string) (string, error) {
	now := time.Now()

	// 建立資料夾路徑: basePath/dataType/年_月_日
	dateFolder := now.Format("2006_01_02")
	folderPath := filepath.Join(s.basePath, dataType, dateFolder)

	// 確保資料夾存在
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", fmt.Errorf("建立資料夾失敗: %v", err)
	}

	// 檔案名稱: 時_分_秒.tspl
	fileName := now.Format("15_04_05") + ".tspl"
	filePath := filepath.Join(folderPath, fileName)

	// 如果檔案已存在,加上毫秒避免衝突
	if _, err := os.Stat(filePath); err == nil {
		fileName = now.Format("15_04_05") + "_" + fmt.Sprintf("%03d", now.Nanosecond()/1000000) + ".tspl"
		filePath = filepath.Join(folderPath, fileName)
	}

	// 寫入檔案
	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		return "", fmt.Errorf("寫入檔案失敗: %v", err)
	}

	return filePath, nil
}

// GetRecentFiles 取得最近的檔案列表
func (s *StorageService) GetRecentFiles(dataType string, limit int) ([]FileInfo, error) {
	basePath := filepath.Join(s.basePath, dataType)

	var files []FileInfo

	// 遍歷日期資料夾
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".tspl" {
			files = append(files, FileInfo{
				Path:     path,
				Name:     info.Name(),
				Size:     info.Size(),
				ModTime:  info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("讀取檔案列表失敗: %v", err)
	}

	// 按修改時間排序（最新的在前）
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i].ModTime.Before(files[j].ModTime) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	// 限制數量
	if limit > 0 && len(files) > limit {
		files = files[:limit]
	}

	return files, nil
}

// FileInfo 檔案資訊
type FileInfo struct {
	Path    string    `json:"path"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}
