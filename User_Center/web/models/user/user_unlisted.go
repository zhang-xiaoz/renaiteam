package user

//邮箱验证码获取
type Mailbox_Dode struct {
	Mailbox string `json:"mailbox" binding:"required,email"`
}

//用户注册所填信息
type Mailbox_Register struct {
	Mailbox  string `json:"mailbox" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

//用户实际所有信息
type User struct {
	Status      int    //状态  1表示正常用户 2表示注销过的用户 3表示未通过审核  4表示改用户不允许注册已经被拉黑
	Uuid        string //唯一标识符
	Mailbox     string //邮箱
	Password    string //密码
	Username    string //用户名
	Name        string //姓名
	Sex         string //性别
	Address     string //地址
	Grade       string //年级
	Direction   string //方向
	QQ          string `json:"qq"` //qq
	Wechat      string //微信
	Position    string //职位
	Motto       string //座右铭
	Creat_time  string //创建时间
	Delete_time string //注销时间
}

//用户登录所填信息
type User_login struct {
	Mailbox  string `json:"mailbox" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Id       string `json:"id" binding:"required"`
	B64s     string `json:"b64s" binding:"required"`
}

//用户验证码验证
type User_Code_Check struct {
	Mailbox string `json:"mailbox" binding:"required,email"`
	Code    string `json:"code" binding:"required"`
}

//用户更改密码
type User_Code_Twice_Check struct {
	Mailbox  string `json:"mailbox" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}
