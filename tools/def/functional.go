package def

type Function struct {
	Kind         string `yaml:"kind" json:"kind"`
	Name         string `json:"name" yaml:"name"`
	Function     string `yaml:"function" json:"function"`
	Requirements string `yaml:"requirements" json:"requirements"`
	Version      int    `yaml:"version" json:"version"`
	Image        string `yaml:"image" json:"image"` //yaml文件无此字段
	ServiceName  string //yaml文件无此字段
	PodName      string //yaml文件无此字段
	URL          string `json:"URL"` //Only for kubectl
}

type FunctionCache struct {
	Name     string
	Version  int    `yaml:"version" json:"version"`
	Image    string `yaml:"image" json:"image"` //yaml文件无此字段
	Services *ClusterIPSvc
}

type StateMachine struct {
	Name    string                 `json:"Name"`
	StartAt string                 `json:"StartAt"`
	States  map[string]interface{} `json:"States"`
	URL     string                 `json:"url"`
}

type Task struct {
	Type     string `json:"Type"`
	Resource string `json:"Resource"`
	Next     string `json:"Next"`
	End      bool   `json:"End"`
}

type Options struct {
	Variable     string `json:"Variable"`
	StringEquals string `json:"StringEquals"`
	Next         string `json:"Next"`
}

type Choice struct {
	Type    string    `json:"Type"`
	Choices []Options `json:"Choices"`
}
