package rtstruct

import (
	"errors"
	"reflect"
	"strings"
)

type TagInterface interface {
	TagStr() string
	TagMap(attr string, value string) TagInterface
}

type Struct struct {
	fields []reflect.StructField
}

// NewStruct 构造函数
func NewStruct() *Struct {
	return &Struct{}
}

// AddField 添加字段
func (rt *Struct) AddField(fieldName string, fieldType SqlType, tags map[string]interface{}) *Struct {
	rt.fields = append(rt.fields, reflect.StructField{
		Name: strings.ToUpper(fieldName[:1]) + fieldName[1:],
		Type: reflect.TypeOf(SqlMappingGo[fieldType]),
		Tag:  reflect.StructTag(NewTag().TagMap(map[string]interface{}{Type: fieldType, Column: fieldName}).TagMap(tags).tagStr),
	})
	return rt
}

// Build 生成结构体
func (rt *Struct) Build() *Instance {
	var (
		index        = make(map[string]int)
		instanceType reflect.Type
	)
	instanceType = reflect.StructOf(rt.fields)
	for i := 0; i < instanceType.NumField(); i++ {
		index[instanceType.Field(i).Name] = i
	}
	return &Instance{Stu: rt, InstanceValue: reflect.New(instanceType).Elem(), InstanceType: instanceType, index: index}
}

// Instance 生成结构体的实例
type Instance struct {
	Stu           *Struct
	index         map[string]int
	InstanceType  reflect.Type
	InstanceValue reflect.Value
}

var (
	FieldNoExist = errors.New("field no exist")
	FieldTypeErr = errors.New("field type error")
)

// Struct 获取结构体实例
func (in *Instance) Struct() interface{} {
	return in.InstanceValue.Addr().Interface()
}

// Field 获取结构体字段
func (in *Instance) Field(name string) (reflect.Value, error) {
	if i, ok := in.index[name]; ok {
		return in.InstanceValue.Field(i), nil
	} else {
		return reflect.Value{}, FieldNoExist
	}
}

// Value 用于给生成的结构体赋值
func (in *Instance) Value(name string, value interface{}) error {
	if i, ok := in.index[name]; ok {
		if in.InstanceValue.Field(i).Kind() != reflect.TypeOf(value).Kind() {
			return FieldTypeErr
		}
		in.InstanceValue.Field(i).Set(reflect.ValueOf(value))
		return nil
	}
	return FieldNoExist
}
