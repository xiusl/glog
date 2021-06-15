# logging

日志库，

## 主要功能
- 支持向不同位置输出日志（os.Stdout，file，db等）
- 日志级别控制
  - Debug/Trace/Info/Warning/Error/Fatal
  - 可控制不同级别的输出
- 日志格式
  - 时间、行号、文件名、函数名、日志级别、日志信息
- 日志切割（文件）
- 异步写日志（文件）

## 使用
- 控制台便捷使用
```
logging.Debug(...)
```
- 文件使用
```
```