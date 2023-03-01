package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
)

var host = "http://127.0.0.1:9200/"

func newPipeline(esclient *elastic.Client, name string, body string) *elastic.IngestPutPipelineResponse {
	res, err := esclient.IngestPutPipeline(name).BodyJson(body).Do(context.Background())
	if err != nil {
		fmt.Println(res)
	}
	return res
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetBasicAuth("elastic", "123456"))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	jsonBody :=
		`{
		  "processors": [
			{
			  "dissect": {
				"field": "message",
				"pattern": "[%{log_time}] %{username} %{database} %{process_id} %{remote_host}:%{remote_port} %{session_id} %{cmd_id} %{vci} %{level} %{audit_type} %{audit_result} <query>%{query}</query> %{effected_rows} %{duration}ms <info>%{detail_info}</info> <tree>%{content}</tree>"
			  }
			},
			{
			  "date": {
				  "field": "log_time",
				  "formats": ["yyyy-MM-dd'T'HH:mm:ss.SSSZ"],
                  "timezone": "Asia/Shanghai", 
				  "target_field": "@timestamp"
			  }
			}
		  ]
        }`
	res := newPipeline(client, "testFileBeat01", jsonBody)
	fmt.Println(res)
}
