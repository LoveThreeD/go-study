package viper

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/spf13/viper"
)

var (
	LevelConf LevelConfig
	TaskConf  TaskConfig
	Conf      Config
)

type LevelConfig struct {
	Level []Level `mapstructure:"list"`
}
type Level struct {
	Id           int   `mapstructure:"id"`
	TaskList     []int `mapstructure:"task_list"`
	FinishReward int   `mapstructure:"finish_reward"`
}

type TaskConfig struct {
	Task []Task `mapstructure:"list"`
}

type Task struct {
	Id          int `mapstructure:"id"`
	RewardScore int `mapstructure:"reward_score"`
}

type Mysql struct {
	Address  string `toml:"Address"`
	Port     int    `toml:"Port"`
	Username string `toml:"Username"`
	Password string `toml:"Password"`
	DbName   string `toml:"DbName"`
	URL      string `toml:"Url"`
}
type Redis struct {
	Address string `toml:"Address"`
	Type    string `toml:"Type"`
	Port    int    `toml:"Port"`
}

type Mongo struct {
	AuthMechanism string `toml:"AuthMechanism"`
	AuthSource    string `toml:"AuthSource"`
	Username      string `toml:"Username"`
	Password      string `toml:"Password"`
	Address       string `toml:"Address"`
	Port          int    `toml:"Port"`
}

type Config struct {
	Mysql Mysql `toml:"Mysql"`
	Redis Redis `toml:"Redis"`
	Mongo Mongo `toml:"Mongo"`
}

func init() {
	viper.SetConfigName("level") //找寻文件的名字
	viper.SetConfigType("json")  // 找寻文件的类型
	viper.AddConfigPath("conf")  //.代表当前文件夹找寻，可以多个目录找寻，生成数组
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if v, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Info(v)
		} else {
			logger.Fatal(err)
		}
	}
	//将配置文件反序列化为结构体
	LevelConf = LevelConfig{}
	if err := viper.Unmarshal(&LevelConf); err != nil {
		logger.Fatal(err)
	}

}

func init() {
	viper.SetConfigName("task") //找寻文件的名字
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if v, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Info(v)
		} else {
			logger.Fatal(err)
		}
	}
	//将配置文件反序列化为结构体
	TaskConf = TaskConfig{}
	if err := viper.Unmarshal(&TaskConf); err != nil {
		logger.Fatal(err)
	}
}

func init() {
	viper.SetConfigName("conf") //找寻文件的名字
	viper.SetConfigType("toml") // 找寻文件的类型
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if v, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Info(v)
		} else {
			logger.Fatal(err)
		}
	}
	//将配置文件反序列化为结构体
	Conf = Config{}
	if err := viper.Unmarshal(&Conf); err != nil {
		logger.Fatal(err)
	}
}
