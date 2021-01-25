# go-excel-tool
用golang实现的excel转json或者lua的配置表

字段类型：int/string/auto/table

其中table是需要组合字段，比如table+string 就是字符串数组

支持指定字段为主key：int+key或者string+key

允许为空，字段会不导出，注意强格式语言解析json可能需要额外处理
