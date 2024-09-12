package api

import (
	"context"
	"encoding/json"
	"os"
	"time"
	"web/log"
	"web/models"
	modelss "web/models/blogs"

	"web/proto_blog"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/google/uuid"
	goredislib "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Add_User_Blog(ctx *gin.Context) {
	//首先判断get里边是否有数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//实体化对象获取数据
	blog := modelss.BLOG{}
	err := ctx.ShouldBind(&blog)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	//判断标题是否符合条件
	i := 0
	for range blog.Title {
		i++
	}
	if i < 5 || i > 20 {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "标题不符合条件",
			"jwt":     jwtString,
		})
		return
	}
	//判断摘要是否小于100个字
	i = 0
	for range blog.Abstract {
		i++
	}
	if i > 100 {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "摘要大于100字",
			"jwt":     jwtString,
		})
		return
	}
	//判断文章类型
	if blog.Article_Type != "原创" && blog.Article_Type != "转载" && blog.Article_Type != "翻译" {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "文章类型不符合",
			"jwt":     jwtString,
		})
		return
	}
	//判断审核是否正确
	//判断文章类型
	if blog.Process != "未审核" && blog.Process != "草稿" {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "文章存储格式不正确",
			"jwt":     jwtString,
		})
		return
	}
	//获取时间数据
	blog.Creat_Time = time.Now().Format("2006-01-02 15:04:05")
	//开始判断到底是有存还是没有从新获取
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	//加锁
	client := goredislib.NewClient(&goredislib.Options{ //连接redis服务器配置//redis包内容
		Addr:     models.Overall_Situation_Redis.Redis.Port,
		Password: models.Overall_Situation_Redis.Redis.Password,
	})
	pool := goredis.NewPool(client)                                                               //用redsync把client封装起来
	rs := redsync.New(pool)                                                                       //redsync.New方法从给定的Redis连接池创建并返回一个新的Redsync实例。
	mutex := rs.NewMutex(models.Overall_Situation_Redisclock.RedisClock.Blog + blog.Article_Uuid) //创建一个redis锁//里边包含锁的一些配置//默认过期时间八秒
	if err := mutex.Lock(); err != nil {                                                          //获取锁
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	quit := make(chan bool)
	go func() { //go线程
		for {
			select {
			case <-quit:
				return
			default:
				time.Sleep(time.Second * 4)
				mutex.Extend() //启动一个协程  重置到锁之前状态 保证一直运行
				log.SugarLogger.Info("注册时redis锁重置8秒")
			}
		}
	}()
	defer mutex.Unlock() //解锁//defer是反着执行的
	defer close(quit)
	if blog.Article_Uuid != "" {
		//那就是查找并且修改信息
		_, err = userBlogsClient.Check_Blog_Uuid_MysqlBlog(context.Background(), &proto_blog.BlogProcessBlog{
			ArticleUuid: blog.Article_Uuid,
			UserUuid:    blog.User_Uuid,
		})
		if err != nil {
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "添加失败",
				"jwt":     jwtString,
			})
			return
		}
		//如果有则进行更改存储
		_, err = userBlogsClient.Revise_Blog_Elastic_Mysql(context.Background(), &proto_blog.BlogBlog{
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
			Readingvolume: 0,
			CreatTime:     time.Now().Format("2006-01-02 15:04:05"), //还是把时间给他吧2012-12-31 11:30:45这样格式
		})
		if err != nil {
			log.SugarLogger.Error(err)
			ctx.JSON(200, map[string]string{
				"code":    "205",
				"message": "修改失败",
				"jwt":     jwtString,
			})
			return
		}
		ctx.JSON(200, map[string]string{
			"code":    "206",
			"message": "修改成功",
			"jwt":     jwtString,
		})
		return
	} else {
		blog.Article_Uuid = uuid.New().String()
		_, err = userBlogsClient.Add_Blog_Message_MysqlBlog(context.Background(), &proto_blog.BlogBlog{
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
			Readingvolume: 0,
			CreatTime:     time.Now().Format("2006-01-02 15:04:05"), //还是把时间给他吧2012-12-31 11:30:45这样格式
		})
		if err != nil {
			log.SugarLogger.Error(err)
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "添加失败",
				"jwt":     jwtString,
			})
			return
		}
		//返回信息
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "添加成功",
			"jwt":     jwtString,
		})
		return
	}
}
func Agree_User_Blog(ctx *gin.Context) { //同意某个文章审核
	//首先判断get里边是否有数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	blog := modelss.Blog_Process{}
	err := ctx.ShouldBind(&blog)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "未成功",
			"jwt":     jwtString,
		})
		return
	}
	if blog.Process != "已经审核" {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "未成功",
			"jwt":     jwtString,
		})
		return
	}
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	//判断这篇文章是否是这个用户的
	result1, err := userBlogsClient.Check_Blog_Uuid_MysqlBlog(context.Background(), &proto_blog.BlogProcessBlog{
		ArticleUuid: blog.Article_Uuid,
		UserUuid:    blog.User_Uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	if !result1.Back {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "该用户没有该文章",
			"jwt":     jwtString,
		})
		return
	}
	//加锁//加添加的锁
	client := goredislib.NewClient(&goredislib.Options{ //连接redis服务器配置//redis包内容
		Addr:     models.Overall_Situation_Redis.Redis.Port,
		Password: models.Overall_Situation_Redis.Redis.Password,
	})
	pool := goredis.NewPool(client)                                                               //用redsync把client封装起来
	rs := redsync.New(pool)                                                                       //redsync.New方法从给定的Redis连接池创建并返回一个新的Redsync实例。
	mutex := rs.NewMutex(models.Overall_Situation_Redisclock.RedisClock.Blog + blog.Article_Uuid) //创建一个redis锁//里边包含锁的一些配置//默认过期时间八秒
	if err := mutex.Lock(); err != nil {                                                          //获取锁
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	quit := make(chan bool)
	go func() { //go线程
		for {
			select {
			case <-quit:
				return
			default:
				time.Sleep(time.Second * 4)
				mutex.Extend() //启动一个协程  重置到锁之前状态 保证一直运行
				log.SugarLogger.Info("注册时redis锁重置8秒")
			}
		}
	}()
	defer mutex.Unlock() //解锁//defer是反着执行的
	defer close(quit)
	_, err = userBlogsClient.Revise_Blog_Process_MysqlBlog(context.Background(), &proto_blog.BlogProcessBlog{
		UserUuid:    blog.User_Uuid,
		Process:     blog.Process,
		ArticleUuid: blog.Article_Uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "204",
		"message": "添加成功",
		"jwt":     jwtString,
	})
}

func Add_Blog_Picture(ctx *gin.Context) {
	//同意某个图片
	//首先判断get里边是否有数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	file, error := ctx.FormFile("file")
	if error != nil {
		log.SugarLogger.Error(error)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "上传失败",
			"jwt":     jwtString,
		})
		return
	}
	if file == nil {
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "上传失败",
			"jwt":     jwtString,
		})
		return
	}
	file_uuid := uuid.New().String()
	if file.Size > 2*1024*1024 {
		//返回信息
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "图片过大",
			"jwt":     jwtString,
		})
		return
	}
	contentType := file.Header.Get("Content-Type")
	attribute := ""
	if contentType == "image/jpeg" {
		attribute = ".jpg"
	} else if contentType == "image/png" {
		attribute = ".png"
	} else {
		//返回信息
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "图片格式不对",
			"jwt":     jwtString,
		})
		return
	}
	//上传到redis里边
	//连接grpc服务
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "添加失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)

	result1, err := userBlogsClient.Save_Redis_StorageBlog(context.Background(), &proto_blog.Redis_Storage_Blog{
		Who: "1",
		Key: "file_uuid" + attribute,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "上传失败",
			"jwt":     jwtString,
		})
		return
	}
	type picturecc struct {
		p string
	}
	//删除临时的照片
	if result1 != nil && result1.Uuid != "" {
		back1 := []picturecc{}
		json.Unmarshal([]byte(result1.Uuid), &back1)
		//删除里边照片
		for _, vv := range back1 {
			os.Remove("./img/temporary/" + vv.p)
		}

	}
	//保存到redis里边并且保存保留时间//这里每次轮询一边查看是否有不用的//得用一个全新的redis数据库
	//保存图片
	err = ctx.SaveUploadedFile(file, "./img/temporary/"+file_uuid+attribute)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "上传失败",
			"jwt":     jwtString,
		})
		return
	}
	//上传到redis里边 key是uuid value是时间
	ctx.JSON(200, map[string]string{
		"code":    "205",
		"message": "上传成功",
		"uuid":    file_uuid + attribute,
		"jwt":     jwtString,
	})
}

func Get_Blog_Picture(ctx *gin.Context) {
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	// 获取文件名
	filename := ctx.Param("filename")
	// 图片所在的文件夹路径
	dir := "./img/blog/"

	// 完整的图片路径
	filePath := dir + filename

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
			"jwt":     jwtString,
		})
		return
	}
	// 使用Gin的c.File方法发送文件
	ctx.File(filePath)
}

func Delete_User_Blog(ctx *gin.Context) { //删除一个数据自动放入回收站//七天保存时间
	//首先先获取数据
	//首先判断get里边是否有数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	blog := modelss.Blog_Delete{}
	err := ctx.ShouldBind(&blog)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
			"jwt":     jwtString,
		})
		return
	}
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	userBlogsClient.Delete_El6_MysqlBlog(context.Background(), &proto_blog.Redis_Storage_Blog{
		Who: "1",
		Key: blog.Article_Uuid,
	})
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "删除成功",
		"jwt":     jwtString,
	})
}

func Completely_Delete_Blog(ctx *gin.Context) { //彻底删除一个数据，回收站
	//删除mysql和el6里边的还有redis里边的
	//获取数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	blog := modelss.Blog_Delete{}
	err := ctx.ShouldBind(&blog)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
			"jwt":     jwtString,
		})
		return
	}
	//连接grpc服务进行删除
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	_, err = userBlogsClient.Completely_El6_MysqlBlog(context.Background(), &proto_blog.Blog_ArticleUuid{
		ArticleUuid: blog.Article_Uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "删除成功",
		"jwt":     jwtString,
	})
}

func Get_Delete_Blog(ctx *gin.Context) { //查看删除/草稿等图片的博客

	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取用户的uuid
	uuid := ctx.GetHeader("uuid")
	if uuid == "" {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
			"jwt":     jwtString,
		})
		return
	}
	//获取上一个数据
	message := modelss.Blog_Search{}
	//解析数据
	err := ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
			"jwt":     jwtString,
		})
		return
	}
	//解析json数据
	type identification struct {
		Time int
		Uuid string
	}
	message1 := identification{}
	err = json.Unmarshal([]byte(message.Previous), &message1)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
			"jwt":     jwtString,
		})
		return
	}
	if message.Code == 0 {
		message1.Time = int(time.Now().Unix())
	} else if message.Code != 1 {
		ctx.JSON(200, map[string]string{
			"code":    "205",
			"message": "code数据不正确",
			"jwt":     jwtString,
		})
		return
	}
	//开启grpc服务
	//连接grpc服务进行删除
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	result1, err := userBlogsClient.Get_Delete_Blog_Mysql(context.Background(), &proto_blog.AAAA{
		A1: uuid,                 //用户uuid
		A2: message.Message,      //搜索信息 按照状态//这里是已经删除
		A3: message1.Uuid,        //上一个搜索的uuid
		A4: int64(message1.Time), //上一个搜索的时间
		A5: int64(message.Code),  //第几次发送
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
			"jwt":     jwtString,
		})
		return
	}
	if result1.Number == 0 {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "获取数据为0",
			"jwt":     jwtString,
		})
	} else {
		datadata, err := json.Marshal(result1.Data)
		if err != nil {
			log.SugarLogger.Error(err)
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "获取失败",
				"jwt":     jwtString,
			})
			return
		}
		ctx.JSON(200, map[string]string{
			"code":     "204",
			"message":  "获取成功",
			"jwt":      jwtString,
			"previous": result1.Previous,
			"data":     string(datadata),
		})
	}
}

func Garbage_Move_Draft(ctx *gin.Context) { //回收站博客存放到草稿箱//也就是更改状态
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取用户的uuid
	uuid := ctx.GetHeader("uuid")
	if uuid == "" {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "保存失败",
			"jwt":     jwtString,
		})
		return
	}
	//获取上一个数据
	message := modelss.Blog_Elastic_loadstorage{}
	//解析数据
	err := ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "保存失败",
			"jwt":     jwtString,
		})
		return
	}
	//开启grpc服务
	//连接grpc服务进行删除
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "保存失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	result, err := userBlogsClient.Revise_Blog_Process_PlusBlog(context.Background(), &proto_blog.BlogProcessPlusBlog{
		Process:     "草稿",
		UserUuid:    uuid,
		ArticleUuid: message.Article_Uuid,
		Code:        1,
	})
	if !result.Back {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "保存失败",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "保存成功",
		"jwt":     jwtString,
	})
}

func Get_Collection_Blog(ctx *gin.Context) { //查看收藏夹里边内容
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取接受到的时间信息  就是我跟他返回的数据
	//如果是第一次加载 //直接传入0即可
	message := modelss.Blog_Time{}
	//解析数据
	err := ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "加载失败",
			"jwt":     jwtString,
		})
		return
	}
	//获取用户的uuid
	uuid := ctx.GetHeader("uuid")
	//连接grpc服务进行获取数据
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "加载失败",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	result, err := userBlogsClient.Get_Mysql_Blog_Collection(context.Background(), &proto_blog.UserUuid_Blog{
		UserUuid: uuid,
		Time:     int64(message.Time),
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "加载失败",
			"jwt":     jwtString,
		})
		return
	}
	//转换一下数据
	backhtml := []modelss.Blog_Elastic_Html{}
	for _, i2 := range result.Data {
		i3 := modelss.Blog_Elastic_Html{
			Article_Uuid: i2.ArticleUuid,
			Cover:        i2.Cover,
			Title:        i2.Title,
			Abstract:     i2.Abstract,
			Visibility:   i2.Visibility,
		}
		backhtml = append(backhtml, i3)
	}
	backhtml1, err := json.Marshal(backhtml)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "加载失败",
			"jwt":     jwtString,
		})
		return
	}
	previous1, err := json.Marshal(result.Previous)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "加载失败",
			"jwt":     jwtString,
		})
		return
	}
	//开始返回数据
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "加载成功",
		"jwt":     jwtString,
		"data":    string(backhtml1),
		"time":    string(previous1),
	})
}

func Get_Change_Blog(ctx *gin.Context) { //获取修改博客文章内容
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取发送过来的uuid和文章内容的uuid
	message := modelss.Blog_Elastic_loadstorage{}
	//解析数据
	err := ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "未获取到有效uuid",
			"jwt":     jwtString,
		})
		return
	}
	//获取用户的uuid
	uuid := ctx.GetHeader("uuid")
	if uuid == "" {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "用户不是有效身份",
			"jwt":     jwtString,
		})
		return
	}
	//连接grpc服务
	//连接grpc服务进行获取数据
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "未获取到有效uuid",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	userBlogsClient := proto_blog.NewBlogsClient(conn)
	//开始检查是否是有效身份并检查文章id
	_, err = userBlogsClient.Check_Blog_Uuid_MysqlBlog(context.Background(), &proto_blog.BlogProcessBlog{
		ArticleUuid: message.Article_Uuid,
		UserUuid:    uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "用户没有权限修改这篇文章",
			"jwt":     jwtString,
		})
		return
	}
	//获取文章信息
	result1, err := userBlogsClient.Get_Blog_Mysql_Message(context.Background(), &proto_blog.Blog_ArticleUuid{
		ArticleUuid: message.Article_Uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "205",
			"message": "未拿到数据",
			"jwt":     jwtString,
		})
		return
	}
	//创建一个返回数据的对象
	blog1 := modelss.BLOG{
		Article_Uuid:   result1.ArticleUuid,
		User_Uuid:      result1.UserUuid,
		Title:          result1.Title,
		Content:        result1.Content,
		Label:          result1.Label,
		Cover:          result1.Cover,
		Abstract:       result1.Abstract,
		Reading_Volume: int(result1.Readingvolume),
		Visibility:     result1.Visibility,
		Article_Type:   result1.ArticleType,
		Process:        result1.Process,
		Creat_Time:     result1.CreatTime,
	}
	blog2, err := json.Marshal(blog1)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "205",
			"message": "未拿到数据",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "206",
		"message": "成功拿到数据",
		"jwt":     jwtString,
		"blog":    string(blog2),
	})
}
