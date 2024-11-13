package data

import (
	"time"
	"context"
	"database/sql"
)


var (
	QueryTimeoutDuration = time.Second * 5
)


type MySQLStorage struct {
	
		// 用户
		Users interface {
			GetOne(ctx context.Context, id int64) (*User, error)
			Exists(ctx context.Context, id int64) (bool, error)
			GetByEmail(ctx context.Context, email string) (*User, error)
			Delete(ctx context.Context, id int64) error

			CreateAndInvite(ctx context.Context, arg User, token string, invitationExp time.Duration) error
			Activate(ctx context.Context, token string) error

			ModPassword(ctx context.Context, arg User) error
		}

		// 帖子


		// 评论


	
}

func NewMySQLStorage(db *sql.DB) MySQLStorage {
	return MySQLStorage{
		Users: &MySQLUserStore{db: db},
	}
}


type MySQLUserStore struct {
	db *sql.DB
}


type MySQLPostStore struct {
	db *sql.DB
}

type MySQLCommentStore struct {
	db *sql.DB
}

type MySQLPlanStore struct {
	db *sql.DB
}


/*

使用 go 操作数据库时, 我们会使用到 go 语言的官方库 database/sql

通过 sql.Open("mysql", ${dsn}) 获取一个数据源操作对象 db.
开启事务时, 使用 db.Begin() 或 db.BeginTx(ctx, &sql.TxOptions{}) 获得事务操作对象 tx, 
执行 sql 查询使用 tx.Query, 执行 sql 新增、修改、删除, 使用 tx.Exec；最后使用 tx.Commit() 提交
或使用 tx.Rollback() 回滚。

 */
