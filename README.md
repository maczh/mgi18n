# mgi18n多语言框架

## 安装

```bash
go get -u github.com/maczh/mgi18n
```

## 使用说明

- 本框架依赖x-lang多语言字典库微服务模块提供数据支持

### 在配置文件中添加

```yaml
go:
  xlang:
    appName: test
    default: zh-cn
```

###  在main.go中，在mgconfig.Init()之后,添加一行

```go
    mgi18n.Init()
```

### 在Router.go中添加中间件引用

```go
	//添加多语种支持
	engine.Use(mgerr.RequestLanguage())
```

### 在代码中错误信息多语言处理

```go
	engine.Any("/err", func(c *gin.Context) {
		m := utils.GinParamMap(c)
		result = mgi18n.Error(-1,m["msgId"])
		c.JSON(http.StatusOK, result)
	})
```

### 在代码中返回带多语言的成功信息

```go
	engine.Any("/test", func(c *gin.Context) {
		m := utils.GinParamMap(c)
		result = mgi18n.Success(mgi18n.String(m["msgId"]))
		c.JSON(http.StatusOK, result)
	})
```

### 在代码中替换语言代码

```go
    msgId := "100001"
    msg := mgi18n.String(msgId)
```

### 其他函数
```go
    func ErrorWithMsg(code int, messageId, msg string) mgresult.Result
    func SuccessWithPage(data interface{},count, index, size, total int) mgresult.Result
```
