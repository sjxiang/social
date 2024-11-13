package data

import (
	"context"
)


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
		plan.PlanAmountFormatted = plan.AmountForDisplay()
		
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

	i.PlanAmountFormatted = i.AmountForDisplay()

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

