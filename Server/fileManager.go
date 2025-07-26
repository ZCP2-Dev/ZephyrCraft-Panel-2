package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileInfo 文件信息结构
type FileInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	IsDir       bool      `json:"isDir"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"modTime"`
	Permissions string    `json:"permissions"`
	IsEditable  bool      `json:"isEditable"`
	Content     string    `json:"content,omitempty"`
}

// FileManager 文件管理器
type FileManager struct {
	rootPath string
}

// NewFileManager 创建新的文件管理器
func NewFileManager(rootPath string) *FileManager {
	// 确保根路径是绝对路径
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		log.Printf("[ERROR]无法获取绝对路径: %v", err)
		absPath = rootPath
	}
	return &FileManager{
		rootPath: absPath,
	}
}

// GetFileList 获取指定路径下的文件列表
func (fm *FileManager) GetFileList(path string) ([]FileInfo, error) {
	// 构建完整路径
	fullPath := filepath.Join(fm.rootPath, path)

	// 安全检查：确保路径在根目录内
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return nil, fmt.Errorf("访问路径超出允许范围")
	}

	// 检查路径是否存在
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	// 如果是文件，返回文件信息
	if !info.IsDir() {
		return []FileInfo{fm.getFileInfo(fullPath, info)}, nil
	}

	// 如果是目录，读取目录内容
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		entryPath := filepath.Join(fullPath, entry.Name())
		entryInfo, err := entry.Info()
		if err != nil {
			continue // 跳过无法获取信息的文件
		}
		files = append(files, fm.getFileInfo(entryPath, entryInfo))
	}

	return files, nil
}

// getFileInfo 获取单个文件信息
func (fm *FileManager) getFileInfo(fullPath string, info os.FileInfo) FileInfo {
	// 计算相对路径
	relPath, _ := filepath.Rel(fm.rootPath, fullPath)
	if relPath == "." {
		relPath = ""
	}

	// 判断是否可编辑
	isEditable := fm.isEditableFile(info.Name())

	return FileInfo{
		Name:        info.Name(),
		Path:        relPath,
		IsDir:       info.IsDir(),
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		Permissions: info.Mode().String(),
		IsEditable:  isEditable,
	}
}

// isEditableFile 判断文件是否可编辑
func (fm *FileManager) isEditableFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	editableExts := []string{
		".json", ".yaml", ".yml", ".ini", ".cfg", ".conf",
		".txt", ".log", ".md", ".xml", ".properties", ".toml",
		".env", ".sh", ".bat", ".ps1", ".py", ".js", ".ts",
		".vue", ".html", ".css", ".scss", ".less",
	}

	for _, ext2 := range editableExts {
		if ext == ext2 {
			return true
		}
	}
	return false
}

// ReadFile 读取文件内容
func (fm *FileManager) ReadFile(path string) (string, error) {
	fullPath := filepath.Join(fm.rootPath, path)

	// 安全检查
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return "", fmt.Errorf("访问路径超出允许范围")
	}

	fmt.Println(fullPath)

	// 检查文件是否存在且可编辑
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", fmt.Errorf("无法读取目录内容")
	}

	if !fm.isEditableFile(info.Name()) {
		return "", fmt.Errorf("文件类型不支持编辑")
	}

	// 读取文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// WriteFile 写入文件内容
func (fm *FileManager) WriteFile(path string, content string) error {
	fullPath := filepath.Join(fm.rootPath, path)

	// 安全检查
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return fmt.Errorf("访问路径超出允许范围")
	}

	// 检查文件是否可编辑
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("无法写入目录")
	}

	if !fm.isEditableFile(info.Name()) {
		return fmt.Errorf("文件类型不支持编辑")
	}

	// 写入文件内容
	return os.WriteFile(fullPath, []byte(content), 0644)
}

// CreateDirectory 创建目录
func (fm *FileManager) CreateDirectory(path string) error {
	fullPath := filepath.Join(fm.rootPath, path)

	// 安全检查
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return fmt.Errorf("访问路径超出允许范围")
	}

	return os.MkdirAll(fullPath, 0755)
}

// DeleteFile 删除文件或目录
func (fm *FileManager) DeleteFile(path string) error {
	fullPath := filepath.Join(fm.rootPath, path)

	// 安全检查
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return fmt.Errorf("访问路径超出允许范围")
	}

	// 防止删除根目录
	if fullPath == fm.rootPath {
		return fmt.Errorf("无法删除根目录")
	}

	return os.RemoveAll(fullPath)
}

// RenameFile 重命名文件或目录
func (fm *FileManager) RenameFile(oldPath, newPath string) error {
	oldFullPath := filepath.Join(fm.rootPath, oldPath)
	newFullPath := filepath.Join(fm.rootPath, newPath)

	// 安全检查
	if !strings.HasPrefix(oldFullPath, fm.rootPath) || !strings.HasPrefix(newFullPath, fm.rootPath) {
		return fmt.Errorf("访问路径超出允许范围")
	}

	return os.Rename(oldFullPath, newFullPath)
}

// GetFileSize 获取文件大小
func (fm *FileManager) GetFileSize(path string) (int64, error) {
	fullPath := filepath.Join(fm.rootPath, path)

	// 安全检查
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return 0, fmt.Errorf("访问路径超出允许范围")
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// ReadFileRaw 读取任意文件的二进制内容
func (fm *FileManager) ReadFileRaw(path string) ([]byte, error) {
	fullPath := filepath.Join(fm.rootPath, path)
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return nil, fmt.Errorf("访问路径超出允许范围")
	}
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("无法读取目录内容")
	}
	return os.ReadFile(fullPath)
}

// WriteFileRaw 写入任意文件的二进制内容
func (fm *FileManager) WriteFileRaw(path string, data []byte) error {
	fullPath := filepath.Join(fm.rootPath, path)
	if !strings.HasPrefix(fullPath, fm.rootPath) {
		return fmt.Errorf("访问路径超出允许范围")
	}
	return os.WriteFile(fullPath, data, 0644)
}

// CreateZipFile 创建zip文件
func (fm *FileManager) CreateZipFile(filesToZip []string, zipFileName string) error {
	// 构建zip文件的完整路径
	zipPath := filepath.Join(fm.rootPath, zipFileName)
	if !strings.HasSuffix(zipPath, ".zip") {
		zipPath += ".zip"
	}

	// 检查zip文件路径是否在允许范围内
	if !strings.HasPrefix(zipPath, fm.rootPath) {
		return fmt.Errorf("zip文件路径超出允许范围")
	}

	// 创建zip文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("创建zip文件失败: %v", err)
	}
	defer zipFile.Close()

	// 创建zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历要压缩的文件
	for _, filePath := range filesToZip {
		// 构建文件的完整路径
		fullPath := filepath.Join(fm.rootPath, filePath)

		// 检查文件路径是否在允许范围内
		if !strings.HasPrefix(fullPath, fm.rootPath) {
			continue // 跳过超出范围的文件
		}

		// 获取文件信息
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			continue // 跳过不存在的文件
		}

		if fileInfo.IsDir() {
			// 如果是目录，递归添加目录中的所有文件
			err = fm.addDirectoryToZip(zipWriter, fullPath, filePath)
			if err != nil {
				log.Printf("添加目录到zip失败 %s: %v", filePath, err)
			}
		} else {
			// 如果是文件，直接添加
			err = fm.addFileToZip(zipWriter, fullPath, filePath)
			if err != nil {
				log.Printf("添加文件到zip失败 %s: %v", filePath, err)
			}
		}
	}

	return nil
}

// addFileToZip 添加单个文件到zip
func (fm *FileManager) addFileToZip(zipWriter *zip.Writer, fullPath, relativePath string) error {
	// 打开文件
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建zip文件条目
	zipEntry, err := zipWriter.Create(relativePath)
	if err != nil {
		return err
	}

	// 复制文件内容到zip
	_, err = io.Copy(zipEntry, file)
	return err
}

// addDirectoryToZip 递归添加目录到zip
func (fm *FileManager) addDirectoryToZip(zipWriter *zip.Writer, fullPath, relativePath string) error {
	return filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(fm.rootPath, path)
		if err != nil {
			return err
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			// 为目录创建条目（zip需要目录条目）
			_, err = zipWriter.Create(relPath + "/")
			return err
		} else {
			// 添加文件
			return fm.addFileToZip(zipWriter, path, relPath)
		}
	})
}

// 读取server.properties为map
func (fm *FileManager) ReadServerProperties() (map[string]string, error) {
	props := make(map[string]string)
	propPath := filepath.Join(fm.rootPath, "server.properties")
	file, err := os.Open(propPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			props[key] = val
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return props, nil
}

// 写入server.properties（全量覆盖）
func (fm *FileManager) WriteServerProperties(props map[string]string) error {
	propPath := filepath.Join(fm.rootPath, "server.properties")
	file, err := os.Create(propPath)
	if err != nil {
		return err
	}
	defer file.Close()
	for k, v := range props {
		if _, err := file.WriteString(k + "=" + v + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// 获取max-players
func (fm *FileManager) GetMaxPlayers() (int, error) {
	props, err := fm.ReadServerProperties()
	if err != nil {
		return 0, err
	}
	val, ok := props["max-players"]
	if !ok {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// 设置max-players
func (fm *FileManager) SetMaxPlayers(n int) error {
	props, err := fm.ReadServerProperties()
	if err != nil {
		return err
	}
	props["max-players"] = strconv.Itoa(n)
	return fm.WriteServerProperties(props)
}

// 全局文件管理器实例
var globalFileManager *FileManager

// InitFileManager 初始化全局文件管理器实例
func InitFileManager(root string) {
	globalFileManager = NewFileManager(root)
}

// GetFileManager 获取全局文件管理器实例
func GetFileManager() *FileManager {
	if globalFileManager == nil {
		// 默认使用当前工作目录
		globalFileManager = NewFileManager(".")
	}
	return globalFileManager
}
