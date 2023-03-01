package _struct

type FilebeatYaml struct {
	FBInputs        []FilebeatInput       `yaml:"filebeat.inputs"`
	FBConfigModules FilebeatConfigModules `yaml:"filebeat.config.modules"`
	SetUpIlmEnabled bool                  `yaml:"setup.ilm.enabled"`
	OutputES        OutputElasticsearch   `yaml:"output.elasticsearch"`
	QueueMemory     QueueMemorySetting    `yaml:"queue.mem"`
}

type FilebeatInput struct {
	InputType           string            `yaml:"type"`
	Enabled             bool              `yaml:"enabled"`
	Paths               []string          `yaml:"paths"`
	Fields              map[string]string `yaml:"fields"`
	ScanFrequency       string            `yaml:"scan_frequency"`
	MaxBytes            int               `yaml:"max_bytes"`
	HarvesterBufferSize int               `yaml:"harvester_buffer_size"`
}

type FilebeatConfigModules struct {
	Path          string `yaml:"path"`
	ReloadEnabled bool   `yaml:"reload.enabled"`
}

type OutputElasticsearch struct {
	ESHosts []string `yaml:"hosts"`
	//DefaultPipeline string          `yaml:"pipeline"`
	TargetIndices []TargetESIndex `yaml:"indices"`
	Pipelines     []Pipeline      `yaml:"pipelines"`
	UserName      string          `yaml:"username"`
	PassWord      string          `yaml:"password"`
	BulkMaxSize   int             `yaml:"bulk_max_size"`
}

type TargetESIndex struct {
	IndexName  string         `yaml:"index"`
	WhenEquals MatchCondition `yaml:"when.equals"`
}

type Pipeline struct {
	PipelineName string         `yaml:"pipeline"`
	WhenEquals   MatchCondition `yaml:"when.equals"`
}

type MatchCondition struct {
	Fields map[string]string `yaml:"fields"`
}

type QueueMemorySetting struct {
	Events         int    `yaml:"events"`
	FlushMinEvents int    `yaml:"flush.min_events"`
	FlushTime      string `yaml:"flush.timeout"`
}
