package data

import (
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)


type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	IsActive  bool      `json:"is_active"`  
	Role      int64     `json:"role"` // 角色, 0 游客, 1 用户, 2 版主, 3 管理员
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Plan      Plan      `json:"plan"`
}


type password struct {
	text *string  
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}


// 会员订阅计划 subscription plans
type Plan struct {
	ID                  int       `json:"id"`
	PlanName            string    `json:"plan_name"`
	PlanAmount          int       `json:"plan_amount"`  // 金额, 单位: 元
	PlanAmountFormatted string    `json:"plan_amount_formatted"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// 格式化为货币字符串
func (p *Plan) AmountForDisplay() string {
	return fmt.Sprintf("¥%d", p.PlanAmount)
}


func (p *Plan) PlanNameForDisplay() string {
	planName := p.PlanName

	if planName == "free" {
		return "免费"
	} else if planName == "basic" {
		return "基础"
	} else if planName == "pro" {
		return "专业"
	} else if planName == "enterprise" {
		return "企业"
	} else {
		return "未知"
	}
}

func (p *Plan) ForDisplay() (string, string) {
	return p.PlanNameForDisplay(), p.AmountForDisplay()
}

