package services

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/astaxie/beego"
	"fmt"
	"awesomeProject/kafa-job/models"
	"encoding/json"
)

const mapping =`{
  "settings": {
    "number_of_shards": 5,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "ik_pinyin_analyzer": {
          "type": "custom",
          "tokenizer": "ik_smart",
          "filter": [
            "my_pinyin",
            "word_delimiter"
          ]
        }
      },
      "filter": {
        "my_pinyin": {
          "type": "pinyin",
          "first_letter": "prefix",
          "padding_char": " "
        }
      }
    }
  },
  "mappings": {
    "item_type": {
      "properties": {
        "item_id": {
          "type": "keyword"
        },
        "item_name": {
          "type": "keyword",
          "store": true,
          "fields": {
            "ikPinyin": {
              "type": "text",
              "analyzer": "ik_pinyin_analyzer"
            }
          }
        },
        "item_no": {
          "type": "keyword"
        },
        "bar_code": {
          "type": "keyword"
        },
        "first_category_id": {
          "type": "keyword"
        },
        "first_category": {
          "type": "keyword",
          "fields": {
            "ikPinyin": {
              "type": "text",
              "analyzer": "ik_pinyin_analyzer"
            }
          }
        },
        "second_category_id": {
          "type": "keyword"
        },
        "second_category": {
          "type": "keyword",
          "fields": {
            "ikPinyin": {
              "type": "text",
              "analyzer": "ik_pinyin_analyzer"
            }
          }
        },
        "item_unit": {
          "type": "keyword"
        },
        "item_size": {
          "type": "keyword"
        },
        "image_url": {
          "type": "keyword"
        },
        "stock_size": {
          "type": "integer"
        },
        "brand": {
          "type": "keyword"
        },
        "create_time": {
          "type": "date"
        },
        "modify_time": {
          "type": "date"
        },
        "image_middle_url": {
          "type": "keyword"
        },
        "completion_snap_shoot": {
          "type": "completion",
          "analyzer": "ik_pinyin_analyzer"
        },
        "term_snap_shoot": {
          "type": "text",
          "analyzer": "ik_pinyin_analyzer"
        }
      }
    }
  }
}`
func index(msg []byte)  {
	ctx:=context.Background()
	client,err:=elastic.NewClient()
	if err!=nil{
		panic(err)

	}
	exist,err:=client.IndexExists("item_index").Do(ctx);
	if err!=nil{
		panic(err)
	}
	if !exist{
		//create index
		createIndex,err:=client.CreateIndex("item_index").Body(mapping).Do(ctx)
		if err !=nil{
			panic(err)
		}
		if !createIndex.Acknowledged {
			beego.Error(fmt.Sprintln("item_index not created successed"))
			return
		}
	}
	item:=&models.ElasticItem{}
	json.Unmarshal(msg,item)
	response,err:=client.Index().Index("item_index").Type("item_type").Id((*item).ItemId).BodyString(string(msg)).Do(ctx)
	if err!=nil{
		panic(err)
	}
	beego.Info(fmt.Sprintf("索引文件返回状态:%d",response.Status))

}