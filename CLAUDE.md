# LarkBase SDK - 项目概览

## 项目简介

LarkBase 是一个 Go SDK，用于操作飞书（Lark）多维表格（Base）。提供类型安全的 ORM 风格接口，支持对 Lark Base 表进行增删改查、批量操作、字段过滤、数据库同步等功能。

- **Module**: `github.com/lincaiyong/larkbase`
- **Go 版本**: 1.23+
- **主要依赖**: `gorm.io/gorm`, `github.com/mitchellh/mapstructure`, `github.com/lincaiyong/log`

---

## 目录结构

```
larkbase/
├── larksuite.go       # 入口：Connect/ConnectUrl/ConnectAny/DescribeTable
├── conn.go            # Connection[T] 类型及全部 CRUD 方法
├── record.go          # Record / AnyRecord 数据结构
├── reflect.go         # 反射工具：struct 解析与 Record 互转
├── field.go           # larkfield.Field 的类型别名（TextField 等）
├── findoption.go      # FindOption 查询构建器（过滤/排序/分页）
├── ext.go             # CreateTable / CreateAll（泛型批量插入）
├── sync.go            # SyncToDatabase（同步到 MySQL）
├── meta.go            # Meta 结构体（RecordId / ModifiedTime）
├── larkfield/         # 字段类型系统
│   ├── field.go       # Field 接口定义
│   ├── type.go        # 字段类型枚举与工厂
│   ├── field_base.go  # BaseField（所有字段的公共实现）
│   ├── condition.go   # 条件构建器（Is/Contains/IsGreater 等）
│   ├── datetime.go    # 时间工具（北京时区 Asia/Shanghai）
│   ├── parse.go       # API 响应 → Go 值（解析）
│   ├── build.go       # Go 值 → API 请求格式（序列化）
│   └── field_*.go     # 各字段类型实现
├── larksuite/         # 飞书官方 SDK 封装
│   ├── client.go      # SDK 客户端初始化
│   ├── bitable/       # Bitable API（自动生成代码）
│   └── core/          # HTTP/Token 管理/日志
└── example/
    └── main.go        # 使用示例
```

---

## 核心概念

### 1. 认证配置

```go
// 通过环境变量（推荐）
os.Setenv("LARK_APP_ID", "xxx")
os.Setenv("LARK_APP_SECRET", "xxx")

// 或代码设置
larkbase.SetAppIdSecret("app_id", "app_secret")
```

### 2. 表结构定义

用 Go struct 描述 Lark Base 表，**规则严格**：

```go
type MyRecord struct {
    // 第一个字段必须是 larkbase.Meta，lark tag 为完整的表 URL
    larkbase.Meta `lark:"https://bytedance.larkoffice.com/base/{appToken}?table={tableId}"`

    // 其余字段必须是 larkfield.*Field 类型，lark tag 为字段名
    Name   larkbase.TextField         `lark:"name"`
    Age    larkbase.NumberField       `lark:"age"`
    Date   larkbase.DateField         `lark:"日期"`
    Multi  larkbase.MultiSelectField  `lark:"multi"`
    Single larkbase.SingleSelectField `lark:"单选"`
    Check  larkbase.CheckboxField     `lark:"check"`
    Link   larkbase.UrlField          `lark:"超链接"`
    No     larkbase.AutoNumberField   `lark:"no"`      // 只读
    Lookup larkbase.LookupField       `lark:"lookup"`  // 只读
    Mtime  larkbase.ModifiedTimeField `lark:"modified_time"` // 只读
}
```

### 3. 字段类型映射

| Go 类型 | Lark 字段类型 |
|---------|-------------|
| `TextField` | 文本、邮箱(EmailField)、条码(BarcodeField) |
| `NumberField` | 数字、货币(CurrencyField)、进度(ProgressField)、评分(RatingField) |
| `SingleSelectField` | 单选 |
| `MultiSelectField` | 多选 |
| `DateField` | 日期 |
| `CheckboxField` | 复选框 |
| `UrlField` | 超链接 |
| `AutoNumberField` | 自动编号（只读） |
| `ModifiedTimeField` | 修改时间（只读） |
| `LookupField` | 引用查找（只读） |
| `FormulaField` | 公式（只读） |

---

## API 参考

### 连接初始化

```go
// 从 struct 的 Meta lark tag 读取 URL
conn, err := larkbase.Connect[MyRecord](ctx)

// 指定 URL
conn, err := larkbase.ConnectUrl[MyRecord](ctx, tableUrl)

// 动态 schema（不需要预定义 struct）
conn, err := larkbase.ConnectAny(ctx, tableUrl)

// 自动生成 struct 定义字符串
structCode, err := larkbase.DescribeTable(ctx, tableUrl)
```

### CRUD 操作

```go
// 查找单条
var r MyRecord
err = conn.Find(&r, larkbase.NewFindOption(conn.FilterAnd(conn.Condition().Name.Is("andy"))))

// 查找多条（支持分页、过滤、排序、限制）
var records []*MyRecord
err = conn.FindAll(&records, larkbase.NewFindOption(conn.FilterOr(cond1, cond2)).Limit(100))

// 统计记录数（仅发一次请求，pageSize=1，读取响应中的 Total 字段）
count, err := conn.Count(nil)                                            // 全表统计
count, err := conn.Count(larkbase.NewFindOption(conn.FilterAnd(...)))   // 带过滤条件统计

// 创建
var r MyRecord
r.Name.SetValue("test")
err = conn.Create(&r)  // r.RecordId 会被填充

// 批量创建
records, err = conn.CreateAll(records)

// 更新（仅更新 dirty 字段）
r.Age.SetIntValue(25)
err = conn.Update(&r)

// 批量更新
err = conn.UpdateAll(records)

// 删除
err = conn.Delete(&r)
err = conn.DeleteAll(records)
```

### 查询构建

```go
opt := larkbase.NewFindOption(
    conn.FilterAnd(
        conn.Condition().Name.Is("andy"),
        conn.Condition().Age.IsGreater(18),
    ),
).Limit(50)

// 排序
opt.Sort(conn.Condition().Age, true) // true=升序

// 视图过滤（创建或使用视图）
err = conn.CreateView("myView", conn.ViewFilterAnd(conn.Condition().Name.IsNotEmpty()))
```

### 条件操作符

| 方法 | 适用类型 |
|-----|---------|
| `Is(v)` / `IsNot(v)` | 所有类型 |
| `Contains(v)` / `DoesNotContain(v)` | Text, MultiSelect |
| `IsEmpty()` / `IsNotEmpty()` | 所有类型 |
| `IsGreater(v)` / `IsGreaterEqual(v)` | Number, Date |
| `IsLess(v)` / `IsLessEqual(v)` | Number, Date |
| `IsToday()` / `IsYesterday()` / `IsTomorrow()` | Date |
| `IsCurrentWeek()` / `IsLastWeek()` | Date |
| `IsCurrentMonth()` / `IsLastMonth()` | Date |

### 序列化

```go
jsonStr, err := conn.MarshalRecord(&r)
jsonStr, err := conn.MarshalRecords(records)
jsonStr = conn.MarshalIgnoreError(records) // 忽略错误
```

### 数据库同步

```go
// 需要表中有 ModifiedTimeField 字段
db, _ := gorm.Open(mysql.Open(dsn))
err = conn.SyncToDatabase(db, batchSize)
```

### 批量操作配置

```go
conn.SetBatchSize(100)  // 设置每批次处理数量，默认不分批
```

---

## 架构要点

### 层次结构

```
用户代码（struct 定义）
    ↓
larkbase（Connection[T]，CRUD）
    ↓
larkfield（Field 类型，Parse/Build，条件）
    ↓
larksuite SDK（Bitable API，HTTP/Token）
    ↓
飞书 OpenAPI（open.feishu.cn）
```

### 泛型设计

`Connection[T any]` 是核心泛型类型，T 须满足：
- 第一个字段为 `larkbase.Meta`
- 其余字段为 `larkfield.*Field` 类型

### 反射使用策略

反射**仅在连接初始化时**使用（解析 struct 结构），CRUD 操作期间不使用反射，以保证性能。

### Dirty 字段追踪

所有字段有 dirty 标记。`Update` 和 `UpdateAll` 只发送被 `SetValue()` 修改过的字段，避免多余 API 调用。

### 时区

所有日期时间处理使用**北京时间（Asia/Shanghai）**。

---

## 已知问题

- `SyncFromApi` 方法基本为桩代码，未完整实现
- `SyncToDatabase` 使用 MySQL 特定语法（`ON DUPLICATE KEY UPDATE`），不支持其他数据库

---

## 开发注意事项

1. **新增字段类型**：需在 `larkfield/type.go`、`parse.go`、`build.go` 同步添加逻辑
2. **只读字段**（AutoNumber/ModifiedTime/Lookup/Formula）的 `Build()` 返回 nil，不会被发送到 API
3. `TextField` 超过 80,000 字符会被截断并打印警告
4. URL 格式：`https://*.larkoffice.com/base/{appToken}?table={tableId}[&view={viewId}]`
5. 批量操作默认每页 100 条（由 Lark API 限制）
