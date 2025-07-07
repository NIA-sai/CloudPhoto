//这样其实挺方便的？
//面像对象不用struct怎么你了

package config

var (
	mode  string
	mysql mysqlConfig
	ai    aiConfig
)

//var config configuration
//type configuration struct {
//	Mysql mysqlConfig
//	Ai    aiConfig
//	Mode  string
//}

type mysqlConfig struct {
}

type aiConfig struct {
}

//func Get() configuration {
//	return config
//}

func Mode() string {
	return mode
}
func Mysql() mysqlConfig {
	return mysql
}

func Ai() aiConfig {
	return ai
}

func Read() {
}
