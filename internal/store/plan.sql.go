package store

import (
	"time"
	"fmt"
)

// 会员订阅计划
type Plan struct {
	ID                  int       `json:"id"`
	PlanName            string    `json:"plan_name"`
	PlanAmount          int       `json:"plan_amount"`  // 金额, 单位: 分
	PlanAmountFormatted string    `json:"plan_amount_formatted"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// 展示货币金额
func (p *Plan) AmountForDisplay() string {
	amount := float64(p.PlanAmount) / 100.00
	return fmt.Sprintf("¥%.2f", amount)
}

// 展示 i18n 计划名词
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

