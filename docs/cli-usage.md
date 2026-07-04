# qlog-search CLI 使用说明

`qlog-search` 是一个腾讯云 CLS 日志查询 CLI 工具，支持**命令行参数模式**（默认）和**交互参数模式**两种使用方式，查询结果中匹配的关键字会以红色高亮显示。

---

## 构建

```bash
go build -o qlog-search cmd/cli/main.go
```

构建完成后会在当前目录生成 `qlog-search` 可执行文件。

---

## 环境变量配置

运行前需要配置 `.env` 文件（或通过系统环境变量设置）：

| 环境变量 | 必填 | 默认值 | 说明 |
|---------|------|--------|------|
| `QCLOUD_SECRET_ID` | **是** | - | 腾讯云 API SecretId |
| `QCLOUD_SECRET_KEY` | **是** | - | 腾讯云 API SecretKey |
| `QCLOUD_TOPIC_ID` | **是** | - | CLS 日志主题 ID |
| `QCLOUD_ENDPOINT` | 否 | `tencentcloudapi.com` | API 端点地址 |
| `QCLOUD_REGION` | 否 | `ap-beijing` | 地域 |
| `QUERY_LOG_LIMIT` | 否 | `500` | 单次查询返回日志条数上限 |
| `INTERACTIVE_PARAM_MODE` | 否 | `false` | 设为 `true` 启用交互参数模式 |

`.env` 文件示例：

```bash
QCLOUD_SECRET_ID=your_secret_id
QCLOUD_SECRET_KEY=your_secret_key
QCLOUD_TOPIC_ID=7597e6d2-5b0b-xxxx-xxxxx-fb25c8dffea7
QCLOUD_ENDPOINT=xx.tencentcloudapi.com
QUERY_LOG_LIMIT=1000
INTERACTIVE_PARAM_MODE=false
```

---

## 模式一：命令行参数模式（默认）

当 `INTERACTIVE_PARAM_MODE=false` 或未设置时，使用命令行 flag 传递查询参数。

### 参数说明

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `-query` | string | **是** | `""` | 日志搜索关键字 |
| `-period` | string | 否 | `last15m` | 日志查询时间范围 |
| `-topicId` | string | 否 | 环境变量 `QCLOUD_TOPIC_ID` | CLS Topic ID，可覆盖环境变量 |
| `-sort` | string | 否 | `asc` | 排序方式，`asc` 正序 / `desc` 倒序 |

### period 可选值

| 值 | 含义 |
|------|------|
| `last15m` | 最近 15 分钟 |
| `last1h` | 最近 1 小时 |
| `last6h` | 最近 6 小时 |
| `last1d` | 最近 1 天 |
| `last7d` | 最近 7 天 |

### 使用示例

```bash
# 基本查询 - 搜索最近15分钟包含 "error" 的日志
./qlog-search -query "error"

# 指定时间范围 - 搜索最近1小时
./qlog-search -query "timeout" -period last1h

# 指定排序方式 - 倒序查看最新日志
./qlog-search -query "panic" -period last6h -sort desc

# 覆盖 Topic ID
./qlog-search -query "connection refused" -topicId "your-custom-topic-id"

# 完整参数
./qlog-search -query "OOM" -period last1d -sort desc -topicId "your-topic-id"
```

### 输出示例

```
2025-04-15 09:30:01 service-A [ERROR] connection timeout after 30s
2025-04-15 09:30:05 service-B [ERROR] failed to process request: context deadline exceeded
```

> 其中 `ERROR` 等匹配关键字部分会以 **红色背景** 高亮显示（类似 `grep --color` 效果）。

---

## 模式二：交互参数模式

当 `.env` 中设置 `INTERACTIVE_PARAM_MODE=true` 时，工具将通过终端交互逐步引导输入查询条件。

### 交互流程

1. **输入查询关键字** — 手动输入日志搜索关键字（不能为空）
2. **选择时间范围** — 从预设列表中选择（last15m / last1h / last6h / last1d / last7d）
3. **输入排序方式** — 输入 `asc` 或 `desc`（可留空，默认 asc）

### 使用示例

```bash
# 直接运行，进入交互模式
./qlog-search
```

交互过程：

```
✔ Query keyword: error
Use the arrow keys to navigate: ↓ ↑ → ←
? Select period format:
  ▸ last15m
    last1h
    last6h
    last1d
    last7d
✔ Sort type (asc/desc), default asc: desc
```

---

## 模式对比

| 特性 | 命令行参数模式（默认） | 交互参数模式 |
|------|----------------------|-------------|
| 触发条件 | `INTERACTIVE_PARAM_MODE=false` 或未设置 | `INTERACTIVE_PARAM_MODE=true` |
| Query 来源 | `-query` flag | 交互输入 |
| Period 来源 | `-period` flag（默认 `last15m`） | 交互选择 |
| Sort 来源 | `-sort` flag（默认 `asc`） | 交互输入 |
| TopicId 来源 | `-topicId` flag（默认取环境变量） | 环境变量 `QCLOUD_TOPIC_ID` |
| 适用场景 | 脚本调用 / 自动化 / CI / SKILL 集成 | 手动查询 / 调试 |

---

## 快速开始

```bash
# 1. 克隆项目并进入目录
# 2. 配置 .env 文件
cp .env.example .env
# 编辑 .env 填入腾讯云凭证和 Topic ID

# 3. 构建
go build -o qlog-search cmd/cli/main.go

# 4. 运行（命令行参数模式）
./qlog-search -query "error" -period last1h -sort desc
```
