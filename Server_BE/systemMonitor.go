package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMonitor 系统监控器
type SystemMonitor struct {
	lastCPUUsage float64
	lastCPUTime  time.Time
	outputBuffer []string
	startTime    string
}

// NewSystemMonitor 创建新的系统监控器
func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{
		lastCPUUsage: 0,
		lastCPUTime:  time.Now(),
		outputBuffer: make([]string, 0, 1000),
	}
}

// GetSystemInfo 获取系统状态信息
func (sm *SystemMonitor) GetSystemInfo() *SystemInfo {
	info := &SystemInfo{}

	// 获取CPU使用率
	if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
		info.CPUUsage = cpuPercent[0]
	}

	// 获取内存信息
	if vmstat, err := mem.VirtualMemory(); err == nil {
		info.MemoryUsage = vmstat.UsedPercent
		info.MemoryTotal = vmstat.Total / 1024 / 1024 // 转换为MB
		info.MemoryUsed = vmstat.Used / 1024 / 1024   // 转换为MB
	}

	// 获取磁盘信息（使用服务器程序所在目录）
	serverDir := filepath.Dir(config.ServerPath)
	if diskStat, err := disk.Usage(serverDir); err == nil {
		info.DiskUsage = diskStat.UsedPercent
		info.DiskTotal = diskStat.Total / 1024 / 1024 // 转换为MB
		info.DiskUsed = diskStat.Used / 1024 / 1024   // 转换为MB
	}

	// 获取系统运行时间
	if hostInfo, err := host.Info(); err == nil {
		info.Uptime = hostInfo.Uptime
	}

	return info
}

// AddOutputLine 添加输出行到缓冲区
func (sm *SystemMonitor) AddOutputLine(line string) {
	sm.outputBuffer = append(sm.outputBuffer, line)
	// 保持缓冲区大小
	if len(sm.outputBuffer) > 1000 {
		sm.outputBuffer = sm.outputBuffer[1:]
	}
}

// SetStartTimeNow 设置启动时间为当前时间
func (sm *SystemMonitor) SetStartTimeNow() {
	sm.startTime = time.Now().Format("2006-01-02 15:04:05")
}

// GetServerInfo 从进程输出中解析服务器信息
func (sm *SystemMonitor) GetServerInfo() *ServerInfo {
	info := &ServerInfo{
		Version:       "未知",
		LoaderVersion: "",
		StartTime:     sm.startTime,
		PlayerCount:   GetPlayerManager().GetPlayerCount(),
		MaxPlayers:    20, // 默认值
		Uptime:        0,
	}

	// 读取max-players
	fileManager := GetFileManager()
	if maxPlayers, err := fileManager.GetMaxPlayers(); err == nil && maxPlayers > 0 {
		info.MaxPlayers = maxPlayers
	}

	// 优先从配置文件读取版本
	if config.Version != "" && config.Version != "未知" {
		info.Version = config.Version
		info.LoaderVersion = config.LoaderVersion
	}

	for i := len(sm.outputBuffer) - 1; i >= 0 && i >= len(sm.outputBuffer)-100; i-- {
		line := sm.outputBuffer[i]
		// 解析版本信息
		if strings.Contains(line, "Version:") {
			ver := ""
			loader := ""
			if idx := strings.Index(line, "Version:"); idx != -1 {
				verStr := line[idx+8:]
				verStr = strings.TrimSpace(verStr)
				// 删除ANSI转义序列等多余文本
				verStr = removeANSI(verStr)
				// 解析Loader版本
				if strings.Contains(verStr, "with ") {
					parts := strings.SplitN(verStr, "with ", 2)
					ver = strings.TrimSpace(parts[0])
					loader = strings.TrimSpace(parts[1])
				} else {
					// 只取第一个空格前内容或括号前内容
					if sp := strings.IndexAny(verStr, " (【"); sp != -1 {
						ver = verStr[:sp]
					} else {
						ver = verStr
					}
				}
				if ver != "" {
					info.Version = ver
					updateConfigVersion(ver)
				}
				if loader != "" {
					info.LoaderVersion = loader
					updateConfigLoaderVersion(loader)
				}
			}
		}

		// 解析玩家数量
		if strings.Contains(line, "players online") {
			parts := strings.Fields(line)
			for j, part := range parts {
				if part == "players" && j > 0 {
					if count, err := strconv.Atoi(parts[j-1]); err == nil {
						info.PlayerCount = count
					}
				}
			}
		}

		// 解析启动时间
		if strings.Contains(line, "Started") && strings.Contains(line, "server") {
			// 尝试从输出中提取启动时间
			info.StartTime = time.Now().Format("2006-01-02 15:04:05")
		}
	}

	return info
}

// 更新config.json中的Version字段
func updateConfigVersion(version string) {
	configPath := ".\\Panel_Setting\\config.json"
	file, err := os.OpenFile(configPath, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	var cfg map[string]interface{}
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return
	}
	cfg["Version"] = version
	file.Seek(0, 0)
	file.Truncate(0)
	json.NewEncoder(file).Encode(cfg)
}

// 更新config.json中的LoaderVersion字段
func updateConfigLoaderVersion(loader string) {
	configPath := ".\\Panel_Setting\\config.json"
	file, err := os.OpenFile(configPath, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	var cfg map[string]interface{}
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return
	}
	cfg["LoaderVersion"] = loader
	file.Seek(0, 0)
	file.Truncate(0)
	json.NewEncoder(file).Encode(cfg)
}

// 删除ANSI转义序列等多余文本
func removeANSI(text string) string {
	// 删除ANSI转义序列
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	text = ansiRegex.ReplaceAllString(text, "")
	// 删除其他控制字符
	controlRegex := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	text = controlRegex.ReplaceAllString(text, "")
	return strings.TrimSpace(text)
}

// GetConsoleHistory 获取终端输出历史（最后100行，顺序为最早到最新）
func (sm *SystemMonitor) GetConsoleHistory() []string {
	history := []string{}
	start := 0
	if len(sm.outputBuffer) > 100 {
		start = len(sm.outputBuffer) - 100
	}
	for i := start; i < len(sm.outputBuffer); i++ {
		history = append(history, sm.outputBuffer[i])
	}
	return history
}

// 全局系统监控器实例
var globalSystemMonitor *SystemMonitor

// GetSystemMonitor 获取全局系统监控器实例
func GetSystemMonitor() *SystemMonitor {
	if globalSystemMonitor == nil {
		globalSystemMonitor = NewSystemMonitor()
	}
	return globalSystemMonitor
}
