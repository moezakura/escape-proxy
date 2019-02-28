package model

type ConfigYaml struct {
	Auth          bool        `yaml:"auth"`
	ProxyServer   string      `yaml:"proxy"`
	GatewayServer string      `yaml:"gateway"`
	Listen        string      `yaml:"listen"`
	Users         []AuthUsers `yaml:"users"`
}

type AuthUsers struct {
	Id       string `yaml:"id"`
	Password string `yaml:"password"`
}
