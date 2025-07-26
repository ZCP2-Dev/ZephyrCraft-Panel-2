package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
