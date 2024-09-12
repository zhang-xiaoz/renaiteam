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
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// 注意  在网络传输中，不能两个都为nil否者有问题
type BlogsServer struct {
	*proto.UnimplementedBlogsServer
}

func (s *BlogsServer) Add_Blog_Message_MysqlBlog(ctx context.Context, req *proto.BlogBlog) (*proto.BACKBlog, error) {
	fmt.Println(req)
	fmt.Println(req.CreatTime)
	back := proto.BACKBlog{}
	//转换一下时间格式
	t, _ := time.Parse("2006-01-02 15:04:05", req.CreatTime)
	blog_elastic := models.Blog_Elastic{
		Article_Uuid:   req.ArticleUuid,
		User_Uuid:      req.UserUuid,
		Title:          req.Title,
		Label:          req.Label,
		Abstract:       req.Abstract,
		Cover:          req.Cover,
		Article_Type:   req.ArticleType,
		Process:        req.Process,
		Visibility:     req.Visibility,
		Reading_Volume: int(req.Readingvolume),
		Creat_Time:     int(t.Unix()),
	}
	fmt.Println(req.CreatTime)
	//转换数据转到blog结构体里边
	blog := models.Blog{
		Article_Uuid:   req.ArticleUuid,
		User_Uuid:      req.UserUuid,
		Title:          req.Title,
		Content:        req.Content,
		Label:          req.Label,
		Cover:          req.Cover,
		Abstract:       req.Abstract,
		Visibility:     req.Visibility,
		Article_Type:   req.ArticleType,
		Process:        req.Process,
		Creat_Time:     req.CreatTime,
		Reading_Volume: int(req.Readingvolume),
	}
	fmt.Println(req.CreatTime)
	fmt.Println(blog)
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	result1 := tx.Create(&blog)     //直接填写进去
	if result1.Error != nil {
		log.SugarLogger.Error(result1.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result1.Error
	}
	//添加到elastic里边
	_, err := initialization.Elastic.Index().Index("renai_blog").Id(req.ArticleUuid).BodyJson(blog_elastic).Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(err)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result1.Error
	}
	tx.Commit() //提交事务
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Revise_Blog_Process_MysqlBlog(ctx context.Context, req *proto.BlogProcessBlog) (*proto.BACKBlog, error) {
	blog := models.Blog{}
	back := proto.BACKBlog{}
	//查找数据库数据看这篇文章是否属于这个用户//或者说是否有这篇博客

	result1, err := initialization.Elastic.Get().Index("renai_blog").Id(req.ArticleUuid).Do(context.Background())
	if err != nil { //没找到
		log.SugarLogger.Error(req.ArticleUuid)
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, errors.New("没有该数据")
	}
	//转一下数据
	result2 := models.Blog_Elastic{}
	err = json.Unmarshal(result1.Source, &result2)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	log.SugarLogger.Error(req.UserUuid)
	log.SugarLogger.Error(result2.Article_Uuid)
	if req.UserUuid != result2.User_Uuid { //说明不是同一人
		back.Back = false
		return &back, errors.New("作者不是同一人")
	}
	result2.Process = "已经审核"
	//修改数据//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&blog).Where("article_uuid=?", req.ArticleUuid).Updates(models.Blog{Process: req.Process})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	} else if result.RowsAffected == 0 {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, errors.New("没有找到这个数据")
	}
	//修改el里边数据
	_, err = initialization.Elastic.Update().Index("renai_blog").Id(req.ArticleUuid).Doc(map[string]interface{}{"process": req.Process}).Refresh("true").Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	tx.Commit() //成功提交事务
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Check_Blog_Uuid_MysqlBlog(ctx context.Context, req *proto.BlogProcessBlog) (*proto.BACKBlog, error) {
	back := proto.BACKBlog{}
	//查找数据库数据看这篇文章是否属于这个用户//或者说是否有这篇博客
	result1, err := initialization.Elastic.Get().Index("renai_blog").Id(req.ArticleUuid).Do(context.Background())
	if err != nil { //没找到
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, errors.New("没有该数据")
	}
	//转一下数据
	result2 := models.Blog_Elastic{}
	err = json.Unmarshal(result1.Source, &result2)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	if req.UserUuid != result2.User_Uuid { //说明不是同一人
		back.Back = false
		return &back, errors.New("作者不是同一人")
	}
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Get_Blog_Message_Elastic_TopBlog(ctx context.Context, req *emptypb.Empty) (*proto.BlogloadingBlog, error) {
	//初始化返回数据
	back1 := proto.BlogloadingBlog{}
	//c/c++的数据
	var sliceOfBlogs []models.Blog_Elastic
	a3 := &sliceOfBlogs
	//获取人工智能的数据
	//获取当前时间
	time1 := time.Now().Unix()
	a1, _, err := get_Blog_Message_Elastic_TopBlog1(int(time1), "人工智能")
	if err != nil {
		return nil, err
	}
	//获取c/c++的数据
	asas := int(time1)
	for {
		a2, c2, err := get_Blog_Message_Elastic_TopBlog1(asas, "c/c++")
		if err != nil {
			return nil, err
		}
		if c2 < 5 {
			//等会直接跳出就行了
			//查找是否有重复数据
			for _, i2 := range *a2 {
				for i3, i4 := range *a1 {
					if i2.Article_Uuid == i4.Article_Uuid {
						break
					}
					if i3 == len(*a1)-1 {
						*a3 = append(*a3, i2)
						if len(*a3) >= 4 {
							break
						}
					}
				}
			}
			break
		} else { //说明数据多着来
			//循环判断
			for _, i2 := range *a2 {
				asas = i2.Creat_Time
				for i3, i4 := range *a1 {
					if i2.Article_Uuid == i4.Article_Uuid {
						break
					}
					if i3 == len(*a1)-1 {
						*a3 = append(*a3, i2)
						if len(*a3) >= 4 {
							break
						}
					}
				}
				if len(*a3) >= 4 {
					break
				}
			}
			if len(*a3) >= 4 {
				break
			}
		}
	}
	//获取现在时间
	back1.B = time.Now().Unix()
	//现在开始获取真正需要的数据
	//创建两个需要返回的数据
	var a1_html []models.Blog_Elastic_loadstorage
	var a2_html []models.Blog_Elastic_loadstorage
	for _, ccc1 := range *a1 {
		a3_html := models.Blog_Elastic_loadstorage{
			Article_Uuid: ccc1.Article_Uuid,
		}
		a1_html = append(a1_html, a3_html)
	}
	for _, ccc2 := range *a3 {
		a3_html := models.Blog_Elastic_loadstorage{
			Article_Uuid: ccc2.Article_Uuid,
		}
		a1_html = append(a2_html, a3_html)
	}
	json1, err := json.Marshal(a1_html)
	if err != nil {
		return nil, err
	}
	json2, err := json.Marshal(a2_html)
	if err != nil {
		return nil, err
	}
	//已经获取到两个数据转json
	back1.A = string(json1)
	back1.C = string(json2)
	return &back1, nil
}

func (s *BlogsServer) Check_Blog_Uuid_Elastic_Blog(ctx context.Context, req *proto.BlogloadingBlog) (*proto.BACKBlog, error) {
	back := proto.BACKBlog{}
	//string类型的json转
	var blog1 []models.Blog_Elastic_loadstorage
	var blog2 []models.Blog_Elastic_loadstorage
	err := json.Unmarshal([]byte(req.A), &blog1)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	err = json.Unmarshal([]byte(req.C), &blog2)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	//如果两个都不为4则直接从新获取
	if len(blog1) < 4 {
		back.Back = false
		return &back, nil
	}
	if len(blog2) < 4 {
		back.Back = false
		return &back, nil
	}
	//开始检测是否有不能使用的
	if !check_Blog_Uuid_Elastic_Blog1(blog1) {
		back.Back = false
		return &back, nil
	}
	if !check_Blog_Uuid_Elastic_Blog1(blog2) {
		back.Back = false
		return &back, nil
	}
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Get_Blog_Message_ArticleUuid(ctx context.Context, req *proto.Blog_ArticleUuid) (*proto.Blog_Elastic_Message_TopBlog, error) {
	return get_Blog_Message_ArticleUuid1(req.ArticleUuid)
}

func (s *BlogsServer) Get_Blog_Message_Elastic_Lower(ctx context.Context, req *proto.BlogloadingBlog) (*proto.Blog_Elastic_Message_TopBlogs, error) {
	//创建返回所用结构体
	back := proto.Blog_Elastic_Message_TopBlogs{}
	var data []*proto.Blog_Elastic_Message_TopBlog
	back.Data = data
	//先把json数据处理好
	var ai []models.Blog_Elastic_loadstorage
	var c []models.Blog_Elastic_loadstorage
	json.Unmarshal([]byte(req.A), &ai)
	// err := json.Unmarshal([]byte(req.A), &ai)
	// if err != nil {
	// 	log.SugarLogger.Error(err)
	// 	return nil, err
	// }
	json.Unmarshal([]byte(req.C), &c)
	// err = json.Unmarshal([]byte(req.C), &c)
	// if err != nil {
	// 	log.SugarLogger.Error(err)
	// 	return nil, err
	// }
	//开始查找数据
	if req.B == 0 {
		req.B = time.Now().Unix()
	}
	for {
		result1, number1, _, err := get_Blog_Message_Elastic_Lower1(req.B)
		if err != nil {
			log.SugarLogger.Error(err)
			return nil, err
		}
		//查询是否重复
		fmt.Println("bbbbbb")
		fmt.Println(result1)
		fmt.Println(number1)
		for _, i1 := range *result1 {
			fmt.Println(i1.Creat_Time, i1.Article_Uuid)
			//开始查找是否有重复的
			i := 0
			for _, i2 := range ai {
				if i1.Article_Uuid == i2.Article_Uuid {
					i = 1
					break
				}
			}
			if i == 1 {
				i = 0
				continue
			}
			for _, i3 := range c {
				if i1.Article_Uuid == i3.Article_Uuid {
					i = 1
					break
				}
			}
			if i == 1 {
				i = 0
				continue
			}
			datadata := proto.Blog_Elastic_Message_TopBlog{
				ArticleUuid: i1.Article_Uuid,
				Title:       i1.Title,
				Cover:       i1.Cover,
				Abstract:    i1.Abstract,
				Visibility:  i1.Visibility,
			}
			req.B = int64(i1.Creat_Time)
			back.Data = append(back.Data, &datadata)
			if len(back.Data) >= 4 { //说明数据够了
				break
			}
		}
		fmt.Println(number1)
		if number1 < 4 {
			break
		}
		fmt.Println(len(back.Data))
		if len(back.Data) >= 4 { //说明数据够了
			break
		}
	}
	//返回数据
	back.Number = req.B
	return &back, nil
}

func (s *BlogsServer) Search_Blog_ElasticBlog(ctx context.Context, req *proto.AAAAA) (*proto.Blog_Elastic_Message_SerachBlogs, error) {
	back := proto.Blog_Elastic_Message_SerachBlogs{}
	var data []*proto.Blog_Elastic_Message_TopBlog
	back.Data = data
	//开始搜索数据
	//开始查找数据
	fmt.Println(req.B)
	fmt.Println(req.D)
	fmt.Println(req)
	if req.D == 0 {
		fmt.Println("aaaaaaaa")
		req.D = 10
	} else {
		req.D = 4
	}
	result1, number1, n1, n2, err := search_Blog_ElasticBlog1(req.A, req.B, int(req.C), int(req.D))
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, err
	}
	back.Number = int64(number1)
	//json化 n1,n2
	type identification struct {
		Uuid           string
		Reading_Volume int `json:"reading_volume"` //阅读量
	}
	message1 := identification{}
	message1.Reading_Volume = n1
	message1.Uuid = n2
	cc1, _ := json.Marshal(message1)
	back.Previous = string(cc1)
	for _, i2 := range *result1 {
		data2 := proto.Blog_Elastic_Message_TopBlog{
			ArticleUuid: i2.Article_Uuid,
			Title:       i2.Title,
			Cover:       i2.Cover,
			Abstract:    i2.Abstract,
			Visibility:  i2.Visibility,
		}
		back.Data = append(back.Data, &data2)
	}
	return &back, nil
}

func (s *BlogsServer) Add_Blog_Reading_Volume(ctx context.Context, req *proto.Blog_ArticleUuid) (*proto.BACKBlog, error) {
	back := proto.BACKBlog{}
	//先增加mysql的
	//开启事务
	tx := initialization.DB.Begin() //开始事务
	blog := models.Blog{}
	result := tx.Model(&blog).Where("article_uuid = ?", req.ArticleUuid).Update("reading_volume", gorm.Expr("reading_volume+?", 1))
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	} else if result.RowsAffected == 0 {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, errors.New("没有找到这个数据")
	}
	err := add_Blog_Reading_Volume1(req.ArticleUuid)
	if err != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	tx.Commit() //提交事务
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Get_Blog_Mysql_Message(ctx context.Context, req *proto.Blog_ArticleUuid) (*proto.BlogBlog, error) {
	//获取数据
	blog := models.Blog{}
	// 使用 First 方法根据姓名获取数据
	result := initialization.DB.Where("article_uuid = ?", req.ArticleUuid).First(&blog)
	// 检查查询结果
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("没有找到这个数据")
	}
	return &proto.BlogBlog{
		ArticleUuid:   blog.Article_Uuid,
		UserUuid:      blog.User_Uuid,
		Title:         blog.Title,
		Content:       blog.Content,
		Label:         blog.Label,
		Cover:         blog.Cover,
		Abstract:      blog.Abstract,
		Visibility:    blog.Visibility,
		ArticleType:   blog.Article_Type,
		Process:       blog.Process,
		Readingvolume: int64(blog.Reading_Volume),
		CreatTime:     blog.Creat_Time,
	}, nil
}

func (s *BlogsServer) Save_Redis_StorageBlog(ctx context.Context, req *proto.Redis_Storage_Blog) (*proto.Redis_Back_Picture, error) {
	//先进行连接
	back := proto.Redis_Back_Picture{}
	back.Uuid = ""
	rediss := initialization.Redis.Get()
	if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		return &back, err
	}
	defer rediss.Close() //关闭数据库连接
	//更改数据库
	_, err := rediss.Do("select", req.Who)
	if err != nil {
		log.SugarLogger.Error(err)
		return &back, err
	}
	//获取时间
	req.Value = time.Now().Unix()
	number := float64(req.Value)
	//进行存储
	_, err = rediss.Do("ZADD", "picture", number, req.Key)
	if err != nil {
		log.SugarLogger.Error(err)
		return &back, err
	}
	//进行删除判断
	//6天不用自动删除
	type picture struct {
		p string
	}
	back1 := []picture{}
	//定义分数范围
	minValue := 0.1                         // 分数大于此值
	maxValue := float64(req.Value) - 432000 // 分数小于或等于正无穷大
	result, err := redis.Values(rediss.Do("ZRANGE", "picture", minValue, maxValue, "WITHSCORES"))
	ccc := 0
	if err == nil {
		for _, v := range result {
			back2 := picture{}
			if ccc == 0 {
				ccc++
				if data, ok := v.([]byte); ok {
					back2.p = string(data)
					back1 = append(back1, back2)
				}
			} else {
				ccc--
			}
		}
		//删除redis里边那么多数据
		rediss.Do("ZREM", "picture", redis.Args{}.AddFlat(back1).Add("UNLINK"))
		//直接json化保存
		c, err := json.Marshal(back1)
		if err == nil {
			back.Uuid = string(c)
		}
	}
	return &back, nil
}

func (s *BlogsServer) Delete_El6_MysqlBlog(ctx context.Context, req *proto.Redis_Storage_Blog) (*proto.BACKBlog, error) {
	//先进行连接
	back := proto.BACKBlog{}
	back.Back = false
	rediss := initialization.Redis.Get()
	if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		return &back, err
	}
	rediss1 := initialization.Redis.Get()
	if _, err := rediss1.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		return &back, err
	}
	defer rediss1.Close() //关闭数据库连接
	//更改数据库
	_, err := rediss.Do("select", req.Who)
	if err != nil {
		log.SugarLogger.Error(err)
		return &back, err
	}
	_, err = rediss1.Do("select", req.Who)
	if err != nil {
		log.SugarLogger.Error(err)
		return &back, err
	}
	wg := sync.WaitGroup{} //开始控制线程
	wg.Add(2)
	var err1 error
	//更改el6状态
	//修改el里边数据
	_, err = initialization.Elastic.Update().Index("renai_blog").Id(req.Key).Doc(map[string]interface{}{"process": "已经删除"}).Refresh("true").Do(context.Background())
	go func() { //进行存储
		//获取时间
		ccdc := time.Now().Unix()
		number := float64(ccdc)
		//进行存储
		_, err1 = rediss1.Do("ZADD", "blog", number, req.Key)
		if err != nil {
			log.SugarLogger.Error(err)
		}
		wg.Done()
	}()
	go func() { //进行删除
		//6天不用自动删除
		type blog struct {
			p string
		}
		//定义分数范围
		minValue := 0.1                                 // 分数大于此值
		maxValue := float64(time.Now().Unix()) - 604800 // 分数小于或等于正无穷大
		result, err2 := redis.Values(rediss.Do("ZREVRANGEBYSCORE", "blog", maxValue, minValue, "WITHSCORES"))
		ccc := 0
		if err2 == nil {
			for _, v := range result {
				back2 := blog{}
				if ccc == 0 {
					ccc++
					if data, ok := v.([]byte); ok {
						back2.p = string(data)
						fmt.Println(back2.p)
						//删除数据mysql与el6
						delete_El6_Mysql_blog1(back2.p)
						rediss.Do("ZREM", "blog", back2.p)
					}
				} else {
					ccc--
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()
	if err1 != nil {
		log.SugarLogger.Error(err1)
		return &back, err1
	}
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Completely_El6_MysqlBlog(ctx context.Context, req *proto.Blog_ArticleUuid) (*proto.BACKBlog, error) {
	back := &proto.BACKBlog{}
	back.Back = false
	err := delete_El6_Mysql_blog1(req.ArticleUuid)
	if err != nil {
		return back, err
	}
	//删除redis数据库里边数据
	back.Back = true
	rediss := initialization.Redis.Get()
	if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
		log.SugarLogger.Error("redis连接失败")
		return back, nil
	}
	defer rediss.Close() //关闭数据库连接
	//更改数据库
	_, err = rediss.Do("select", "1")
	if err != nil {
		log.SugarLogger.Error(err)
		return back, nil
	}
	rediss.Do("ZREM", "blog", req.ArticleUuid)
	return back, nil
}

func (s *BlogsServer) Get_Delete_Blog_Mysql(ctx context.Context, req *proto.AAAA) (*proto.Blog_Elastic_Message_SerachBlogs, error) {
	back := proto.Blog_Elastic_Message_SerachBlogs{}
	var data []*proto.Blog_Elastic_Message_TopBlog
	back.Data = data
	//开始查找数据
	if req.A5 == 0 {
		req.A5 = 10
	} else {
		req.A5 = 4
	}
	result1, n2, n1, number1, err := get_Delete_Blog_Mysql1(req.A4, req.A5, req.A1, req.A2, req.A3)
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, err
	}
	back.Number = int64(number1)
	//json化 n1,n2
	type identification struct {
		Time int
		Uuid string
	}
	message1 := identification{}
	message1.Time = n1
	message1.Uuid = n2
	cc1, err := json.Marshal(message1)
	if err != nil {
		return nil, err
	}
	back.Previous = string(cc1)
	for _, i2 := range *result1 {
		data2 := proto.Blog_Elastic_Message_TopBlog{
			ArticleUuid: i2.Article_Uuid,
			Title:       i2.Title,
			Cover:       i2.Cover,
			Abstract:    i2.Abstract,
			Visibility:  i2.Visibility,
		}
		back.Data = append(back.Data, &data2)
	}
	return &back, nil
}

func (s *BlogsServer) Revise_Blog_Process_PlusBlog(ctx context.Context, req *proto.BlogProcessPlusBlog) (*proto.BACKBlog, error) {
	//先更改状态
	blog := models.Blog{}
	back := proto.BACKBlog{}
	//查找数据库数据看这篇文章是否属于这个用户//或者说是否有这篇博客
	result1, err := initialization.Elastic.Get().Index("renai_blog").Id(req.ArticleUuid).Do(context.Background())
	if err != nil { //没找到
		log.SugarLogger.Error(req.ArticleUuid)
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, errors.New("没有该数据")
	}
	//转一下数据
	result2 := models.Blog_Elastic{}
	err = json.Unmarshal(result1.Source, &result2)
	if err != nil {
		log.SugarLogger.Error(err)
		back.Back = false
		return &back, err
	}
	if req.UserUuid != result2.User_Uuid { //说明不是同一人
		back.Back = false
		return &back, errors.New("作者不是同一人")
	}
	result2.Process = req.Process
	//修改数据//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&blog).Where("article_uuid=?", req.ArticleUuid).Updates(models.Blog{Process: req.Process})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	} else if result.RowsAffected == 0 {
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, errors.New("没有找到这个数据")
	}
	//修改el里边数据
	_, err = initialization.Elastic.Update().Index("renai_blog").Id(req.ArticleUuid).Doc(map[string]interface{}{"process": req.Process}).Refresh("true").Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	tx.Commit() //成功提交事务
	if req.Code == 1 {
		//删除redis里边的数据
		rediss := initialization.Redis.Get()
		if _, err := rediss.Do("AUTH", "wasd2002"); err != nil { //密码认证
			log.SugarLogger.Error("redis连接失败")
			return &back, nil
		}
		defer rediss.Close() //关闭数据库连接
		//更改数据库
		rediss.Do("select", "1")
		rediss.Do("ZREM", "blog", req.ArticleUuid)
	}
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Add_Mysql_Blog_Collection(ctx context.Context, req *proto.Blog_UserUuid) (*proto.BACKBlog, error) {
	fmt.Println(req.Status)
	back := proto.BACKBlog{}
	back.Back = false
	//开启两个线程
	wg := sync.WaitGroup{} //开始控制线程
	wg.Add(2)
	var err1 error
	var err2 error
	go func() { //查找是否有这个博客
		blog := models.Blog{}
		result := initialization.DB.Where("article_uuid=?", req.Status).First(&blog)
		if result.Error != nil {
			err1 = result.Error
			wg.Done()
			return
		} else if result.RowsAffected == 0 {
			err1 = errors.New("没有找到数据")
			wg.Done()
			return
		}
		if blog.User_Uuid == req.UserUuid {
			err1 = errors.New("作者是同一个")
			wg.Done()
			return
		}
		wg.Done()
	}()
	go func() { //查找是否已经存储过
		result := initialization.DB.Where("article_uuid=?", req.Status).Where("user_uuid=?", req.UserUuid).First(&models.Blog_Collection{})
		if result.RowsAffected != 0 {
			err2 = errors.New("找到数据")
			wg.Done()
			return
		}
		wg.Done()
	}()
	wg.Wait()
	if err1 != nil {
		return &back, err1
	}
	if err2 != nil {
		return &back, err2
	}
	//进行保存操作
	collection := models.Blog_Collection{}
	collection.Article_Uuid = req.Status
	collection.User_Uuid = req.UserUuid
	collection.Time = time.Now().Format("2006-01-02 15:04:05")
	result := initialization.DB.Create(&collection)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return &back, result.Error
	}
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Delete_Mysql_Blog_Collection(ctx context.Context, req *proto.Blog_UserUuid) (*proto.BACKBlog, error) {
	back := proto.BACKBlog{}
	back.Back = false
	//进行删除操作
	collection := models.Blog_Collection{}
	collection.Article_Uuid = req.Status
	collection.User_Uuid = req.UserUuid
	//db.Where("id = ?", 1).Delete(&YourModel{})
	result := initialization.DB.Where("article_uuid=? and user_uuid=?", collection.Article_Uuid, collection.User_Uuid).Delete(&models.Blog_Collection{})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return &back, result.Error
	}
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Get_Mysql_Blog_Collection(ctx context.Context, req *proto.UserUuid_Blog) (*proto.Blog_Elastic_Message_CollectionBlogs, error) {
	back := proto.Blog_Elastic_Message_CollectionBlogs{}
	back1 := []*proto.Blog_Elastic_Message_TopBlog{}
	back.Data = back1
	//先判断是第几次开始查找
	fmt.Println(req)
	var code int
	code = 4
	if req.Time == 0 {
		req.Time = time.Now().Unix()
		code = 10
	}
	//开始查找数据
	blog_collection := []models.Blog_Collection{}
	//先改变时间
	time1 := time.Unix(req.Time, 0)
	time2 := time1.Format("2006-01-02 15:04:05")
	//构建 查询体
	result := initialization.DB.Where("time<?", time2).Where("user_uuid=?", req.UserUuid).Order("time DESC").Limit(code).Find(&blog_collection)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	fmt.Println(blog_collection)
	for _, i2 := range blog_collection {
		i4 := models.Blog{}
		//开始查找数据
		result1 := initialization.DB.Where("article_uuid=?", i2.Article_Uuid).
			Where("visibility IN (?)", []string{"全体可见", "仁爱成员可见"}).
			Where("process=?", "已经审核").First(&i4)
		if result1.Error != nil {
			log.SugarLogger.Error(result1.Error)
			continue
		} else {
			if i4.Article_Uuid == "" {
				continue
			} else {
				back.Number++
				// 定义日期时间的布局格式，RFC3339 是一种常用的日期时间格式
				const layout = time.RFC3339
				// 解析字符串为 time.Time 对象
				loc, _ := time.LoadLocation("Asia/Shanghai")
				t, _ := time.ParseInLocation(layout, i2.Time, loc)
				// 转换为 UNIX 时间戳（秒）
				back.Previous = t.Unix()
				i3 := proto.Blog_Elastic_Message_TopBlog{
					ArticleUuid: i4.Article_Uuid,
					Title:       i4.Title,
					Cover:       i4.Cover,
					Abstract:    i4.Abstract,
					Visibility:  i4.Visibility,
				}
				back1 = append(back1, &i3)
			}
		}
	}
	//开始返回数据
	back.Data = back1
	fmt.Println(&back)
	return &back, nil
}

func (s *BlogsServer) Revise_Blog_Elastic_Mysql(ctx context.Context, req *proto.BlogBlog) (*proto.BACKBlog, error) {
	//开始修改数据
	//先更改状态
	blog := models.Blog{}
	back := proto.BACKBlog{}
	//修改数据//开启事务
	tx := initialization.DB.Begin() //开始事务
	result := tx.Model(&blog).Where("article_uuid=?", req.ArticleUuid).Updates(models.Blog{
		Title:        req.Title,
		Content:      req.Content,
		Label:        req.Label,
		Cover:        req.Cover,
		Abstract:     req.Abstract,
		Visibility:   req.Visibility,
		Article_Type: req.ArticleType,
		Process:      req.Process,
		Creat_Time:   req.CreatTime,
	})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	} else if result.RowsAffected == 0 {
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, errors.New("没有找到这个数据")
	}
	//修改el里边数据
	_, err := initialization.Elastic.Update().Index("renai_blog").Id(req.ArticleUuid).Doc(map[string]interface{}{"title": req.Title, "label": req.Label, "abstract": req.Abstract, "cover": req.Cover, "article_type": req.ArticleType, "visibility": req.Visibility, "process": req.Process, "creat_time": req.CreatTime}).Refresh("true").Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		back.Back = false
		return &back, result.Error
	}
	tx.Commit() //成功提交事务
	back.Back = true
	return &back, nil
}

func (s *BlogsServer) Search_Blog_Label_Elastic(ctx context.Context, req *proto.AAAAA) (*proto.Blog_Elastic_Message_SerachBlogs, error) {
	back := proto.Blog_Elastic_Message_SerachBlogs{}
	var data []*proto.Blog_Elastic_Message_TopBlog
	back.Data = data
	//开始搜索数据
	//开始查找数据
	fmt.Println(req.B)
	fmt.Println(req.D)
	fmt.Println(req)
	if req.D == 0 {
		req.D = 10
	} else {
		req.D = 4
	}
	result1, number1, n1, n2, err := search_Blog_Label_Elastic1(req.A, req.B, int(req.C), int(req.D))
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, err
	}
	back.Number = int64(number1)
	//json化 n1,n2
	type identification struct {
		Uuid string
		Time int
	}
	message1 := identification{}
	message1.Time = n1
	message1.Uuid = n2
	cc1, _ := json.Marshal(message1)
	back.Previous = string(cc1)
	for _, i2 := range *result1 {
		data2 := proto.Blog_Elastic_Message_TopBlog{
			ArticleUuid: i2.Article_Uuid,
			Title:       i2.Title,
			Cover:       i2.Cover,
			Abstract:    i2.Abstract,
			Visibility:  i2.Visibility,
		}
		back.Data = append(back.Data, &data2)
	}
	return &back, nil
}

func (s *BlogsServer) Get_Blog_Collection_Status(ctx context.Context, req *proto.Blog_UserUuid) (*proto.Status, error) {
	blog := models.Blog{}
	back := proto.Status{}
	//先判断文章作者是否是他
	//构建 查询体
	result1 := initialization.DB.Where("article_uuid=?", req.Status).First(&blog)
	if result1.Error != nil {
		log.SugarLogger.Error(result1.Error)
		back.Status = 4
		return &back, result1.Error
	}
	if blog.User_Uuid == req.UserUuid {
		back.Status = 1
		return &back, result1.Error
	}
	//开始查找该用户是否收藏
	result2 := initialization.DB.Where("article_uuid=?", req.Status).Where("user_uuid=?", req.UserUuid).First(&models.Blog_Collection{})
	if result2.RowsAffected != 0 {
		//找到数据
		back.Status = 2
	} else {
		back.Status = 3
	}
	return &back, result1.Error
}

func (s *BlogsServer) Get_Blog_Manage_All(ctx context.Context, req *proto.Blog_Page) (*proto.Blog_Message, error) {
	//根据时间倒序开始获取数据内容
	back1 := []models.Blog_Background_Html{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("blog").Where("process=?", req.Message).Order("creat_time DESC").Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.Blog{}).Where("process=?", req.Message).Count(&count)
	}
	//开始循环改变数据
	for _, v := range back1 {
		user := models.User{}
		//开始改名字
		initialization.DB.Where("uuid = ?", v.User_Uuid).First(&user)
		v.User_Uuid = user.Name
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
	return &proto.Blog_Message{
		Message: string(bb),
	}, nil
}

func (s *BlogsServer) Delete_El6_Mysql_Process(ctx context.Context, req *proto.BlogProcessBlog) (*proto.BACKBlog, error) {
	blog := models.Blog{}
	//开启事务进行修改
	tx := initialization.DB.Begin()
	result := tx.Model(&blog).Where("article_uuid=?", req.ArticleUuid).Updates(models.Blog{Process: req.Process})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return &proto.BACKBlog{
			Back: false,
		}, errors.New("未找到")
	}
	if result.Error != nil {
		tx.Rollback()
		return &proto.BACKBlog{
			Back: false,
		}, result.Error
	}
	//开始修改elastic里边内容
	_, err := initialization.Elastic.Update().Index("renai_blog").Id(req.ArticleUuid).Doc(map[string]interface{}{"process": req.Process}).Refresh("true").Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		return &proto.BACKBlog{
			Back: false,
		}, err
	}
	tx.Commit() //成功提交事务
	return &proto.BACKBlog{
		Back: false,
	}, nil
}

func (s *BlogsServer) Get_Blog_Manage_Name(ctx context.Context, req *proto.Blog_Page) (*proto.Blog_Message, error) {
	//根据时间倒序开始获取数据内容
	back1 := []models.Blog_Background_Html{}
	user := []models.User{}
	result2 := initialization.DB.Where("name like ?", "%"+req.Message+"%").First(&user)
	if result2.RowsAffected == 0 {
		return nil, errors.New("没有找到")
	}
	if result2.Error != nil {
		return nil, result2.Error
	}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 假设你有一个包含多个搜索词的切片
	searchTerms := []string{}
	for vvv1, vvv2 := range user {
		searchTerms[vvv1] = vvv2.Uuid
	}

	// 构造一个包含所有LIKE条件的字符串，使用OR连接
	var conditions []string
	for _, term := range searchTerms {
		conditions = append(conditions, fmt.Sprintf("title LIKE '%%%s%%'", term))
	}
	fullCondition := strings.Join(conditions, " OR ")
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("blog").Where(fullCondition).Where("process = ?", "已经审核").Order("creat_time DESC").Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.Blog{}).Where(fullCondition).Where("process = ?", "已经审核").Count(&count)
	}
	//开始循环改变数据
	for i := range back1 {
		back1[i].User_Uuid = req.Message
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
	return &proto.Blog_Message{
		Message: string(bb),
	}, nil
}

func (s *BlogsServer) Get_Blog_Manage_Title(ctx context.Context, req *proto.Blog_Page) (*proto.Blog_Message, error) {
	//根据时间倒序开始获取数据内容
	back1 := []models.Blog_Background_Html{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("blog").Where("label like ?", "%"+req.Message+"%").Where("process = ?", "已经审核").Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.Blog{}).Where("label like ?", "%"+req.Message+"%").Where("process = ?", "已经审核").Count(&count)
	}
	//开始循环改变数据
	for _, v := range back1 {
		user := models.User{}
		//开始改名字
		initialization.DB.Where("uuid = ?", v.User_Uuid).First(&user)
		v.User_Uuid = user.Name
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
	return &proto.Blog_Message{
		Message: string(bb),
	}, nil
}

func (s *BlogsServer) Get_Blog_Manage_Label(ctx context.Context, req *proto.Blog_Page) (*proto.Blog_Message, error) {
	//根据时间倒序开始获取数据内容
	back1 := []models.Blog_Background_Html{}
	// 计算偏移量
	offset := (req.P - 1) * req.Pn
	// 使用Limit和Offset进行分页查询
	result := initialization.DB.Table("blog").Where("label like ?", "%"+req.Message+"%").Where("process = ?", "已经审核").Limit(int(req.Pn)).Offset(int(offset)).Find(&back1)
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		return nil, result.Error
	}
	var count int64
	if req.P == 1 { //查找一共有多少数据
		initialization.DB.Model(&models.Blog{}).Where("label like ?", "%"+req.Message+"%").Where("process = ?", "已经审核").Count(&count)
	}
	//开始循环改变数据
	for _, v := range back1 {
		user := models.User{}
		//开始改名字
		initialization.DB.Where("uuid = ?", v.User_Uuid).First(&user)
		v.User_Uuid = user.Name
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
	return &proto.Blog_Message{
		Message: string(bb),
	}, nil
}
