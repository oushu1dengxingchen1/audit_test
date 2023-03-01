package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)
import "github.com/olivere/elastic"

func getClient() *elastic.Client {
	var host = "http://127.0.0.1:9200/"
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetBasicAuth("elastic", "123456"))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return client
}

func searchForAll() {

	// 输入
	client := getClient()
	// 构造query,不指定字段 全局搜索
	q := elastic.NewMultiMatchQuery("con258 con563").Type("best_fields")
	// 构造时间query
	time_q := elastic.NewRangeQuery("@timestamp")
	time_q = time_q.Gte("2023-02-28T16:04:54.256+08:00").Lte("2023-02-28T16:06:54.256+08:00")
	var bool_q = elastic.NewBoolQuery().Must(q).Filter(time_q)
	src, err := bool_q.Source()
	if err != nil {
		fmt.Println(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

	searchResult, err := client.Search().Index("myindex01").Query(bool_q).From(0).Size(20).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(searchResult.Hits.TotalHits)
	fmt.Println(searchResult.Hits.Hits)
}

func AdvancedSearch() {
	seachKeyword := ""
	index_name := "myindex01"
	pageNum := 0
	pageSize := 10
	sortField := "@timestamp"
	sortOrder := -1
	sortOrderBool := false
	if sortOrder > 0 {
		sortOrderBool = true
	}
	start_time := "2023-03-01T14:48:59.190+08:00"
	end_time := "2023-03-01T14:50:59.190+08:00"
	user_name := "deng"
	database := "postgres"
	process_id := "p"
	host := "127.0."
	port := "5432"
	session_id := "con"
	cmd_id := "cmd"
	vci := "vc_default"
	audit_type := "AUDIT_LOGIN_SUCCESS"
	result := "FAILED"
	query_sql := "SELECT * FR"
	duration_start := 0
	duration_end := 50
	detail_info := "user oush"
	//
	//user_name = ""
	//database = ""
	//process_id = ""
	//host = ""
	//port = ""
	//session_id = ""
	//cmd_id = ""
	//vci = ""
	//audit_type = ""
	//result = ""
	//query_sql = ""
	//duration_start = -1
	//duration_end = -1
	//detail_info = ""

	// 输入
	client := getClient()
	var queries []elastic.Query // 构造query,不指定字段 全局搜索
	if seachKeyword != "" {
		q := elastic.NewMultiMatchQuery(seachKeyword).Type("best_fields")
		queries = append(queries, q)
	}

	// 构造时间query
	if start_time != "" && end_time != "" {
		time_q := elastic.NewRangeQuery("@timestamp")
		time_q = time_q.Gte(start_time).Lte(end_time)
		queries = append(queries, time_q)
	}

	// 用户查询  精准匹配
	if user_name != "" {
		q := elastic.NewMatchQuery("username", user_name)
		queries = append(queries, q)
	}

	// 数据库查询 精准匹配
	if database != "" {
		q := elastic.NewMatchQuery("database", database)
		queries = append(queries, q)
	}

	// process id查询
	if process_id != "" {
		q := elastic.NewWildcardQuery("process_id", getLikelyMatchString(process_id))
		queries = append(queries, q)
	}

	// host查询
	if host != "" {
		q := elastic.NewWildcardQuery("remote_host", getLikelyMatchString(host))
		queries = append(queries, q)
	}

	// port 查询
	if port != "" {
		q := elastic.NewMatchQuery("remote_port", port)
		queries = append(queries, q)
	}

	// session 查询
	if session_id != "" {
		q := elastic.NewWildcardQuery("session_id", getLikelyMatchString(session_id))
		queries = append(queries, q)
	}

	// cmd_id 查询
	if cmd_id != "" {
		q := elastic.NewWildcardQuery("cmd_id", getLikelyMatchString(cmd_id))
		queries = append(queries, q)
	}

	// vci 查询
	if vci != "" {
		q := elastic.NewWildcardQuery("vci", getLikelyMatchString(vci))
		queries = append(queries, q)
	}

	// audit_type 查询  精准匹配
	if audit_type != "" {
		audit_type = strings.ToLower(audit_type)
		q := elastic.NewTermQuery("audit_type", audit_type)
		queries = append(queries, q)
	}

	// result 查询  精准匹配
	if result != "" {
		q := elastic.NewMatchQuery("audit_result", result)
		queries = append(queries, q)
	}

	// query_sql 查询
	if query_sql != "" {
		q := elastic.NewMatchQuery("query", getLikelyMatchString(query_sql))
		queries = append(queries, q)
	}

	// detail_info 查询
	if detail_info != "" {
		q := elastic.NewMatchQuery("detail_info", getLikelyMatchString(detail_info))
		queries = append(queries, q)
	}

	// duration 查询  范围查询
	if duration_start >= 0 && duration_end >= 0 && duration_end >= duration_start {
		q := elastic.NewRangeQuery("duration")
		q = q.Gte(duration_start).Lte(duration_end)
		queries = append(queries, q)
	}

	var bool_q = elastic.NewBoolQuery()
	for _, q := range queries {
		bool_q = bool_q.Must(q)
	}

	src, err := bool_q.Source()
	if err != nil {
		fmt.Println(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

	searchResult, err := client.Search().Index(index_name).Query(bool_q).Sort(sortField, sortOrderBool).From(pageNum).Size(pageSize).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(searchResult.Hits.TotalHits)
	fmt.Println(searchResult.Hits.Hits)
}

func getLikelyMatchString(s string) string {
	return "*" + s + "*"
}

func main() {
	AdvancedSearch()
}
