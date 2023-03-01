package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"io/ioutil"
)

var host = "http://127.0.0.1:9200/"

func main() {
	templateName := "temptest3"
	indexNamePattern := "myindex01*"
	ilmPolicy := "myindex01_policy"
	rolloverAlias := "myindex01"
	numShards := 2
	numReplication := 2

	// policy config
	rolloverDays := 8
	maxRetentDays := 33

	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetBasicAuth("elastic", "123456"))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// ilm policy 注册
	// 检查policy是否存在
	getilm, err := client.XPackIlmGetLifecycle().Policy(ilmPolicy).Pretty(true).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	if getilm == nil {
		fmt.Println("didn't get lifecycle response")
	}

	// Check the policy exists
	_, found := getilm[ilmPolicy]
	if found {
		fmt.Printf("ilm policy %s exist, skip create \n", ilmPolicy)
	} else {
		// 读取json
		data, err := ioutil.ReadFile("/home/deng/code/oushu/test/policy.json")
		if err != nil {
			fmt.Println("read file err:", err.Error())
			return
		}
		mappingJSon := string(data)
		mappingJSon = fmt.Sprintf(mappingJSon, rolloverDays, maxRetentDays)
		fmt.Println(mappingJSon)

		// Create the policy
		putilm, err := client.XPackIlmPutLifecycle().Policy(ilmPolicy).BodyString(mappingJSon).Do(context.TODO())
		if err != nil {
			fmt.Println(err)
		}
		if putilm == nil {
			fmt.Println("didn't get put policy response")
		}
		if !putilm.Acknowledged {
			fmt.Println("put ilm policy ack false", putilm.Acknowledged)
		}
	}

	// template注册
	// 检查template是否存在
	templateExists, err := client.IndexTemplateExists(templateName).Do(context.TODO())
	if templateExists {
		fmt.Printf("template %s exists, skip create\n", templateName)
	} else {
		// 读取json
		data, err := ioutil.ReadFile("/home/deng/code/oushu/test/template.json")
		if err != nil {
			fmt.Println("read file err:", err.Error())
			return
		}
		mappingJSon := string(data)
		mappingJSon = fmt.Sprintf(mappingJSon, indexNamePattern, numShards, numReplication, ilmPolicy, rolloverAlias)
		fmt.Println(mappingJSon)
		putresp, err := client.IndexPutTemplate(templateName).BodyString(mappingJSon).Create(true).Do(context.TODO())
		if err != nil {
			fmt.Println(err)
		}
		if putresp == nil {
			fmt.Printf("expected put mapping response; got: %v", putresp)
		}
		if !putresp.Acknowledged {
			fmt.Printf("expected put mapping ack; got: %v", putresp.Acknowledged)
		}
	}

	// 查看template
	res, err := client.IndexGetTemplate(templateName).Pretty(true).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	if res == nil {
		fmt.Printf("expected result; got: %v", res)
	}
	template := res[templateName]
	if template == nil {
		fmt.Printf("expected template %q to be != nil; got: %v", "template_1", template)
	}
	if len(template.IndexPatterns) != 1 || template.IndexPatterns[0] != indexNamePattern {
		fmt.Printf("expected index settings of %q to be [\"index1\"]; got: %v", indexNamePattern, template.IndexPatterns)
	}

	// 创建第一个Index
	// 检查index是否存在
	indexExists, err := client.IndexExists(rolloverAlias).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	if indexExists {
		fmt.Printf("index %s exist, skip create\n", rolloverAlias)
	} else {
		// 创建index
		testIndexName := fmt.Sprintf("<%s-{now/d}-1>", rolloverAlias)
		mappingJSon := fmt.Sprintf(`{"aliases": 
										{
											"%s":{
											  "is_write_index": true
											}
										}
                                      }`, rolloverAlias)
		createIndex, err := client.CreateIndex(testIndexName).BodyJson(mappingJSon).Do(context.TODO())
		if err != nil {
			fmt.Println(err)
		}
		if !createIndex.Acknowledged {
			fmt.Println("expected IndicesCreateResult.Acknowledged %v; got %v", true, createIndex.Acknowledged)
		}
	}

}
