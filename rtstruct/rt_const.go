package rtstruct

const (
	Type          = "type"           // 类型
	Null          = "null"           // 可为空
	Column        = "column"         // 名称
	Comment       = "comment"        // 备注
	Default       = "default"        // 默认值
	NotNull       = "not null"       // 不可为空
	UniqueKey     = "unique key"     // 唯一键
	PrimaryKey    = "primary_key"    // 主键
	AutoIncrement = "auto_increment" // 自动递增，适用于整数类型
)

type SqlType string

// 数值类型
const (
	TINYINT   = SqlType("TINYINT")   // 小整数值
	SMALLINT  = SqlType("SMALLINT")  // 大整数值
	MEDIUMINT = SqlType("MEDIUMINT") // 大整数值
	INTEGER   = SqlType("INTEGER")   // 大整数值
	BIGINT    = SqlType("BIGINT")    // 极大整数值
	FLOAT     = SqlType("FLOAT")     // 单精度浮点数值
	DOUBLE    = SqlType("DOUBLE")    // 双精度浮点数值
	DECIMAL   = SqlType("DECIMAL")   // 小数值
)

// 日期和时间类型
const (
	DATE      = SqlType("DATE")      // 日期值
	TIME      = SqlType("TIME")      // 时间值或持续时间
	YEAR      = SqlType("YEAR")      // 年份值
	DATETIME  = SqlType("DATETIME")  // 混合日期和时间值
	TIMESTAMP = SqlType("TIMESTAMP") // 混合日期和时间值，时间戳
)

// 字符串类型SqlType(
const (
	CHAR       = SqlType("CHAR")       // 定长字符串
	VARCHAR    = SqlType("VARCHAR")    // 变长字符串
	TINYBLOB   = SqlType("TINYBLOB")   // 不超过 255 个字符的二进制字符串
	TINYTEXT   = SqlType("TINYTEXT")   // 短文本字符串
	BLOB       = SqlType("BLOB")       // 二进制形式的长文本数据
	TEXT       = SqlType("TEXT")       // 长文本数据
	MEDIUMBLOB = SqlType("MEDIUMBLOB") // 二进制形式的中等长度文本数据
	MEDIUMTEXT = SqlType("MEDIUMTEXT") // 中等长度文本数据
	LONGBLOB   = SqlType("LONGBLOB")   // 二进制形式的极大文本数据
	LONGTEXT   = SqlType("LONGTEXT")   // 极大文本数据
)

// SqlMappingGo 数据库类型与Go类型映射
var SqlMappingGo = map[SqlType]interface{}{

	// 日期和时间类型
	DATE:      "0",
	TIME:      "0",
	YEAR:      "0",
	DATETIME:  "0",
	TIMESTAMP: "0",

	// 字符串类型
	CHAR:       "0",
	VARCHAR:    "0",
	TINYBLOB:   "0",
	TINYTEXT:   "0",
	BLOB:       "0",
	TEXT:       "0",
	MEDIUMBLOB: "0",
	MEDIUMTEXT: "0",
	LONGBLOB:   "0",
	LONGTEXT:   "0",

	// 数值类型
	TINYINT:   int64(0),
	SMALLINT:  int64(0),
	MEDIUMINT: int64(0),
	INTEGER:   int64(0),
	BIGINT:    int64(0),
	FLOAT:     int64(0),
	DOUBLE:    int64(0),
	DECIMAL:   int64(0),
}
