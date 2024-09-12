package user

type JWT_Check struct { //Jwt验证附加一些信息(后边变成了很多函数所用结构体)
	Access_token  string            `json:"access_token"`     //可以设置个长度
	Refresh_token string            `json:"refresh_token"`    //可以设置个长度
	UUID          string            `json:"uuid" form:"uuid"` //可以设置个长度
	Change_User   map[string]string `json:"change_user"`
	Code          string            `json:"code"`        //验证码
	Reason        string            `json:"reason"`      //注销理由
	OldPassword   string            `json:"oldpassword"` //原来的密码
	NewPassword   string            `json:"newpassword"` //要被修改的密码
	Mailbox       string            `json:"mailbox"`     //新邮箱//要修改的邮箱
}
