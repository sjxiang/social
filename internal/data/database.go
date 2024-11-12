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
		// 动态
		Post interface {
			Create(context.Context) error
		}
		
		// 评论
		Comments interface {
			Create(context.Context) error
		}

		// 用户
		User interface {
			GetByID(context.Context, int64) (*User, error)
			GetByEmail(context.Context, string) (*User, error)
			Exists(context.Context, int64) (bool, error)
			Create(context.Context, *sql.Tx, *User) error
			CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
			Activate(context.Context, string) error
			Delete(context.Context, int64) error
			ModPassword(context.Context, *User) error
		}
	
		// 角色
		Role interface {
			GetByName(context.Context, string) (*Role, error)
		}

}

func NewMySQLStorage(db *sql.DB) MySQLStorage {
	return MySQLStorage{
		// UserStore:     &MySQLUserStore{db},
		Role:   NewMySQLRoleStore(db),
	}
}

/*

使用 go 操作数据库时, 我们会使用到 go 语言的官方库 database/sql

通过 sql.Open("mysql", ${dsn}) 获取一个数据源操作对象 db.
开启事务时, 使用 db.Begin() 或 db.BeginTx(ctx, &sql.TxOptions{}) 获得事务操作对象 tx, 
执行 sql 查询使用 tx.Query, 执行 sql 新增、修改、删除, 使用 tx.Exec；最后使用 tx.Commit() 提交
或使用 tx.Rollback() 回滚。

 */
