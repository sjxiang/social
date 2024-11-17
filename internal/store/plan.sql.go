package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
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


type PlanStoreImpl struct {
	db *sql.DB
}

func (p *PlanStoreImpl) GetAll(ctx context.Context) ([]*Plan, error) {
	query := `
		select id, plan_name, plan_amount, created_at, updated_at
		from plans 
		order by id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	var plans []*Plan

	for rows.Next() {
	
		var i Plan

		err := rows.Scan(
			&i.ID,
			&i.PlanName,
			&i.PlanAmount,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		i.PlanName, i.PlanAmountFormatted = i.ForDisplay()
		
		if err != nil {
			return nil, err
		}

		plans = append(plans, &i)
	}

	return plans, nil
}


func (p *PlanStoreImpl) GetOne(ctx context.Context, planID int64) (*Plan, error) {
	query := `
		select id, plan_name, plan_amount, created_at, updated_at 
		from plans 
		where id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var i Plan
	row := p.db.QueryRowContext(ctx, query, planID)

	err := row.Scan(
		&i.ID,
		&i.PlanName,
		&i.PlanAmount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	i.PlanName, i.PlanAmountFormatted = i.ForDisplay()

	return &i, nil
}

// 订阅
func (p *PlanStoreImpl) SubscribeUserToPlan(ctx context.Context, params User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// 取消用户的所有订阅
	stmt := `delete from user_plans where user_id = ?`	
	
	_, err := p.db.ExecContext(ctx, stmt, params.ID)
	if err != nil {
		return err
	}

	// 如果以前没有订阅, 没有影响 rowsAffected == 0 嘛


	// 创建新的订阅
	stmt = `
		insert into user_plans (user_id, plan_id, created_at, updated_at)
		values (?, ?, ?, ?)
	`

	_, err = p.db.ExecContext(ctx, stmt, 
		params.ID, 
		params.Plan.ID, 
		params.Plan.CreatedAt, 
		params.Plan.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

