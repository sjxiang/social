package validate

import (
	"testing"
	"time"
)


func TestValidateParams(t *testing.T) {
	t.Log(time.DateTime)

	var params struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"` 
	}

	params.Email    = "gua@vip.cn"
	params.Password = "1nidqkbdiovcf2"

	if err := Check(params); err != nil {
		t.Fatal(err)
	}

	t.Log("符合规则")
}