package conf

import (
	"encoding/json"
	"github.com/asim/go-micro/v3/logger"
	"os"
)

type GAMECONFIG struct {
	LevelConfig LevelConfig
	TaskConfig  TaskConfig
}

type LevelConfig struct {
	Level []Level `json:"list"`
}

type Level struct {
	ID           int   `json:"id"`
	TaskList     []int `json:"task_list"`
	FinishReward int   `json:"finish_reward"`
}

type TaskConfig struct {
	Task []Task `json:"list"`
}

type Task struct {
	ID          int `json:"id"`
	RewardScore int `json:"reward_score"`
}

var GameConfig GAMECONFIG

func InitConfig() {
	// 打开文件
	file, _ := os.Open("conf/level.json")
	// NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	var levelConfig LevelConfig
	decoder := json.NewDecoder(file)
	// Decode从输入流读取下一个json编码值并保存在v指向的值里
	if err := decoder.Decode(&levelConfig); err != nil {
		logger.Error(err)
	}
	file.Close()

	file, _ = os.Open("conf/task.json")
	defer file.Close()
	// NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	var taskConfig TaskConfig
	decoder = json.NewDecoder(file)
	// Decode从输入流读取下一个json编码值并保存在v指向的值里
	if err := decoder.Decode(&taskConfig); err != nil {
		logger.Error(err)
	}
	GameConfig.LevelConfig = levelConfig
	GameConfig.TaskConfig = taskConfig
}
