package models

type Register struct {
	Mailbox string `json:"mailbox"` //邮箱
	Grade   string `json:"grade"`   //年级
}

type Mailbox struct {
	Mailbox string `json:"mailbox"` //邮箱
}

type Revise_User struct {
	Change_User map[string]string `json:"change_user"`
	Mailbox     string            `json:"mailbox"` //要修改数据的邮箱
}

type Blacklist_User struct {
	Uuid     string `json:"uuid"`     //唯一标识符
	Mailbox  string `json:"mailbox"`  //账户
	Password string `json:"password"` //密码(也就是注销原因)
}

type Message struct {
	Message string `json:"message"`
}
