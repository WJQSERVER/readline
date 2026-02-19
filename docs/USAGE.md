# PureGo Readline 使用指南

`purego-readline` 是一个纯 Go 编写的跨平台行读取库，支持 Linux, macOS 和 Windows。

## 安装

由于这是一个本地项目，你可以直接在你的 Go 项目中通过 go mod 引用它。

```bash
go get github.com/your-username/purego-readline
```

## 快速开始

以下是一个简单的使用示例：

```go
package main

import (
	"fmt"
	"io"
	"purego-readline"
)

func main() {
	// 1. 初始化配置
	cfg := &readline.Config{
		Prompt: "> ",
	}

	// 2. 创建 Readline 实例
	rl, err := readline.NewInstance(cfg)
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// 3. 进入读取循环
	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt { // 处理 Ctrl+C
				fmt.Println("^C")
				continue
			} else if err == io.EOF { // 处理 Ctrl+D
				break
			}
			fmt.Println("Error:", err)
			break
		}

		fmt.Printf("输入内容: %s\n", line)
		if line == "exit" {
			break
		}
	}
}
```

## 功能特性

### 1. 自动补全 (Auto-completion)

你可以通过实现 `Completer` 接口来自定义补全逻辑：

```go
type Completer interface {
	Do(line []rune, pos int) (candidates [][]rune, length int)
}
```

库内置了一个简单的 `PrefixCompleter`：

```go
completer := &readline.PrefixCompleter{
    Candidates: []string{"help", "exit", "status", "version"},
}
cfg := &readline.Config{
    Prompt:    "> ",
    Completer: completer,
}
```

### 2. 历史记录 (History)

库默认在内存中管理历史记录。你可以通过上下方向键导航。

### 3. 多平台支持

- **Unix**: 完美支持所有主流终端。
- **Windows**: 自动尝试开启虚拟终端处理（VT Processing），支持 ANSI 转义序列。

## 常用快捷键

| 快捷键 | 功能 |
| --- | --- |
| `Enter` | 确认输入并返回 |
| `Ctrl+C` | 中断当前输入 (`ErrInterrupt`) |
| `Ctrl+D` | 在行首时发送 EOF，或者删除光标后的字符 |
| `Ctrl+L` | 清理当前行内容（逻辑上） |
| `方向键上/下` | 浏览历史记录 |
| `方向键左/右` | 移动光标 |
| `Home / End` | 移动光标到行首或行尾 |
| `Tab` | 触发自动补全 |
| `Backspace` | 删除光标前字符 |
| `Delete` | 删除光标后字符 |
