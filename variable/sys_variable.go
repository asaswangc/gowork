package variable

import (
	"sync"
)

// Global SyncMap 全局容器[这个主要存全局通用的变量]
var Global = syncMap{}

// SyncMap 全局容器
type syncMap struct {
	syncMap sync.Map
}

// Get 从全局容器中获取值
func (this *syncMap) Get(key string) interface{} {
	value, ok := this.syncMap.Load(key)
	if ok {
		return value
	}
	return nil
}

// Set 向全局容器中添加
func (this *syncMap) Set(key string, value interface{}) {
	this.syncMap.Store(key, value)
}
