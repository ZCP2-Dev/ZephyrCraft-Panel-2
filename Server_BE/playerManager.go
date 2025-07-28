package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Player 玩家信息结构
type Player struct {
	Name string `json:"name"` // 玩家名称
	XUID string `json:"xuid"` // 玩家XUID（可能为空）
}

// PlayerManager 玩家管理器
type PlayerManager struct {
	players map[string]*Player // 在线玩家列表，key为玩家名称
	mu      sync.RWMutex       // 读写锁，保护玩家列表的并发访问
}

// NewPlayerManager 创建新的玩家管理器
func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[string]*Player),
	}
}

// GetPlayers 获取所有在线玩家
func (pm *PlayerManager) GetPlayers() []*Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	players := make([]*Player, 0, len(pm.players))
	for _, player := range pm.players {
		players = append(players, player)
	}
	return players
}

// GetPlayerCount 获取在线玩家数量
func (pm *PlayerManager) GetPlayerCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.players)
}

// AddPlayer 添加玩家到在线列表
func (pm *PlayerManager) AddPlayer(name, xuid string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.players[name] = &Player{
		Name: name,
		XUID: xuid,
	}

	log.Printf("[PlayerManager] 玩家加入: %s (XUID: %s), 当前在线: %d", name, xuid, len(pm.players))
}

// RemovePlayer 从在线列表中移除玩家
func (pm *PlayerManager) RemovePlayer(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.players[name]; exists {
		delete(pm.players, name)
		log.Printf("[PlayerManager] 玩家离开: %s, 当前在线: %d", name, len(pm.players))
	}
}

// ClearPlayers 清空玩家列表（服务器重启时使用）
func (pm *PlayerManager) ClearPlayers() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.players = make(map[string]*Player)
	log.Printf("[PlayerManager] 玩家列表已清空")
}

// ParsePlayerEvent 解析玩家事件（连接/断开）
func (pm *PlayerManager) ParsePlayerEvent(line string) bool {
	// 玩家连接事件的正则表达式
	// 匹配格式: "12:19:34.368 INFO [Server] Player connected: win81pro, xuid: 2535421504983964"
	connectedRegex := regexp.MustCompile(`Player connected:\s*([^,]+)(?:,\s*xuid:\s*([^\s]+))?`)
	// 玩家生成事件的正则表达式
	// 匹配格式: "14:19:11.664 INFO [Server] Player Spawned: LiteZero2614929 xuid: 2535451043182905, pfid: bc8f95ca6b684025"
	connectedRegex4 := regexp.MustCompile(`Player Spawned:\s*([^\s]+)\s+xuid:\s*([^\s,]+)`)

	// 玩家断开事件的正则表达式
	// 匹配格式: "12:20:27.624 INFO [Server] Player disconnected: win81pro, xuid: 2535421504983964, pfid: c1e893d6b8ec3e71"
	disconnectedRegex := regexp.MustCompile(`Player disconnected:\s*([^,]+)(?:,\s*xuid:\s*([^\s]+))?`)

	// 另一种断开事件格式
	// 匹配格式: "12:20:27.624 INFO [Server] Player win81pro left the game"
	disconnectedRegex2 := regexp.MustCompile(`Player\s+([^\s]+)\s+left\s+the\s+game`)

	// 第三种断开事件格式
	// 匹配格式: "12:20:27.624 INFO [Server] win81pro left the game"
	disconnectedRegex3 := regexp.MustCompile(`([^\s]+)\s+left\s+the\s+game`)

	// kick命令成功执行的正则表达式
	// 匹配格式: "12:20:27.624 INFO [Server] Kicked player: win81pro"
	kickSuccessRegex := regexp.MustCompile(`Kicked player:\s*([^\s]+)`)

	// 另一种kick命令输出格式
	// 匹配格式: "12:20:27.624 INFO [Server] Player win81pro was kicked"
	kickSuccessRegex2 := regexp.MustCompile(`Player\s+([^\s]+)\s+was\s+kicked`)

	// 检查是否为玩家连接事件
	if matches := connectedRegex.FindStringSubmatch(line); matches != nil {
		playerName := matches[1]
		xuid := ""
		if len(matches) > 2 {
			xuid = matches[2]
		}
		log.Printf("[PlayerManager][DEBUG]检测到玩家连接事件: %s (XUID: %s)", playerName, xuid)

		return true
	}

	// 检查是否为玩家生成事件
	if matches := connectedRegex4.FindStringSubmatch(line); matches != nil {
		playerName := matches[1]
		xuid := ""
		if len(matches) > 2 {
			xuid = matches[2]
		}

		if isDebug {
			log.Printf("[PlayerManager][DEBUG]检测到玩家生成事件: %s (XUID: %s)", playerName, xuid)
		}

		if xuid != "" {
			if isDebug {
				log.Printf("[调试] 开始处理玩家 XUID: %s", xuid)
			}

			go func(name, xuid string) {
				if isDebug {
					log.Printf("[调试] 启动协程处理玩家 %s (XUID: %s) 的云黑查询", name, xuid)
				}

				client := &http.Client{Timeout: 5 * time.Second}
				url := "http://uniteban.xyz:19132/api.php?xuid=" + xuid

				if isDebug {
					log.Printf("[调试] 准备请求云黑API: %s", url)
				}

				resp, err := client.Get(url)
				if err != nil {
					log.Printf("[Uniteban] 请求失败: %v", err)
					if isDebug {
						log.Printf("[调试] API请求失败，错误: %v，正常添加玩家", err)
					}
					pm.AddPlayer(name, xuid) // API失败，正常添加
					return
				}
				defer resp.Body.Close()

				if isDebug {
					log.Printf("[调试] 云黑API响应状态码: %d", resp.StatusCode)
				}

				var result struct {
					Exists bool   `json:"exists"`
					Reason string `json:"reason"`
				}

				if isDebug {
					log.Printf("[调试] 开始解析API响应JSON")
				}

				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					log.Printf("[Uniteban] 解析响应失败: %v", err)
					if isDebug {
						log.Printf("[调试] JSON解析失败，错误: %v，正常添加玩家", err)
					}
					pm.AddPlayer(name, xuid) // 解析失败，正常添加
					return
				}

				if isDebug {
					log.Printf("[调试] 云黑API查询结果: Exists=%v, Reason=%s", result.Exists, result.Reason)
				}

				if result.Exists {
					reason := result.Reason
					if reason == "" {
						reason = "云黑封禁"
					}

					if isDebug {
						log.Printf("[调试] 玩家 %s (XUID: %s) 在云黑名单中，理由: %s", name, xuid, reason)
					}

					if config.Uniteban {
						log.Printf("[Uniteban] 玩家 %s (XUID: %s) 在云黑名单中，踢出，理由: %s", name, xuid, reason)
						pm.kickPlayer(name, "\n你因为在uniteban云黑系统中因理由: "+reason+" 被踢出游戏。")
						getProcessManager().sendMessage(Message{Content: fmt.Sprintf("[Uniteban] 玩家 %s (XUID: %s) 在云黑名单中，但因配置原因未自动踢出！", name, xuid)})
						if isDebug {
							log.Printf("[调试] 已执行踢人操作并发送广播")
						}
					} else {
						log.Printf("[Uniteban][警告] 玩家 %s (XUID: %s) 在云黑名单中，但因配置原因未自动踢出！", name, xuid)
						pm.AddPlayer(name, xuid+"_危险玩家(云黑)")
						getProcessManager().sendMessage(Message{Content: fmt.Sprintf("[Uniteban] 玩家 %s (XUID: %s) 在云黑名单中，但因配置原因未自动踢出！", name, xuid)})

						if isDebug {
							log.Printf("[调试] 未执行踢人操作，已添加标记并发送警告广播")
						}
					}
				} else {
					pm.AddPlayer(name, xuid)

					if isDebug {
						log.Printf("[调试] 玩家 %s (XUID: %s) 不在云黑名单中，正常添加", name, xuid)
					}
				}
			}(playerName, xuid)
		}
		return true
	}

	// 检查是否为玩家断开事件
	if matches := disconnectedRegex.FindStringSubmatch(line); matches != nil {
		playerName := matches[1]
		log.Printf("[PlayerManager][DEBUG]检测到玩家断开事件: %s", playerName)
		pm.RemovePlayer(playerName)
		return true
	}

	// 检查是否为第二种断开事件格式
	if matches := disconnectedRegex2.FindStringSubmatch(line); matches != nil {
		playerName := matches[1]
		log.Printf("[PlayerManager][DEBUG]检测到玩家断开事件(格式2): %s", playerName)
		pm.RemovePlayer(playerName)
		return true
	}

	// 检查是否为第三种断开事件格式
	if matches := disconnectedRegex3.FindStringSubmatch(line); matches != nil {
		playerName := matches[1]
		log.Printf("[PlayerManager][DEBUG]检测到玩家断开事件(格式3): %s", playerName)
		pm.RemovePlayer(playerName)
		return true
	}

	// 检查是否为kick命令成功执行
	if matches := kickSuccessRegex.FindStringSubmatch(line); matches != nil {
		playerName := matches[1]
		log.Printf("[PlayerManager][DEBUG]检测到kick命令成功: %s", playerName)
		pm.RemovePlayer(playerName)
		return true
	}

	// 检查是否为第二种kick命令输出格式
	if matches := kickSuccessRegex2.FindStringSubmatch(line); matches != nil {
		playerName := matches[1]
		log.Printf("[PlayerManager][DEBUG]检测到kick命令成功(格式2): %s", playerName)
		pm.RemovePlayer(playerName)
		return true
	}

	// 调试：记录所有包含"Player"、"Kicked"、"left"、"disconnected"的行，帮助识别其他格式
	if isDebug && (strings.Contains(line, "Player") || strings.Contains(line, "player") ||
		strings.Contains(line, "Kicked") || strings.Contains(line, "kicked") ||
		strings.Contains(line, "left") || strings.Contains(line, "disconnected")) {
		log.Printf("[PlayerManager][DEBUG]包含Player/Kicked/left/disconnected的行: %s", strings.TrimSpace(line))
	}

	// 额外检查：如果行包含"left"或"disconnected"但没有被上面的正则表达式匹配，记录详细信息
	if isDebug && (strings.Contains(line, "left") || strings.Contains(line, "disconnected")) {
		log.Printf("[PlayerManager][DEBUG]可能的玩家退出事件未被匹配: %s", strings.TrimSpace(line))
	}

	// 通用玩家退出检测：如果行包含"left"或"disconnected"，尝试提取玩家名称
	if strings.Contains(line, "left") || strings.Contains(line, "disconnected") {
		// 尝试从行中提取玩家名称
		words := strings.Fields(line)
		for i, word := range words {
			// 如果找到"left"或"disconnected"，检查前面的词是否是玩家名称
			if (word == "left" || word == "disconnected") && i > 0 {
				potentialPlayerName := words[i-1]
				// 检查这个名称是否在当前玩家列表中
				pm.mu.RLock()
				_, exists := pm.players[potentialPlayerName]
				pm.mu.RUnlock()

				if exists {
					log.Printf("[PlayerManager][DEBUG]通过通用检测发现玩家退出: %s", potentialPlayerName)
					pm.RemovePlayer(potentialPlayerName)
					return true
				}
			}
		}
	}

	return false
}

// 新增kickPlayer方法
func (pm *PlayerManager) kickPlayer(name, reason string) {
	// 通过ProcessManager发送kick命令
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	proc := getProcessManager()
	if proc != nil {
		cmd := "kick " + name + " " + reason
		_ = proc.SendCommand(cmd)
	}
}

// 全局玩家管理器实例
var globalPlayerManager *PlayerManager

// GetPlayerManager 获取全局玩家管理器实例
func GetPlayerManager() *PlayerManager {
	if globalPlayerManager == nil {
		globalPlayerManager = NewPlayerManager()
	}
	return globalPlayerManager
}
