# Readline (Pure Go) 详尽使用指南

`github.com/WJQSERVER/readline` 是一个纯 Go 编写的、无 cgo 依赖的多平台行编辑库。它提供了类似 GNU Readline 的交互体验，支持 Unicode 字符、历史记录导航、自动补全以及跨平台（Unix/Windows）兼容性。

> **⚠️ 注意**：这是一个**实验性原型项目**。API 可能会在不经通知的情况下发生更改，建议在生产环境中谨慎使用。

---

## 目录
1. [安装](#安装)
2. [快速开始](#快速开始)
3. [核心配置 (Config)](#核心配置-config)
4. [自动补全 (Auto-completion)](#自动补全-auto-completion)
5. [历史记录管理](#历史记录管理)
6. [终端支持与渲染](#终端支持与渲染)
7. [错误处理](#错误处理)
8. [快捷键列表](#快捷键列表)

---

## 安装

你可以直接通过 `go get` 安装：

```bash
go get github.com/WJQSERVER/readline
```

确保你的 Go 版本在 **1.26** 或以上，因为本项目利用了一些最新的语言特性。

---

## 快速开始

一个最小化的交互式 CLI 示例：

```go
package main

import (
    "fmt"
    "io"
    "github.com/WJQSERVER/readline"
)

func main() {
    rl, err := readline.NewInstance(&readline.Config{
        Prompt: "> ",
    })
    if err != nil {
        panic(err)
    }
    defer rl.Close()

    for {
        line, err := rl.Readline()
        if err != nil {
            if err == readline.ErrInterrupt {
                fmt.Println("^C")
                continue
            } else if err == io.EOF {
                fmt.Println("Goodbye!")
                break
            }
            fmt.Println("Error:", err)
            break
        }
        fmt.Printf("你输入了: %s\n", line)
    }
}
```

---

## 核心配置 (Config)

`Config` 结构体定义了 Readline 实例的行为：

```go
type Config struct {
    // Prompt 是每一行开始时的提示符
    Prompt string

    // History 接口实现，默认为内存历史记录
    History History

    // Completer 接口实现，用于处理 Tab 补全
    Completer Completer

    // Stdin 指定输入源，默认为 os.Stdin
    Stdin *os.File

    // Stdout 指定输出源，默认为 os.Stdout
    Stdout *os.File
}
```

---

## 自动补全 (Auto-completion)

### 使用内置的前缀补全器
库内置了 `PrefixCompleter`，适用于简单的静态关键字补全：

```go
completer := &readline.PrefixCompleter{
    Candidates: []string{"help", "exit", "list", "status", "你好"},
}
cfg := &readline.Config{
    Completer: completer,
}
```

### 实现自定义补全器
你可以通过实现 `Completer` 接口来实现更复杂的逻辑（如基于上下文的补全）：

```go
type MyCompleter struct{}

// Do 方法接收当前整行内容 (runes) 和光标位置 (pos)
// 返回匹配的建议列表 (candidates) 和建议替换的原始字符串长度 (length)
func (m *MyCompleter) Do(line []rune, pos int) (candidates [][]rune, length int) {
    // 逻辑示例：只对单词前缀进行补全
    // ... 解析逻辑 ...
    return [][]rune{[]rune("suggestion1"), []rune("suggestion2")}, 3
}
```

---

## 历史记录管理

库默认提供 `NewHistory()` 生成一个基于内存的历史记录管理器。
- **Append**: 在成功读取一行后，Readline 会自动将其加入历史记录。
- **导航**: 使用方向键 `↑` 和 `↓` 可以在历史记录中往返切换。
- **限制**: 当前原型版本仅支持内存存储，持久化到文件功能已列入 [Roadmap](../roadmap.md)。

---

## 终端支持与渲染

### Unicode 与宽字符支持
我们使用 `github.com/mattn/go-runewidth` 来计算字符的显示宽度。这意味着中文、日文、韩文以及 Emoji 在终端中都能正确计算光标位置，不会出现对齐偏差。

### Windows 适配
在 Windows 系统上，库会尝试调用 API 开启 `ENABLE_VIRTUAL_TERMINAL_PROCESSING`。
- 如果你的 Windows 版本较新（Windows 10 1607+），它将完美支持 ANSI 转义序列。
- 库在内部抽象了终端交互，确保在 Unix 和 Windows 下有一致的 Raw Mode 体验。

---

## 错误处理

`Readline()` 方法可能返回以下特殊错误：
- `readline.ErrInterrupt`: 用户按下了 `Ctrl+C`。
- `io.EOF`: 用户按下了 `Ctrl+D`（通常表示输入流结束）。
- 其他系统错误：如终端状态重置失败等。

建议始终通过 `if err == readline.ErrInterrupt` 等方式明确处理这些信号。

---

## 快捷键列表

| 快捷键 | 功能描述 |
| :--- | :--- |
| **Enter** | 提交当前行并换行 |
| **Ctrl + C** | 发送中断信号 (ErrInterrupt) |
| **Ctrl + D** | 在行首发送 EOF，或删除光标处字符 |
| **Ctrl + L** | 清除行显示（逻辑重绘） |
| **Tab** | 触发自动补全 |
| **↑ / ↓** | 切换上/下一条历史记录 |
| **← / →** | 光标左右移动 |
| **Home** / **Ctrl + A** | 移动到行首 |
| **End** / **Ctrl + E** | 移动到行尾 |
| **Backspace** | 删除光标前的一个字符 |
| **Delete** | 删除光标处的一个字符 |
