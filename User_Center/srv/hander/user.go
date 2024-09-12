package hander

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"srv/initialization"
	"srv/log"
	"srv/models"
	"srv/proto"
	"strconv"

	redisss "github.com/gomodule/redigo/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 注意  在网络传输中，不能两个都为nil否者有问题
type UserServer struct {
	*proto.UnimplementedUsersServer
}

func (s *UserServer) Register_Mailbox_Back(ctx context.Context, req *proto.Mailbox) (*proto.Mailbox_Back, error) {
	backk := proto.Mailbox_Back{} //需要返回的数据，记住返回的数据不能两个都为空
	//查询mysql数据库
	backk.Mailbox = req.Mailbox
	result := initialization.DB.Table("users").Select("uuid,status").Where("mailbox = ?", req.Mailbox).First(&backk)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "用户不存在")
		}
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	return &backk, nil
}

func (s *UserServer) Get_Redis_Storage(ctx context.Context, req *proto.Redis_Storage) (*proto.Redis_Storage, error) {
	//先进行连接
	var back proto.Redis_Storage
	rediss := initialization.Redis.Get()
	if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		return nil, err
	}
	defer rediss.Close() //关闭数据库连接
	//更改数据库
	_, err := rediss.Do("select", req.Who)
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, err
	}
	//获取数据
	value1, _ := redisss.String(rediss.Do("get", req.Key))
	log.SugarLogger.Error(req.Key)
	log.SugarLogger.Error(value1)
	if value1 == "" {
		return nil, status.Errorf(codes.NotFound, "不存在该信息")
	}
	//获取数据
	value2, err := redisss.Int(rediss.Do("ttl", req.Key))
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, err
	}
	back.Key = req.Key
	back.Value = value1
	back.Who = req.Who
	back.Time = int64(value2)
	return &back, nil
}

func (s *UserServer) Save_Redis_Storage(ctx context.Context, req *proto.Redis_Storage) (*proto.BACK, error) { //往redis里边存储存储信息
	//先进行连接
	var back proto.BACK
	back.Back = true
	rediss := initialization.Redis.Get()
	if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		back.Back = false
		return &back, err
	}
	defer rediss.Close() //关闭数据库连接
	//更改数据库
	_, err := rediss.Do("select", req.Who)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	_, err = rediss.Do("setex", req.Key, req.Time, req.Value)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	return &back, nil
}

// func (s *UserServer) Renew_User_Message(ctx context.Context, req *proto.User_Message) (*proto.BACK, error) {
// 	back := proto.BACK{}
// 	//数据转换
// 	u := models.User{
// 		Uuid:        req.Uuid,
// 		Status:      int(req.Status),
// 		Mailbox:     req.Mailbox,
// 		Password:    req.Password,
// 		Sex:         req.Sex,
// 		Username:    req.Username,
// 		Name:        req.Name,
// 		Address:     req.Address,
// 		Creat_time:  req.CreatTime,
// 		Delete_time: req.DeleteTime,
// 	}
// 	//开启事务
// 	tx := initialization.DB.Begin()                              //开始事务
// 	err := tx.Model(&u).Where("mailbox=?", u.Mailbox).Updates(u) //尝试是否能简化成功
// 	if err.Error != nil {
// 		log.SugarLogger.Error(err)
// 		tx.Rollback() // 发生错误回滚事务
// 		back.Back = false
// 		return &back, err.Error
// 	}
// 	tx.Commit() //提交事务
// 	back.Back = true
// 	return &back, nil
// }

func (s *UserServer) Add_User_Message(ctx context.Context, req *proto.User_Message) (*proto.BACK, error) {
	back := proto.BACK{}
	//数据转换
	u := models.User{
		Uuid:        req.Uuid,
		Status:      int(req.Status),
		Mailbox:     req.Mailbox,
		Password:    req.Password,
		Sex:         req.Sex,
		Username:    req.Username,
		Name:        req.Name,
		Address:     req.Address,
		Grade:       req.Grade,
		Direction:   req.Direction,
		QQ:          req.Qq,
		Wechat:      req.Wechat,
		Position:    req.Position,
		Motto:       req.Motto,
		Creat_time:  req.CreatTime,
		Delete_time: req.DeleteTime,
	}
	//开启事务
	fmt.Println(u)
	tx := initialization.DB.Begin() //开始事务
	err := tx.Create(&u)
	if err.Error != nil {
		log.SugarLogger.Error(err)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, err.Error
	}
	tx.Commit() //提交事务
	back.Back = true
	return &back, nil
}

func (s *UserServer) Delete_Redis_Storage(ctx context.Context, req *proto.Redis_Storage) (*proto.BACK, error) {
	//先进行连接
	var back proto.BACK
	back.Back = true
	rediss := initialization.Redis.Get()
	if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		back.Back = false
		return &back, err
	}
	defer rediss.Close() //关闭数据库连接
	//更改数据库
	_, err := rediss.Do("select", req.Who)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	_, err = rediss.Do("del", req.Key)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	return &back, nil
}

func (s *UserServer) Get_User_Mesaage_Mysql(ctx context.Context, req *proto.Mailbox) (*proto.User_Message, error) {
	back := proto.User_Message{}
	back.Mailbox = req.Mailbox
	result := initialization.DB.Table("users").Where("mailbox=?", req.Mailbox).First(&back)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "用户不存在")
		}
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	return &back, nil
}

func (s *UserServer) Revise_User_Message_Mysql(ctx context.Context, req *proto.User_Change) (*proto.BACK, error) {
	fmt.Println(req)
	back := proto.BACK{}
	u1 := models.User{} //实体化一个对象
	u2 := make(map[string]interface{})
	var mailbox string
	for key, value := range req.UserMessage {
		if key == "mailbox" {
			mailbox = value
			continue
		}
		u2[key] = value
	}
	fmt.Println(u2)
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	err := tx.Model(&u1).Where("mailbox=?", mailbox).Updates(u2)
	if err.Error != nil {
		log.SugarLogger.Error(err.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, err.Error
	}
	tx.Commit() //成功提交事务
	back.Back = true
	return &back, nil
}

func (s *UserServer) Get_User_Mailbox(ctx context.Context, req *proto.Mailbox) (*proto.Mailbox, error) {
	mailbox := ""
	result := initialization.DB.Table("users").Select("mailbox").Where("uuid=?", req.Mailbox).Scan(&mailbox)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "用户不存在")
		}
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	return &proto.Mailbox{Mailbox: mailbox}, nil
}

func (s *UserServer) Revise_User_Password_Mysql(ctx context.Context, req *proto.User_Password) (*proto.BACK, error) {
	back := proto.BACK{}
	u := models.User{} //实体化一个对象
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&u).Where("uuid=?", req.Uuid).Updates(models.User{Password: req.Password})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	//将redis里边数据剔除
	rediss := initialization.Redis.Get()
	if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, err
	}
	defer rediss.Close() //关闭数据库连接
	//更改数据库
	_, err := rediss.Do("select", "1")
	if err != nil {
		log.SugarLogger.Error(err)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, err
	}
	//删除数据
	_, err = rediss.Do("del", req.Uuid+"access_token", req.Uuid+"refresh_token")
	if err != nil {
		log.SugarLogger.Error(err)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, err
	}
	tx.Commit() //成功提交事务
	back.Back = true
	return &back, nil
}

func (s *UserServer) Revise_User_Status_Mysql(ctx context.Context, req *proto.Mailbox_Back) (*proto.BACK, error) {
	back := proto.BACK{}
	u := models.User{} //实体化一个对象
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&u).Where("uuid=?", req.Uuid).Updates(models.User{Status: int(req.Status)})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	tx.Commit() //成功提交事务
	back.Back = true
	return &back, nil
}

func (s *UserServer) Revise_User_Mailbox_Mysql(ctx context.Context, req *proto.Mailbox_Back) (*proto.BACK, error) {
	back := proto.BACK{}
	u := models.User{} //实体化一个对象
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&u).Where("uuid=?", req.Mailbox).Updates(models.User{Mailbox: req.Uuid})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	tx.Commit() //成功提交事务
	back.Back = true
	return &back, nil
}

func (s *UserServer) Get_User_Password_Mysql(ctx context.Context, req *proto.Mailbox) (*proto.Mailbox, error) {
	p := proto.Mailbox{}
	result := initialization.DB.Table("users").Select("password").Where("uuid= ?", req.Mailbox).Scan(&p.Mailbox)
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}

func (s *UserServer) Delete_Mysql_Mailbox(ctx context.Context, req *proto.Mailbox) (*proto.BACK, error) {
	back := proto.BACK{}
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Where("mailbox = ?", req.Mailbox).Delete(&models.User{})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	tx.Commit() //成功提交事务
	back.Back = true
	return &back, nil
}

func (s *UserServer) Get_Register_User(ctx context.Context, req *proto.MemberPaging) (*proto.Message, error) {
	back1 := []models.Register{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("users").Select("uuid,status,mailbox,name").Where("status=?", 3).Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("status = ?", 3).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	type Back struct {
		Member string
		Number string
	}
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.Message{
		Message: string(bb),
	}, nil
}

func (s *UserServer) Revise_User_Status(ctx context.Context, req *proto.User_Revise_Status) (*emptypb.Empty, error) {
	//开始修改用户状态
	u := models.User{} //实体化一个对象
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&u).Where("mailbox=?", req.Mailbox).Updates(models.User{Status: int(req.Status), Grade: req.Grade})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, errors.New("没有找到数据")
	}
	tx.Commit() //成功提交事务
	return &emptypb.Empty{}, nil
}

func (s *UserServer) Refuse_User_Status(ctx context.Context, req *proto.Mailbox) (*emptypb.Empty, error) {
	//根据邮箱删除数据
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Where("mailbox = ?", req.Mailbox).Delete(&models.User{})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, errors.New("没有找到数据")
	}
	tx.Commit() //成功提交事务
	return &emptypb.Empty{}, nil
}

func (s *UserServer) Get_User(ctx context.Context, req *proto.MemberPaging) (*proto.Message, error) {
	back1 := []models.Normal_User{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("users").Where("status=?", 1).Order("grade DESC").Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("status = ?", 1).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	type Back struct {
		Member string
		Number string
	}
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.Message{
		Message: string(bb),
	}, nil
}

func (s *UserServer) Seek_Grade_User(ctx context.Context, req *proto.MemberPaging) (*proto.Message, error) {
	//开始查找数据
	back1 := []models.Normal_User{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("users").Where("status=? AND grade=?", 1, req.Message).Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("status=? AND grade=?", 1, req.Message).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	type Back struct {
		Member string
		Number string
	}
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.Message{
		Message: string(bb),
	}, nil
}

func (s *UserServer) Seek_Name_User(ctx context.Context, req *proto.MemberPaging) (*proto.Message, error) {
	//开始查找数据
	back1 := []models.Normal_User{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("users").Where("status=? AND name LIKE ?", 1, req.Message).Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("status=? AND name LIKE ?", 1, req.Message).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	type Back struct {
		Member string
		Number string
	}
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.Message{
		Message: string(bb),
	}, nil
}

func (s *UserServer) Get_Cancel_User(ctx context.Context, req *proto.MemberPaging) (*proto.Message, error) {
	//开始查找数据
	back1 := []models.Cancel_User{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("users").Where("status=?", 2).Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("status=?", 2).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	type Back struct {
		Member string
		Number string
	}
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.Message{
		Message: string(bb),
	}, nil
}

func (s *UserServer) Get_Blacklist_User(ctx context.Context, req *proto.MemberPaging) (*proto.Message, error) {
	//开始查找数据
	back1 := []models.Blacklist_User{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("users").Where("status=?", 4).Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("status=?", 4).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	type Back struct {
		Member string
		Number string
	}
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.Message{
		Message: string(bb),
	}, nil
}

func (s *UserServer) Delete_Blacklist_User(ctx context.Context, req *proto.Mailbox) (*emptypb.Empty, error) {
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Where("mailbox = ?", req.Mailbox).Delete(&models.User{})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		return &emptypb.Empty{}, result.Error
	}
	tx.Commit() //成功提交事务
	return &emptypb.Empty{}, nil
}

func (s *UserServer) Add_Blacklist_User(ctx context.Context, req *proto.Blacklist_User) (*emptypb.Empty, error) {
	//开启事务
	user := models.User{
		Uuid:     req.Uuid,
		Status:   int(req.Status),
		Mailbox:  req.Mailbox,
		Password: req.Password,
	}
	tx := initialization.DB.Begin() //开始事务
	result := tx.Create(&user)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		return &emptypb.Empty{}, result.Error
	}
	tx.Commit() //成功提交事务
	return &emptypb.Empty{}, nil
}

func (s *UserServer) Add_Register_Blacklist(ctx context.Context, req *proto.Blacklist_User) (*emptypb.Empty, error) {
	u := models.User{} //实体化一个对象
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&u).Where("mailbox=?", req.Mailbox).Updates(models.User{Status: int(req.Status), Password: req.Password})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		return &emptypb.Empty{}, result.Error
	}
	tx.Commit() //成功提交事务
	return &emptypb.Empty{}, nil
}
