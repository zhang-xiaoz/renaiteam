package models

//用户登录所填信息
type User_login struct {
	Mailbox  string `json:"mailbox" binding:"required"`
	Password string `json:"password" binding:"required"`
	Id       string `json:"id" binding:"required"`
	B64s     string `json:"b64s" binding:"required"`
}

type Grade struct {
	Grade string `json:"grade"`
}

type Name struct {
	Name string `json:"name"`
}
