package _struct

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

// 配置历史
type ConfHistory struct {
	gorm.Model
	HistoryVersion     string              `json:"history_version"`
	UserID             int                 `json:"user_id"`
	Product            string              `json:"product"`
	ClusterGroupId     int                 `json:"cluster_group_id"` // 冗余字段，与ConfGroupId一同用来变更配置组
	InUse              bool                `json:"in_use"`           // true:当前集群正在使用；false表示还没有生效到集群上
	ExpiredMark        bool                `json:"expired_mark"`     // ture:表示配置已经过期，出现了新的配置
	ConfHistoryDetails []ConfHistoryDetail `json:"conf_history_details" gorm:"foreignKey:ConfHistoryId;references:ID"`
	Machines           []Machine           `gorm:"many2many:conf_history_refer_machine;constraint:OnDelete:CASCADE;"`
}

// 历史详情
type ConfHistoryDetail struct {
	gorm.Model
	ConfHistoryId int    `json:"conf_history_id"`
	FileDir       string `json:"file_dir"`
	FileName      string `json:"file_name"`
	FileType      string `json:"file_type"`
	FileOwner     string `json:"file_owner"`
	FileContent   string `json:"file_content"`
}

type Machine struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func InitDB() {
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
	err = db.AutoMigrate(&ConfHistory{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&ConfHistoryDetail{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&Machine{})
	if err != nil {
		fmt.Println(err)
	}

	// 插入数据
	//db.Create(&Machine{ID: 1, Name: "machine1"})
	//db.Create(&Machine{ID: 2, Name: "machine2"})

	cluserGroupID := 33345
	// 其他记录记为false
	result := db.Model(ConfHistory{}).Where("cluster_group_id = ?", cluserGroupID).Updates(map[string]interface{}{"in_use": false, "expired_mark": false})
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	fmt.Println("更新行数:", result.RowsAffected)

	newConf := ConfHistory{
		HistoryVersion: "v1", UserID: 12345, Product: "Audit", ClusterGroupId: cluserGroupID,
		InUse: true, ExpiredMark: true,
		ConfHistoryDetails: []ConfHistoryDetail{{
			FileDir: "./", FileName: "filebeat2.yml",
			FileContent: `{"paths":["/home/deng/environment/filebeat/log/test_mapping.log"],
                          "shards": 2, "replications": 2,
                          "scan_frequency":10, "max_bytes":5242880, "harvester_buffer_size": 16384,
                          "index_name":"myindex02", "pipeline":"testFileBeat01"}`,
		}},
		Machines: []Machine{{ID: 1}, {ID: 2}},
	}

	result = db.Create(&newConf)
	if result.Error != nil {
		fmt.Println(err)
	}
	if result.RowsAffected > 0 {
		fmt.Println("插入数据成功")
	}
	db.Commit()

}
