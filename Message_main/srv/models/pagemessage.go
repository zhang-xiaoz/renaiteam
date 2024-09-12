package models

//首页成员信息
type MemberMessage struct {
	Name      string `json:"name"`      //姓名
	Sex       string `json:"sex"`       //性别
	QQ        string `json:"qq"`        //qq
	Grade     string `json:"grade"`     //年级
	Direction string `json:"direction"` //方向
	Motto     string `json:"motto"`     //座右铭
}

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

type Prize struct {
	Uuid   string `json:"uuid"`   //唯一标识符
	Time   string `json:"time"`   //获取时间
	Name   string `json:"name"`   //谁获取的
	Awards string `json:"awards"` //奖项名称
}

func (Prize) TableName() string { //gorm更改表名
	return "prize"
}

type ClubDirection struct {
	Name      string `json:"name"`      //方向名称
	Introduce string `json:"introduce"` //方向介绍
	Picture   string `json:"picture"`   //图片
}

func (ClubDirection) TableName() string { //gorm更改表名
	return "club_direction"
}

type Training_Personnel_Message struct {
	Name  string `json:"name"`                //姓名
	Sex   string `json:"sex"`                 //性别
	QQ    string `gorm:"column:qq" json:"qq"` //qq
	Grade string `json:"grade"`               //年级
}

type Message struct {
	Key     string `json:"key"`
	Keyform string `json:"keyform"`
	Value1  string `json:"value1"`
	Value2  string `json:"value2"`
}

func (Message) TableName() string { //gorm更改表名
	return "message"
}
