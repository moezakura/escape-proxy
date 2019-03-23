package model

type ConfigYaml struct {
	Auth              bool        `yaml:"auth"`
	AutoDirectConnect bool        `yaml:"auto_direct_connect"`
	ProxyServer       string      `yaml:"proxy"`
	GatewayServer     string      `yaml:"gateway"`
	Listen            string      `yaml:"listen"`
	LogMode           string      `yaml:"log_mode"`
	ExcludeIps        []string    `yaml:"exclude"`
	Users             []AuthUsers `yaml:"users"`
}

type AuthUsers struct {
	Id       string `yaml:"id"`
	Password string `yaml:"password"`
}
