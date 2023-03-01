package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.OpenFile("/home/deng/environment/filebeat/log/test_mapping.log", os.O_WRONLY|os.O_APPEND, 0666)
	file2, err := os.OpenFile("/home/deng/environment/filebeat/log/test_filebeat.log", os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	defer file2.Close()

	write := bufio.NewWriter(file)
	//write2 := bufio.NewWriter(file2)
	users := []string{"oushu", "oushutest", "deng"}
	databases := []string{"postgres", "magmaap", "123abc_"}
	processID := "p"
	remoteHost := "127.0.0.1"
	remotePort := "5432"
	sessionID := "con"
	cmdID := "cmd"
	vci := "vc_default_1"
	level := "AUDIT"
	auditTypes := []string{"AUDIT_LOGIN_SUCCESS", "AUDIT_LOGIN_FAILURE", "AUDIT_LOGOUT", "AUDIT_DDL_XXX", "AUDIT_DML_NOT_SELECT", " AUDIT_DML_SELECT"}
	auditResults := []string{"OK", "FAILED"}
	querys := []string{
		"UPDATE student SET sname = '张三' WHERE  sid = 1",
		"DELETE FROM 表student WHERE id = 1",
		"SELECT * FROM Persons WHERE FirstName='Thomas' AND LastName='Carter'",
		"SELECT * FROM Persons WHERE firstname='Thomas' OR lastname='Carter'",
		"INSERT INTO locations(location_id,street_address,postal_code,city,state_province,country_id) VALUES (1400,'2014 Jabberwocky Rd','26192','Southlake','Texas','US');",
		"DROP TABLE regions;",
	}
	effectedRows := 0
	duration := 0

	for {
		timeStr := time.Now().Format("2006-01-02T15:04:05.000+08:00")
		userName := users[rand.Intn(len(users))]
		database := databases[rand.Intn(len(databases))]
		processID = "p" + strconv.Itoa(rand.Intn(1000))
		sessionID = "con" + strconv.Itoa(rand.Intn(1000))
		cmdID = "cmd" + strconv.Itoa(rand.Intn(1000))
		auditType := auditTypes[rand.Intn(len(auditTypes))]
		auditResult := auditResults[rand.Intn(2)]
		query := querys[rand.Intn(len(querys))]
		effectedRows = rand.Intn(100)
		duration = rand.Intn(50)
		detailInfo := "user oushu login "
		parseTree := "[{\"sql\": \"select * from t\", \"session_id\": 1048621, \"sql_type\": \"SELECT\", \"table_schema\": [{\"schema_name\": \"public\", \"table_name\": \"t\", \"alias_name\": null, \"oid\": 16526}], \"query_tree\": [{\"QUERY\": {\"commandType\": 1, \"querySource\": 0, \"canSetTag\": true, \"utilityStmt\": null, \"resultRelation\": 0, \"intoClause\": null, \"hasAggs\": false, \"hasWindFuncs\": false, \"hasSubLinks\": false, "
		s := fmt.Sprintf("[%s] %s %s %s %s:%s %s %s %s %s %s %s <query>%s</query> %d %dms <info>%s</info> <tree>%s</tree>\n",
			timeStr, userName, database, processID, remoteHost, remotePort, sessionID, cmdID, vci, level, auditType, auditResult, query, effectedRows, duration, detailInfo, parseTree)

		write.WriteString(s)
		//write2.WriteString(s)
		fmt.Println("写入成功: " + s)
		//Flush将缓存的文件真正写入到文件中
		write.Flush()
		//write2.Flush()
		time.Sleep(1 * time.Second)
	}

}
