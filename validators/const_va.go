package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var IpAddrVf = func() (tag string, fn validator.Func) {
	return "asd", func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if ok {
			today := time.Now()
			if today.After(date) {
				return false
			}
		}
		return true
	}
}
