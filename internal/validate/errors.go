package validate

import "encoding/json"

// 字段错误
type FieldError struct {
	Field string `json:"field"`  // 字段名
	Err   string `json:"error"`  // 错误信息
}

// 集合
type FieldErrors []FieldError

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}
