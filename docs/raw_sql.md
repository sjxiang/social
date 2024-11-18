




# raw sql

> "databse/sql", 定义接口 
> "github.com/go-sql-driver/mysql", 负责具体实现


通过 sql.Open("mysql", ${dsn}) 获取一个数据源操作对象 db
开启事务时, 通过 db.Begin() 或 db.BeginTx(ctx, &sql.TxOptions{}) 获得事务操作对象 tx.
执行 sql 查询使用 tx.Query;
执行 sql 新增、修改、删除, 使用 tx.Exec;
最后, 使用 tx.Commit() 提交或使用 tx.Rollback() 回滚.



# "database/sql" 细节

```go

// 1. Result
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}


// 2. Rows
func (rs *Rows) Close() error                    // 关闭结果集
func (rs *Rows) Err() ([]string, error)          // 错误集
func (rs *Rows) Next() bool                      // 游标, 下一条记录
func (rs *Rows) Scan(dest ...interface{}) error  // 扫描 struct

```
