package data

import 	"golang.org/x/crypto/bcrypt"


type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleID    int64    `json:"role_id"` 
	Role      Role     `json:"role"`
}

func (u *User) SetRole() {
	
	// 遍历 RoleSet, 如果用户角色存在, 则设置用户角色
	for _, e := range RoleSet {
		if e.Name == u.Role.Name {
			// 设置用户角色
			u.RoleID = e.ID
			u.Role = e
			break
		}
	}

	// 如果用户角色不存在, 则设置为游客
	if u.Role.IsAdmin() && u.Role.IsModerator() && u.Role.IsUser() && u.Role.IsGuest() {
		u.RoleID = 4
		u.Role = RoleSet[3]
	}

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

var RoleSet = []Role{
	{	
		ID: 1,
		Name: "admin",
		Level: 4,
		Description: "管理员, 可以修改和删掉阿婆主的动态",
	},
	{
		ID: 2,
		Name: "moderator",
		Level: 3,
		Description: "版主, 可以修改阿婆主的动态",
	},
	{
		ID: 3,
		Name: "user",
		Level: 2,
		Description: "用户, 可以发动态和评论",
	},
	{
		ID: 4,
		Name: "guest",
		Level: 1,
		Description: "游客, 只能浏览动态",
	},
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int64  `json:"level"`
}


func (role *Role) IsAdmin() bool {
	return role.Name == "admin"
}

func (role *Role) IsModerator() bool {
	return role.Name == "moderator"
}

func (role *Role) IsUser() bool {
	return role.Name == "user"
}

func (role *Role) IsGuest() bool {
	return role.Name == "guest"

}