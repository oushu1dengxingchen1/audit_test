package main

import (
	"dxctest.com/m/struct"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"strconv"
)

func main() {
	dsn := "host=localhost user=deng password=123456 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("数据库连接成功")

	// 创建table
	err = db.AutoMigrate(&_struct.ConfHistory{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&_struct.ConfHistoryDetail{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&_struct.Machine{})
	if err != nil {
		fmt.Println(err)
	}

	machineID := "1"
	sql := " SELECT m.name , history_version, product, cluster_group_id, file_content"
	sql += " FROM (conf_history cf LEFT JOIN conf_history_detail cfd ON cf.id = cfd.conf_history_id"
	sql += " LEFT JOIN conf_history_refer_machine cfm ON cf.id = cfm.conf_history_id"
	sql += " LEFT JOIN machine m ON cfm.machine_id = m.id)"
	sql += " WHERE cf.in_use=true AND cfm.machine_id=" + machineID

	fmt.Println(sql)

	rows, err := db.Raw(sql).Rows()
	defer rows.Close()

	FilebeatConfig := new(_struct.FilebeatYaml)
	FilebeatConfig.FBConfigModules.Path = "${path.config}/modules.d/*.yml"
	FilebeatConfig.FBConfigModules.ReloadEnabled = false
	FilebeatConfig.SetUpIlmEnabled = false
	FilebeatConfig.QueueMemory.Events = 4096
	FilebeatConfig.QueueMemory.FlushTime = "1s"
	FilebeatConfig.QueueMemory.FlushMinEvents = 2048

	for rows.Next() {
		var version string
		var product string
		var cluster_group_id int
		var file_content string
		var configMap map[string]interface{}
		// 取值
		err := rows.Scan(&version, &product, &cluster_group_id, &file_content)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(file_content)
		// 解析配置
		err = json.Unmarshal([]byte(file_content), &configMap)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(configMap)

		matchTagStr := product + "_" + strconv.Itoa(cluster_group_id)
		matchTag := map[string]string{"match_tag": matchTagStr}
		var input _struct.FilebeatInput
		input.InputType = "log"
		input.Enabled = true

		for _, item := range configMap["paths"].([]interface{}) {
			input.Paths = append(input.Paths, item.(string))
		}
		input.Fields = matchTag
		input.ScanFrequency = strconv.Itoa(int(configMap["scan_frequency"].(float64))) + "s"
		input.MaxBytes = int(configMap["max_bytes"].(float64))
		input.HarvesterBufferSize = int(configMap["harvester_buffer_size"].(float64))
		FilebeatConfig.FBInputs = append(FilebeatConfig.FBInputs, input)

		// output
		FilebeatConfig.OutputES.ESHosts = []string{"localhost:9200"}
		FilebeatConfig.OutputES.BulkMaxSize = 1000

		var index _struct.TargetESIndex
		index.IndexName = configMap["index_name"].(string)
		index.WhenEquals.Fields = matchTag
		FilebeatConfig.OutputES.TargetIndices = append(FilebeatConfig.OutputES.TargetIndices, index)

		var pipeline _struct.Pipeline
		pipeline.PipelineName = configMap["pipeline"].(string)
		pipeline.WhenEquals.Fields = matchTag
		FilebeatConfig.OutputES.Pipelines = append(FilebeatConfig.OutputES.Pipelines, pipeline)

		FilebeatConfig.OutputES.UserName = "elastic"
		FilebeatConfig.OutputES.PassWord = "123456"
	}

	fmt.Println(FilebeatConfig)
	bytes, err := yaml.Marshal(FilebeatConfig)
	if err != nil {
		return
	}

	err = ioutil.WriteFile("/home/deng/code/oushu/test/filebeat2.yml", bytes, 0666)
	if err != nil {
		fmt.Println(err)
	}
}
