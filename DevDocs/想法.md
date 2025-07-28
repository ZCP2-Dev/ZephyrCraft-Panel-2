<!--
这里是开发文档，用于存放开发者给其他开发者写的注意事项和当作备忘录使用。
-->

> 使用协议文档定义一个协议版本，同一协议版本下前后端随便连接，握手时也要检查协议版本

- 因为可能会增加新的东西或重构协议，同时拥有一份协议文档也方便其他开发者开发新的端，而不仅仅局限在MC

by n15421

根据代码库中的信息，WebSocket 通信时发送和接收的消息大多以 JSON 格式传递。以下是不同接口触发时发送到前端的 JSON 文本示例，以及对应的接口说明：

此处内容不一定准确

### 1. 服务器控制接口

#### 1.1 启动服务器
- **命令**：`start`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "start",
    "status": "success",
    "message": "服务器已启动"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "start",
    "status": "failed",
    "error": "启动服务器时发生错误"
}
```

#### 1.2 停止服务器
- **命令**：`stop`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "stop",
    "status": "success",
    "message": "服务器已停止"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "stop",
    "status": "failed",
    "error": "停止服务器时发生错误"
}
```

#### 1.3 重启服务器
- **命令**：`restart`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "restart",
    "status": "success",
    "message": "服务器已重启"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "restart",
    "status": "failed",
    "error": "重启服务器时发生错误"
}
```

#### 1.4 查询服务器状态
- **命令**：`status`
- **发送到前端的 JSON 示例（运行中）**：
```json
{
    "command": "status",
    "status": "running",
    "message": "服务器正在运行"
}
```
- **发送到前端的 JSON 示例（已停止）**：
```json
{
    "command": "status",
    "status": "stopped",
    "message": "服务器已停止"
}
```

### 2. 文件操作接口

#### 2.1 创建压缩文件
- **命令**：`createZip`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "createZip",
    "status": "success",
    "message": "压缩文件创建成功",
    "filePath": "/path/to/archive.zip"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "createZip",
    "status": "failed",
    "error": "创建压缩文件时发生错误"
}
```

#### 2.2 获取文件列表
- **命令**：`getFiles`
- **发送到前端的 JSON 示例**：
```json
{
    "command": "getFiles",
    "status": "success",
    "fileList": [
        {
            "name": "file1.txt",
            "path": "/path/to/file1.txt",
            "size": 1024,
            "modified": "2024-01-01T12:00:00Z"
        },
        {
            "name": "file2.txt",
            "path": "/path/to/file2.txt",
            "size": 2048,
            "modified": "2024-01-02T12:00:00Z"
        }
    ]
}
```

#### 2.3 读取文件内容
- **命令**：`readFile`
- **发送到前端的 JSON 示例**：
```json
{
    "command": "readFile",
    "status": "success",
    "filePath": "/path/to/file.txt",
    "fileContent": "这是文件的内容。"
}
```

#### 2.4 写入文件
- **命令**：`writeFile`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "writeFile",
    "status": "success",
    "message": "文件写入成功",
    "filePath": "/path/to/file.txt"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "writeFile",
    "status": "failed",
    "error": "写入文件时发生错误",
    "filePath": "/path/to/file.txt"
}
```

### 3. 系统信息接口

#### 3.1 获取系统信息
- **命令**：`getSystemInfo`
- **发送到前端的 JSON 示例**：
```json
{
    "command": "getSystemInfo",
    "status": "success",
    "systemInfo": {
        "cpuUsage": 20,
        "memoryUsage": 50,
        "diskUsage": 30
    }
}
```

#### 3.2 获取服务器信息
- **命令**：`getServerInfo`
- **发送到前端的 JSON 示例**：
```json
{
    "command": "getServerInfo",
    "status": "success",
    "serverInfo": {
        "version": "1.0.0",
        "startTime": "2024-01-01T12:00:00Z",
        "playerCount": 10
    }
}
```

### 4. 终端命令接口

#### 4.1 发送终端命令
- **命令**：`input`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "input",
    "status": "success",
    "output": "命令执行结果"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "input",
    "status": "failed",
    "error": "执行命令时发生错误"
}
```

### 5. 玩家操作接口

#### 5.1 踢出玩家
- **命令**：`input`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "input",
    "content": "kick player1",
    "status": "success",
    "message": "玩家 player1 已被踢出"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "input",
    "content": "kick player1",
    "status": "failed",
    "error": "踢出玩家时发生错误"
}
```

#### 5.2 发送消息给玩家
- **命令**：`input`
- **发送到前端的 JSON 示例（成功响应）**：
```json
{
    "command": "input",
    "content": "tell player1 Hello!",
    "status": "success",
    "message": "消息已发送给玩家 player1"
}
```
- **发送到前端的 JSON 示例（失败响应）**：
```json
{
    "command": "input",
    "content": "tell player1 Hello!",
    "status": "failed",
    "error": "发送消息时发生错误"
}
```

>
