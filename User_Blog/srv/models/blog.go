package models

type Blog struct {
	Article_Uuid   string `gorm:"column:article_uuid"` //文章唯一标识符
	User_Uuid      string `gorm:"column:user_uuid"`    //用户唯一标识符
	Title          string //标题
	Content        string //内容
	Label          string `gorm:"column:label"` //标签
	Cover          string //封面
	Abstract       string //摘要
	Reading_Volume int    `gorm:"column:reading_volume"` //阅读量
	Visibility     string //可见范围
	Article_Type   string `gorm:"column:article_type"` //文章类型
	Process        string //是否过审  未审核,已经审核,未通过审核
	Creat_Time     string `gorm:"column:creat_time"` //创建时间
}

func (Blog) TableName() string { //gorm更改表名
	return "blog"
}

type Blog_Collection struct {
	Article_Uuid string `gorm:"column:article_uuid"` //文章唯一标识符
	User_Uuid    string `gorm:"column:user_uuid"`    //用户唯一标识符
	Time         string `gorm:"column:time"`         //创建时间
}

func (Blog_Collection) TableName() string { //gorm更改表名
	return "blog_collection"
}

// elastic所用blog
type Blog_Elastic struct {
	Article_Uuid   string `json:"article_uuid"`   //文章唯一标识符
	User_Uuid      string `json:"user_uuid"`      //用户唯一标识符
	Title          string `json:"title"`          //标题
	Label          string `json:"label"`          //标签
	Abstract       string `json:"abstract"`       //摘要
	Cover          string `json:"cover"`          //封面
	Article_Type   string `json:"article_type"`   //文章类型
	Visibility     string `json:"visibility"`     //可见范围 //仅自己可见 全体可见 仁爱成员可见
	Process        string `json:"process"`        //是否过审  未审核,已经审核,未通过审核,草稿
	Reading_Volume int    `json:"reading_volume"` //阅读量
	Creat_Time     int    `json:"creat_time"`     //创建时间
}

// 仅页面显示所用结构
type Blog_Elastic_Html struct {
	Cover        string `json:"cover"`        //封面
	Article_Uuid string `json:"article_uuid"` //文章唯一标识符
	Title        string `json:"title"`        //标题
	Abstract     string `json:"abstract"`     //摘要
	Visibility   string `json:"visibility"`   //可见范围 //仅自己可见 全体可见 仁爱成员可见
}

// 页面主页使用存储所用结构体
type Blog_Elastic_loadstorage struct {
	Article_Uuid string `json:"article_uuid"` //文章唯一标识符
}

// 仅页面显示所用结构
type Blog_Background_Html struct {
	Article_Uuid   string `json:"article_uuid"`   //文章唯一标识符
	Title          string `json:"title"`          //标题
	Creat_Time     string `json:"creat_time"`     //创建时间
	User_Uuid      string `json:"user_uuid"`      //用户唯一标识符
	Reading_Volume int    `json:"reading_volume"` //阅读量
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
