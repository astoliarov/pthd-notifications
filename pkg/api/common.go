package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"reflect"
	"strings"
)

func initializeValidator() *validator.Validate {
	var validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate
}

func initializeDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	return decoder
}
