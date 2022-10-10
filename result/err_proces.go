package result

import (
	"errors"
	"github.com/asaswangc/gowork/validators"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// GinBindErr 客制化ErrorProcess
func GinBindErr(err error, _ ...interface{}) {
	if err != nil {
		if err.Error() == "EOF" {
			panic(NewConstErr("参数body不可为空", ParamParseErrCode, ""))
		}
		if ok, errs := validators.InterceptError(err); ok {
			panic(NewConstErr(errs, ParamCheckErrCode, ""))
		}
		panic(err)
	}
}

// SqlCrudErr mysql增删改查Err处理
func SqlCrudErr(err error, _ ...interface{}) {
	if err != nil {
		if value, ok := err.(*mysql.MySQLError); ok {
			switch value.Number {
			case 1062:
				panic(NewConstErr("数据已经存在", SelectExistErrCode, value.Message))
			case 1066:
				panic(NewConstErr("数据查询为空", SelectEmptyErrCode, value.Message))
			}
		}
		panic(err)
	}
}

// SqlCrudErrIgnoreRecordNotFound mysql增删改查Err处理
func SqlCrudErrIgnoreRecordNotFound(err error, _ ...interface{}) {
	if err != nil {
		if value, ok := err.(*mysql.MySQLError); ok {
			switch value.Number {
			case 1062:
				panic(NewConstErr("数据已经存在", SelectExistErrCode, value.Message))
			case 1066:
				panic(NewConstErr("数据查询为空", SelectEmptyErrCode, value.Message))
			}
		}

		// 忽略这个Err
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		panic(err)
	}
}
