package rtstruct

import (
	"fmt"
)

type Tag struct {
	seq    []string
	tagStr string
	tagMap map[string]string
}

func NewTag() *Tag {
	return &Tag{
		tagStr: `gorm:"`,
		seq:    []string{PrimaryKey, AutoIncrement, Null, NotNull, Default, Comment, UniqueKey},
	}
}

func (t *Tag) TagMap(tags map[string]interface{}) *Tag {
	for i := 0; i < len(t.seq); i++ {
		if value, ok := tags[t.seq[i]]; ok {
			if value == "" {
				t.tagStr = t.tagStr + fmt.Sprintf(`%s`, t.seq[i])
			}
			if value != "" {
				t.tagStr = t.tagStr + fmt.Sprintf(`%s:%s`, t.seq[i], value)
			}
		}
	}
	return t
}
