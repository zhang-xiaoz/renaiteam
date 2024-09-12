package initialization

import (
	"context"
	"srv/log"

	"github.com/olivere/elastic/v7"
)

var (
	Elastic *elastic.Client
	err     error
	mapping = `
	{
		"mappings":{
			"properties": {
				"article_uuid": {
					"type": "keyword"
				},
				"user_uuid": {
					"type": "keyword"
				},
				"title": {
					"type": "text",
					"analyzer": "ik_max_word"
				},
				"label": {
					"type": "text",
					"analyzer": "ik_smart"
				},
				"abstract": {
					"type": "text",
					"analyzer": "ik_smart"
				},
				"article_type": {
					"type": "keyword"
				},
				"visibility":{
					"type": "keyword"
				},
				"process": {
					"type": "keyword"
				},
				"reading_volume": {
					"type": "long"
				},
				"creat_time": {
					"type": "long"
				}
			}
		}
	}`
	// "settings":{
	// 	"analysis": {
	// 		"analyzer": {
	// 		  "ik_max_word": {
	// 			"type": "ik",
	// 			"use_smart": false
	// 		  },
	// 		  "ik_smart": {
	// 			"type": "ik",
	// 			"use_smart": true
	// 		  }
	// 		}
	// 	}
	// },
)

func InitElastic() {
	Elastic, err = elastic.NewClient(
		elastic.SetSniff(false),                     // SetSniff启用或禁用嗅探器（默认情况下启用）。
		elastic.SetURL("http://106.14.30.173:9200"), // URL地址
		elastic.SetBasicAuth("elastic", "wasd2002"), // 账号密码
	)
	if err != nil {
		log.SugarLogger.Panic(err)
		return
	}
	exists, err := Elastic.IndexExists("renai_blog").Do(context.Background())
	if err != nil {
		log.SugarLogger.Panic(err)
		return
	}
	if !exists { //不存在从新创建
		_, err = Elastic.CreateIndex("renai_blog").BodyString(mapping).Do(context.Background())
		if err != nil {
			log.SugarLogger.Panic(err)
			return
		}
	}
}
