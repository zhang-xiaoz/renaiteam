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
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 注意  在网络传输中，不能两个都为nil否者有问题
type PagemessageServer struct {
	*proto.UnimplementedPagemessageServer
}

func (*PagemessageServer) Revise_Message(ctx context.Context, req *proto.ReviseMessage) (*emptypb.Empty, error) {
	//开始事务
	tx := initialization.DB.Begin() //开始事务
	//开始修改数据

	result1 := tx.Table("message").Where(" `key` = ? AND keyform = ?", req.Key, req.Keyform).Updates(map[string]interface{}{
		"value1": req.Value1,
		"value2": req.Value2,
	})
	fmt.Println(req.Key)
	fmt.Println(req.Keyform)
	fmt.Println(result1.Error)
	if result1.RowsAffected == 0 {
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, errors.New("没有找到数据")
	}
	if result1.Error != nil {
		log.SugarLogger.Error(result1.Error)
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, result1.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *PagemessageServer) GetMemberMessage(ctx context.Context, req *proto.PageMemberPaging) (*proto.MessagePage, error) {
	//开始查询
	back1 := []models.MemberMessage{}
	//判断当前时间是不是超过9月
	now := time.Now()
	targetDate := time.Date(now.Year(), 9, 1, 0, 0, 0, 0, now.Location())
	grade := time.Now().Year()
	if !now.After(targetDate) { //没超过
		grade = grade - 1
	}
	if req.Message == "大一成员" {

	} else if req.Message == "大二成员" {
		grade = grade - 1
	} else if req.Message == "大三成员" {
		grade = grade - 2
	} else if req.Message == "大四成员" {
		grade = grade - 3
	} else { //历年成员
		var yearsToExclude = []string{strconv.Itoa(grade), strconv.Itoa(grade - 1), strconv.Itoa(grade - 2), strconv.Itoa(grade - 3)} // 要排除的年份列表
		// 计算偏移量
		offset := (req.P - 1) * req.Pn
		// 使用Limit和Offset进行分页查询
		initialization.DB.Table("users").Select("name", "sex", "qq", "direction", "motto").Where("grade  NOT IN (?)", yearsToExclude).Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
		//创建一个新结构体
		type Back struct {
			Member string
			Number string
			P      string
		}
		var count int64
		if req.P == 1 { //查找一共有多少数据
			initialization.DB.Model(&models.User{}).Where("grade  NOT IN (?)", yearsToExclude).Count(&count)
		}
		jsonback1, _ := json.Marshal(back1)
		backk := Back{
			Member: string(jsonback1),
			Number: strconv.FormatInt(count, 10),
			P:      strconv.FormatInt(req.P, 10),
		}
		bb, _ := json.Marshal(backk)
		return &proto.MessagePage{
			Message: string(bb),
		}, nil
	}
	gradestring := strconv.Itoa(grade) // 使用模运算符
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	initialization.DB.Table("users").Select("name", "sex", "qq", "direction", "motto").Where("grade = ?", gradestring).Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	//创建一个新结构体
	type Back struct {
		Member string
		Number string
		P      string
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("grade = ?", gradestring).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
		P:      strconv.FormatInt(req.P, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.MessagePage{
		Message: string(bb),
	}, nil
}

func (s *PagemessageServer) GetPrizeMessage(ctx context.Context, req *proto.PageMemberPaging) (*proto.MessagePage, error) {
	back1 := []models.Prize{}
	//开始查询获奖信息根据时间降序
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	initialization.DB.Order("time desc").Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	//创建一个新结构体
	type Back struct {
		Member string
		Number string
		P      string
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.Prize{}).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
		P:      strconv.FormatInt(req.P, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.MessagePage{
		Message: string(bb),
	}, nil
}

func (s *PagemessageServer) GetClubDirectionMessage(ctx context.Context, req *emptypb.Empty) (*proto.MessagePage, error) {
	back1 := []models.Message{}
	// 查询数据
	initialization.DB.Where("keyform = ?", "主攻方向").Find(&back1)
	fmt.Println(back1)
	jsonback1, _ := json.Marshal(back1)
	return &proto.MessagePage{
		Message: string(jsonback1),
	}, nil
}

func (s *PagemessageServer) GetTrainingPersonnel(ctx context.Context, req *proto.PageMemberPaging) (*proto.MessagePage, error) {
	//开始查询
	back1 := []models.Training_Personnel_Message{}
	//判断当前时间是不是超过9月
	now := time.Now()
	targetDate := time.Date(now.Year(), 9, 1, 0, 0, 0, 0, now.Location())
	grade := time.Now().Year() - 1
	if !now.After(targetDate) { //没超过
		grade = grade - 1
	}
	gradestring := strconv.Itoa(grade) // 使用模运算符
	fmt.Println(gradestring)
	// 使用Limit和Offset进行分页查询
	initialization.DB.Table("users").Select("name", "sex", "qq", "grade").Where("grade = ?", gradestring).Find(&back1)
	//创建一个新结构体
	type Back struct {
		Member string
		Number string
		P      string
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.User{}).Where("grade = ?", gradestring).Count(&count)
	}
	jsonback1, _ := json.Marshal(back1)
	backk := Back{
		Member: string(jsonback1),
		Number: strconv.FormatInt(count, 10),
		P:      strconv.FormatInt(req.P, 10),
	}
	bb, _ := json.Marshal(backk)
	return &proto.MessagePage{
		Message: string(bb),
	}, nil
}

func (s *PagemessageServer) GetTrainingMessage(ctx context.Context, req *emptypb.Empty) (*proto.MessagePage, error) {
	back1 := []models.Message{}
	// 查询数据
	initialization.DB.Where("keyform = ?", "培训信息").Find(&back1)
	jsonback1, _ := json.Marshal(back1)
	return &proto.MessagePage{
		Message: string(jsonback1),
	}, nil
}

func (s *PagemessageServer) GetTrainingTime(ctx context.Context, req *emptypb.Empty) (*proto.MessagePage, error) {
	back1 := []models.Message{}
	back2 := []models.Message{}
	// 查询数据
	initialization.DB.Where("keyform = ?", "培训时间").Find(&back1)
	//开始处理顺序问题
	for i := 1; i <= 6; i++ {
		for _, b := range back1 {
			if string(b.Key[0]) == strconv.FormatInt(int64(i), 10) {
				b.Key = b.Key[1:]
				back2 = append(back2, b)
				break
			}
		}
	}
	//获取当前时间戳
	timestampseconds := time.Now().Unix()
	//返回的在哪
	back3 := 0
	//开始处理时间 时间格式2024-06-03 并且判断时间在哪里
	for _, bb := range back2 {
		t1, err := time.Parse("2006-01-02", bb.Value1)
		if err != nil {
			log.SugarLogger.Error(err)
			return nil, err
		}
		t2, err := time.Parse("2006-01-02", bb.Value2)
		if err != nil {
			log.SugarLogger.Error(err)
			return nil, err
		}
		// 将时间转换成 Unix 时间戳（秒）
		timestamp1 := t1.Unix()
		timestamp2 := t2.Unix()
		if timestamp1 >= timestampseconds {
			break
		}
		back3++
		if timestamp1 <= timestampseconds && timestampseconds <= timestamp2 {
			break
		}
		back3++
	}
	jsonback2, _ := json.Marshal(back2)
	jsonback3, _ := json.Marshal(back3)
	//创建一个新结构体
	type Back struct {
		Member string
		Number string
	}
	backk := Back{
		Member: string(jsonback2),
		Number: string(jsonback3),
	}
	bb, _ := json.Marshal(backk)
	return &proto.MessagePage{
		Message: string(bb),
	}, nil
}

func (s *PagemessageServer) GetAboutUs(ctx context.Context, req *emptypb.Empty) (*proto.MessagePage, error) {
	back1 := models.Message{}
	// 查询数据
	initialization.DB.Where(" `key` = ?", "关于我们").First(&back1)
	jsonback1, _ := json.Marshal(back1)
	return &proto.MessagePage{
		Message: string(jsonback1),
	}, nil
}

func (s *PagemessageServer) GetLearningStyle(ctx context.Context, req *emptypb.Empty) (*proto.MessagePage, error) {
	back1 := []models.Message{}
	// 查询数据
	initialization.DB.Where("keyform = ?", "学习方式").Find(&back1)
	//开始按照逻辑添加数据
	back2 := [4]models.Message{}
	for _, b := range back1 {
		if b.Key == "大一" {
			back2[0] = b
		}
		if b.Key == "大二" {
			back2[1] = b
		}
		if b.Key == "大三" {
			back2[2] = b
		}
		if b.Key == "大四" {
			back2[3] = b
		}

	}
	jsonback1, _ := json.Marshal(back2)
	return &proto.MessagePage{
		Message: string(jsonback1),
	}, nil
}

func (s *PagemessageServer) GetClubLocation(ctx context.Context, req *emptypb.Empty) (*proto.MessagePage, error) {
	back1 := models.Message{}
	// 查询数据
	initialization.DB.Where(" `key` = ?", "社团位置").First(&back1)
	jsonback1, _ := json.Marshal(back1)
	return &proto.MessagePage{
		Message: string(jsonback1),
	}, nil
}

func (s *PagemessageServer) Add_Training_Time(ctx context.Context, req *proto.ReviseMessage) (*emptypb.Empty, error) {
	//先查有多少数据
	var count int64
	initialization.DB.Model(&models.Message{}).Where("keyform = ?", "培训时间").Count(&count)
	count++
	//增加一个数据
	result := initialization.DB.Create(&models.Message{
		Key:     strconv.FormatInt(count, 10) + req.Key,
		Keyform: "培训时间",
		Value1:  req.Value1,
		Value2:  req.Value2,
	})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return &emptypb.Empty{}, result.Error
	}
	return &emptypb.Empty{}, nil
}

func (s *PagemessageServer) Del_Training_Time(ctx context.Context, req *proto.DelMessage) (*emptypb.Empty, error) {
	result := []models.Message{}
	for _, b := range req.Message {
		c := models.Message{
			Key:     b.Key,
			Keyform: b.Keyform,
			Value1:  b.Value1,
			Value2:  b.Value2,
		}
		result = append(result, c)
	}
	//先删除开始事务
	tx := initialization.DB.Begin() //开始事务
	//开始先删除数据
	result1 := tx.Where("keyform = ?", "培训时间").Delete(&models.Message{})
	if result1.Error != nil {
		tx.Rollback()
		log.SugarLogger.Error(result1.Error)
		return &emptypb.Empty{}, result1.Error
	}
	//开始增加数据
	result2 := tx.Create(&result)
	if result2.Error == nil {
		tx.Rollback()
		log.SugarLogger.Error(result1.Error)
		return &emptypb.Empty{}, result2.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *PagemessageServer) Add_Club_Direction(ctx context.Context, req *proto.ReviseMessage) (*emptypb.Empty, error) {
	//先查有多少数据
	var count int64
	initialization.DB.Model(&models.Message{}).Where("keyform = ?", "培训时间").Count(&count)
	count++
	//增加一个数据
	result := initialization.DB.Create(&models.Message{
		Key:     req.Key,
		Keyform: "主攻方向",
		Value1:  req.Value1,
		Value2:  req.Value2,
	})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return &emptypb.Empty{}, result.Error
	}
	return &emptypb.Empty{}, nil
}

func (s *PagemessageServer) Del_Club_Direction(ctx context.Context, req *proto.ReviseMessage) (*emptypb.Empty, error) {
	result := models.Message{
		Key: req.Key,
	}
	//先删除开始事务
	tx := initialization.DB.Begin() //开始事务
	//开始先删除数据
	result1 := tx.Where(" `key` = ? AND keyform = ?", result.Key, "主攻方向").Delete(&models.Message{})
	if result1.RowsAffected == 0 {
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, errors.New("没有找到数据")
	}
	if result1.Error != nil {
		tx.Rollback()
		log.SugarLogger.Error(result1.Error)
		return &emptypb.Empty{}, result1.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *PagemessageServer) Revise_Award_Information(ctx context.Context, req *proto.Prize) (*emptypb.Empty, error) {
	//开始事务
	tx := initialization.DB.Begin() //开始事务
	//开始修改数据
	//fmt.Println(req)
	result1 := tx.Table("prize").Where(" uuid = ? ", req.Uuid).Updates(map[string]interface{}{
		"name":   req.Name,
		"awards": req.Adards,
		"time":   req.Time,
	})
	if result1.RowsAffected == 0 {
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, errors.New("没有找到数据")
	}
	if result1.Error != nil {
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, result1.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *PagemessageServer) Del_Award_Information(ctx context.Context, req *proto.Prize) (*emptypb.Empty, error) {
	//先删除开始事务
	tx := initialization.DB.Begin() //开始事务
	//开始先删除数据
	result1 := tx.Where(" uuid = ? ", req.Uuid).Delete(&models.Prize{})
	if result1.RowsAffected == 0 {
		tx.Rollback()
		//如果有错误
		return &emptypb.Empty{}, errors.New("没有找到数据")
	}
	if result1.Error != nil {
		tx.Rollback()
		log.SugarLogger.Error(result1.Error)
		return &emptypb.Empty{}, result1.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *PagemessageServer) Add_Award_Information(ctx context.Context, req *proto.Prize) (*emptypb.Empty, error) {
	tx := initialization.DB.Begin() //开始事务
	//增加一个数据
	result := tx.Create(&models.Prize{
		Uuid:   req.Uuid,
		Name:   req.Name,
		Awards: req.Adards,
		Time:   req.Time,
	})
	if result.Error != nil {
		tx.Rollback()
		log.SugarLogger.Error(result.Error)
		return &emptypb.Empty{}, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
