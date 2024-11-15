package data

import (
	"context"
)


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



func (m *MySQLPlanStore) GetAll(ctx context.Context) ([]*Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		select 
			id, plan_name, plan_amount, created_at, updated_at
		from 
			plans 
		order by 
			id`

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*Plan

	for rows.Next() {
		var plan Plan

		err := rows.Scan(
			&plan.ID,
			&plan.PlanName,
			&plan.PlanAmount,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		plan.PlanName, plan.PlanAmountFormatted = plan.ForDisplay()
		
		if err != nil {
			return nil, err
		}

		plans = append(plans, &plan)
	}

	return plans, nil
}


func (m *MySQLPlanStore) GetOne(ctx context.Context, id int64) (*Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		select 
			id, plan_name, plan_amount, created_at, updated_at 
		from 
			plans 
		where 
			id = ?`

	var i Plan
	row := m.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&i.ID,
		&i.PlanName,
		&i.PlanAmount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	i.PlanName, i.PlanAmountFormatted = i.ForDisplay()

	return &i, nil
}

// 订阅用户到计划
func (m *MySQLPlanStore) SubscribeUserToPlan(ctx context.Context, arg User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// 取消用户的所有订阅
	stmt := `delete from user_plans where user_id = ?`
	_, err := m.db.ExecContext(ctx, stmt, arg.ID)
	if err != nil {
		return err
	}

	// 创建新的订阅
	stmt = `insert into user_plans (user_id, plan_id, created_at, updated_at)
			values (?, ?, ?, ?)`

	_, err = m.db.ExecContext(ctx, stmt, arg.ID, arg.Plan.ID, arg.Plan.CreatedAt, arg.Plan.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

