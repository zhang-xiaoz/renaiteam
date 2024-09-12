package models

type User struct {
	Uuid        string //唯一标识符
	Status      int    `gorm:"column:status"` //用户状态码
	Mailbox     string `gorm:"primary_key"`   //邮箱
	Password    string //密码
	Sex         string //性别
	Username    string //用户名//可以重复
	Name        string //姓名
	Address     string //住址
	Grade       string //年级
	Direction   string //方向
	QQ          string `gorm:"column:qq" json:"qq"` //qq
	Wechat      string //微信
	Position    string //职位
	Motto       string //座右铭
	Creat_time  string `gorm:"column:creat_time"`  //创建时间
	Delete_time string `gorm:"column:delete_time"` //删除时间
}

func (User) TableName() string { //gorm更改表名
	return "users"
}

type Register struct {
	Uuid    string `json:"uuid" gorm:"column:uuid"`     //唯一标识符
	Status  int    `json:"status" gorm:"column:status"` //用户状态码
	Mailbox string `json:"mailbox" gorm:"primary_key"`  //邮箱
	Name    string `json:"name" gorm:"column:name"`     //姓名
}

type Normal_User struct {
	Uuid       string `gorm:"column:uuid" json:"uuid"`           //唯一标识符
	Mailbox    string `gorm:"primary_key" json:"mailbox"`        //邮箱
	Sex        string `gorm:"column:sex" json:"sex"`             //性别
	Username   string `gorm:"column:username" json:"username"`   //用户名//可以重复
	Name       string `gorm:"column:name" json:"name"`           //姓名
	Address    string `gorm:"column:adress" json:"address"`      //住址
	Grade      string `gorm:"column:grade" json:"grade"`         //年级
	Direction  string `gorm:"column:direction" json:"direction"` //方向
	QQ         string `gorm:"column:qq" json:"qq"`               //qq
	Wechat     string `gorm:"column:wechat" json:"wechat"`       //微信
	Position   string `gorm:"column:position" json:"position"`   //职位
	Motto      string `gorm:"column:motto" json:"motto"`         //座右铭
	Creat_time string `gorm:"column:creat_time"`                 //创建时间
}

type Cancel_User struct {
	Uuid        string `gorm:"column:uuid" json:"uuid"`           //唯一标识符
	Password    string `gorm:"column:password" json:"password"`   //密码
	Delete_time string `gorm:"column:delete_time"`                //删除时间
	Mailbox     string `gorm:"primary_key" json:"mailbox"`        //邮箱
	Sex         string `gorm:"column:sex" json:"sex"`             //性别
	Username    string `gorm:"column:username" json:"username"`   //用户名//可以重复
	Name        string `gorm:"column:name" json:"name"`           //姓名
	Address     string `gorm:"column:adress" json:"adress"`       //住址
	Grade       string `gorm:"column:grade" json:"grade"`         //年级
	Direction   string `gorm:"column:direction" json:"direction"` //方向
	QQ          string `gorm:"column:qq" json:"qq"`               //qq
	Wechat      string `gorm:"column:wechat" json:"wechat"`       //微信
	Position    string `gorm:"column:position" json:"position"`   //职位
	Motto       string `gorm:"column:motto" json:"motto"`         //座右铭
	Creat_time  string `gorm:"column:creat_time"`                 //创建时间
}

type Blacklist_User struct {
	Mailbox  string `gorm:"primary_key" json:"mailbox"`      //邮箱
	Password string `gorm:"column:password" json:"password"` //密码
}
