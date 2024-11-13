package data


// import (
// 	"context"
// 	"database/sql"
// )

// type MySQLRoleStore struct {
// 	db *sql.DB
// }

// func NewMySQLRoleStore(db *sql.DB) *MySQLRoleStore {
// 	return &MySQLRoleStore{db: db}
// }

// func (m *MySQLRoleStore) GetByName(ctx context.Context, slug string) (*Role, error) {
// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	query := `
// 		SELECT 
// 			id, name, description, level 
// 		FROM 
// 			roles 
// 		WHERE 
// 			name = ?`
	
// 	var i Role

// 	row := m.db.QueryRowContext(ctx, query, slug)
	
// 	err := row.Scan(
// 		&i.ID, 
// 		&i.Name, 
// 		&i.Description, 
// 		&i.Level,
// 	)
	
// 	if err!= nil {
// 		return nil, err
// 	}

// 	return &i, nil
// }