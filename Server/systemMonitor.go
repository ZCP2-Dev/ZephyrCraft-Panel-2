package main

import (
	"path/filepath"
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

// GetServerInfo 从进程输出中解析服务器信息
func (sm *SystemMonitor) GetServerInfo() *ServerInfo {
	info := &ServerInfo{
		Version:     "未知",
		StartTime:   "",
		PlayerCount: 0,
		MaxPlayers:  20, // 默认值
		TPS:         20.0,
		Uptime:      0,
	}

	// 从最近的输出中解析信息
	for i := len(sm.outputBuffer) - 1; i >= 0 && i >= len(sm.outputBuffer)-100; i-- {
		line := sm.outputBuffer[i]

		// 解析版本信息
		if strings.Contains(line, "Version") && strings.Contains(line, "Bedrock") {
			if versionIndex := strings.Index(line, "Version"); versionIndex != -1 {
				parts := strings.Fields(line[versionIndex:])
				if len(parts) >= 2 {
					info.Version = parts[1]
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

		// 解析最大玩家数
		if strings.Contains(line, "max") && strings.Contains(line, "players") {
			parts := strings.Fields(line)
			for j, part := range parts {
				if part == "max" && j+1 < len(parts) {
					if max, err := strconv.Atoi(parts[j+1]); err == nil {
						info.MaxPlayers = max
					}
				}
			}
		}

		// 解析TPS
		if strings.Contains(line, "TPS") {
			parts := strings.Fields(line)
			for j, part := range parts {
				if part == "TPS" && j+1 < len(parts) {
					if tps, err := strconv.ParseFloat(parts[j+1], 64); err == nil {
						info.TPS = tps
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

// 全局系统监控器实例
var globalSystemMonitor *SystemMonitor

// GetSystemMonitor 获取全局系统监控器实例
func GetSystemMonitor() *SystemMonitor {
	if globalSystemMonitor == nil {
		globalSystemMonitor = NewSystemMonitor()
	}
	return globalSystemMonitor
}
