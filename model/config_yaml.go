package model

type ConfigYaml struct {
	Auth  bool        `yaml:"auth"`
	Users []AuthUsers `yaml:"users"`
}

type AuthUsers struct {
	Id       string `yaml:"id"`
	Password string `yaml:"password"`
}
