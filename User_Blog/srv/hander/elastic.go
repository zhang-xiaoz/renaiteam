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

	"github.com/olivere/elastic/v7"
)

//所有elastic服务放在这里

func get_Blog_Message_Elastic_TopBlog1(a int, d string) (*[]models.Blog_Elastic, int64, error) {
	// 构建查询
	searchQuery := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("label", d),
		elastic.NewMatchQuery("process", "已经审核"),
		elastic.NewTermsQuery("visibility", "全体可见", "仁爱成员可见"),
	)
	//创建搜索服务
	searchService := elastic.NewSearchService(initialization.Elastic).
		Index("renai_blog").       // 指定索引名
		Query(searchQuery).        // 设置查询体
		Sort("creat_time", false). // 根据 creat_time 字段排序，false 表示降序
		Size(4).                   // 设置返回结果的数量
		SearchAfter(a)             // 使用 search_after 参数
		// 执行搜索请求
	res, err := searchService.Do(context.Background())
	if err != nil {
		return nil, 0, err
	}
	var b []models.Blog_Elastic
	number := len(res.Hits.Hits)
	fmt.Println(number)
	for _, hit := range res.Hits.Hits {
		blog := models.Blog_Elastic{}
		err := json.Unmarshal(hit.Source, &blog)
		if err != nil {
			return nil, 0, err
		}
		b = append(b, blog)
	}
	return &b, int64(number), nil
}

func check_Blog_Uuid_Elastic_Blog1(blog []models.Blog_Elastic_loadstorage) bool {
	for _, i2 := range blog {
		get_reponse, err := initialization.Elastic.Get().Index("renai_blog").Id(i2.Article_Uuid).Do(context.Background())
		if err != nil {
			log.SugarLogger.Error(err)
			return false
		}
		if !get_reponse.Found { //没有找到
			return false
		}
	}
	return true
}

func get_Blog_Message_ArticleUuid1(a string) (*proto.Blog_Elastic_Message_TopBlog, error) {
	get_reponse, err := initialization.Elastic.Get().Index("renai_blog").Id(a).Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, err
	}
	if !get_reponse.Found { //没有找到
		return nil, errors.New("未找到")
	}
	c := get_reponse.Source
	cc := models.Blog_Elastic{}
	err = json.Unmarshal(c, &cc)
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, err
	}
	back1 := proto.Blog_Elastic_Message_TopBlog{
		ArticleUuid: cc.Article_Uuid,
		Title:       cc.Title,
		Cover:       cc.Cover,
		Abstract:    cc.Abstract,
		Visibility:  cc.Visibility,
	}
	return &back1, nil
}

func get_Blog_Message_Elastic_Lower1(a int64) (*[]models.Blog_Elastic, int64, int64, error) {
	// 构建查询
	var n int64
	searchQuery := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("process", "已经审核"),
		elastic.NewTermsQuery("visibility", "全体可见", "仁爱成员可见"),
	)
	//创建搜索服务
	searchService := elastic.NewSearchService(initialization.Elastic).
		Index("renai_blog").       // 指定索引名
		Query(searchQuery).        // 设置查询体
		Sort("creat_time", false). // 根据 creat_time 字段排序，false 表示降序
		Size(4).                   // 设置返回结果的数量
		SearchAfter(a)             // 使用 search_after 参数
		// 执行搜索请求
	res, err := searchService.Do(context.Background())
	if err != nil {
		return nil, 0, a, err
	}
	var b []models.Blog_Elastic
	number := len(res.Hits.Hits)
	for _, hit := range res.Hits.Hits {
		blog := models.Blog_Elastic{}
		err := json.Unmarshal(hit.Source, &blog)
		if err != nil {
			return nil, 0, a, err
		}
		n = int64(blog.Creat_Time)
		b = append(b, blog)
	}
	return &b, int64(number), n, nil
}

func search_Blog_ElasticBlog1(A string, B string, C int, D int) (*[]models.Blog_Elastic_Html, int, int, string, error) {
	var n1 int
	var n2 string
	fmt.Println(A)
	fmt.Println(D)
	// 构建查询
	searchQuery := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("process", "已经审核"),
		elastic.NewTermsQuery("visibility", "全体可见", "仁爱成员可见"),
	).Should(
		elastic.NewMatchQuery("title", A),
		elastic.NewMatchQuery("label", A),
		elastic.NewMatchQuery("abstract", A),
	).MinimumShouldMatch("1")
	//创建搜索服务
	searchService := elastic.NewSearchService(initialization.Elastic).
		Index("renai_blog").           // 指定索引名
		Query(searchQuery).            // 设置查询体
		Sort("reading_volume", false). // 根据 creat_time 字段排序，false 表示升序
		Sort("article_uuid", false).   // 根据 _id 字段升序排序
		Size(D).                       // 设置返回结果的数量
		SearchAfter(C, B)              // 使用 search_after 参数
		// 执行搜索请求
	res, err := searchService.Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, 0, C, B, err
	}
	var b2 []models.Blog_Elastic_Html
	number := len(res.Hits.Hits)
	for _, hit := range res.Hits.Hits {
		blog1 := models.Blog_Elastic{}
		err := json.Unmarshal(hit.Source, &blog1)
		if err != nil {
			log.SugarLogger.Error(err)
			return nil, 0, C, B, err
		}
		blog2 := models.Blog_Elastic_Html{
			Cover:        blog1.Cover,
			Article_Uuid: blog1.Article_Uuid,
			Title:        blog1.Title,
			Abstract:     blog1.Abstract,
			Visibility:   blog1.Visibility,
		}
		n1 = blog1.Reading_Volume
		n2 = blog1.Article_Uuid
		b2 = append(b2, blog2)
	}
	return &b2, number, n1, n2, nil
}

func add_Blog_Reading_Volume1(a string) error {
	updateRequest := elastic.NewUpdateByQueryService(initialization.Elastic).
		Index("renai_blog").
		Query(elastic.NewTermQuery("id", a)).
		Script(elastic.NewScript("ctx._source.reading_volume += 1").
			Type("painless"))
	res, err := updateRequest.Do(context.Background())
	if err != nil {
		return err
	}
	// 检查更新是否成功
	if res.Updated != 1 {
		return errors.New("未更新成功")
	} else {
		return nil
	}
}

func delete_El6_Mysql_blog1(a string) error {
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaa3")
	fmt.Println(a)
	tx := initialization.DB.Begin() //开始事务
	result := tx.Where("article_uuid=?", a).Delete(&models.Blog{})
	if result.Error != nil {
		log.SugarLogger.Error(result.Error)
		tx.Rollback() // 发生错误回滚事务
		return result.Error
	}
	//提交事务
	tx.Commit()
	_, err := initialization.Elastic.Delete().Index("renai_blog").Id(a).Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(err)
		return err
	}
	return nil
}

// 后边分别返回 得到的数据 得到数据最后一个数据  得到数据最后一个数据的时间 总共得到数据的数量 错误
func get_Delete_Blog_Mysql1(A4 int64, A5 int64, A1, A2, A3 string) (*[]models.Blog_Elastic_Html, string, int, int, error) {
	var n1 int
	var n2 string
	// 构建查询
	searchQuery := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("user_uuid", A1), //用户id
		elastic.NewMatchQuery("process", A2),   //审核状态
	)
	//创建搜索服务
	searchService := elastic.NewSearchService(initialization.Elastic).
		Index("renai_blog").         // 指定索引名
		Query(searchQuery).          // 设置查询体
		Sort("creat_time", false).   // 根据 creat_time 字段排序 false 降序 54321
		Sort("article_uuid", false). // 根据 _id 字段升序排序
		Size(int(A5)).               // 设置返回结果的数量
		SearchAfter(A4, A3)          // 使用 search_after 参数
	res, err := searchService.Do(context.Background()) // 执行搜索请求
	if err != nil {
		return nil, "", 0, 0, err
	}
	var b2 []models.Blog_Elastic_Html
	number := len(res.Hits.Hits)
	for _, hit := range res.Hits.Hits {
		blog1 := models.Blog_Elastic{}
		err := json.Unmarshal(hit.Source, &blog1)
		if err != nil {
			return nil, "", 0, 0, err
		}
		blog2 := models.Blog_Elastic_Html{
			Cover:        blog1.Cover,
			Article_Uuid: blog1.Article_Uuid,
			Title:        blog1.Title,
			Abstract:     blog1.Abstract,
			Visibility:   blog1.Visibility,
		}
		n1 = blog1.Creat_Time
		n2 = blog1.Article_Uuid
		b2 = append(b2, blog2)
	}
	return &b2, n2, n1, number, nil
}

func search_Blog_Label_Elastic1(A string, B string, C int, D int) (*[]models.Blog_Elastic_Html, int, int, string, error) {
	var n1 int
	var n2 string
	fmt.Println(A)
	fmt.Println(D)
	// 构建查询
	searchQuery := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("process", "已经审核"),
		elastic.NewTermsQuery("visibility", "全体可见", "仁爱成员可见"),
		elastic.NewMatchQuery("label", A),
	)
	//创建搜索服务
	searchService := elastic.NewSearchService(initialization.Elastic).
		Index("renai_blog").         // 指定索引名
		Query(searchQuery).          // 设置查询体
		Sort("creat_time", false).   // 根据 creat_time 字段排序，false 表示升序
		Sort("article_uuid", false). // 根据 _id 字段升序排序
		Size(D).                     // 设置返回结果的数量
		SearchAfter(C, B)            // 使用 search_after 参数
		// 执行搜索请求
	res, err := searchService.Do(context.Background())
	if err != nil {
		log.SugarLogger.Error(err)
		return nil, 0, C, B, err
	}
	var b2 []models.Blog_Elastic_Html
	number := len(res.Hits.Hits)
	for _, hit := range res.Hits.Hits {
		blog1 := models.Blog_Elastic{}
		err := json.Unmarshal(hit.Source, &blog1)
		if err != nil {
			log.SugarLogger.Error(err)
			return nil, 0, C, B, err
		}
		blog2 := models.Blog_Elastic_Html{
			Cover:        blog1.Cover,
			Article_Uuid: blog1.Article_Uuid,
			Title:        blog1.Title,
			Abstract:     blog1.Abstract,
			Visibility:   blog1.Visibility,
		}
		n1 = blog1.Creat_Time
		n2 = blog1.Article_Uuid
		b2 = append(b2, blog2)
	}
	return &b2, number, n1, n2, nil
}
