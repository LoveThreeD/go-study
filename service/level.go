package service

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/golang/protobuf/proto"
	"strconv"
	"study/pkg/response"
	"study/pkg/viper"
	pb "study/proto"
	"study/repository/data"
	"study/repository/document"
)

var (
	// gameData的默认初始数据,用来初始化用户的游戏数据
	gameDataInitData []byte
	taskData         map[int]int
	levelData        map[int]int
)

func init() {
	var err error
	// 1.初始化游戏数据
	levelInit := pb.LevelData{
		CurLevel: 1,
	}
	materialInit := pb.MaterialData{
		Warehouse: make(map[uint32]uint32),
	}
	gameData := pb.GameData{
		LevelData:  &levelInit,
		ShopData:   make(map[uint32]uint32),
		Statistics: &pb.StatisticsData{},
		Setting:    &pb.GameSetting{},
		Material:   &materialInit,
	}

	gameDataInitData, err = proto.Marshal(&gameData)
	if err != nil {
		logger.Fatal(err)
	}

	// 2.初始化配置数据
	taskData = make(map[int]int, 10)
	for _, v := range viper.TaskConf.Task {
		taskData[v.Id] = v.RewardScore
	}
	levelData = make(map[int]int, 10)
	for _, v := range viper.LevelConf.Level {
		taskData[v.Id] = v.FinishReward
	}
}

/*
	进入关卡： 要求  超出需要提示后续开放   只能进入当前关卡
*/
func EnterLevel(levelId int, userID int) (*pb.GameData, error) {
	// 获取与判断
	gameData, err := getGameData(userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	levelLen := len(viper.LevelConf.Level)
	levels := viper.LevelConf.Level

	// 判断关卡是否在关卡表中
	if _, v := levelData[int(gameData.LevelData.CurLevel)]; !v {
		err = errors.New(response.MsgTaskNotFoundError)
		return nil, err
	}

	// 判断是否超出
	if levelId > levels[levelLen-1].Id {
		err = errors.New(response.MsgNotSubsequentError)
		return nil, err
	}

	// 只能进入当前关卡
	if uint32(levelId) != gameData.LevelData.CurLevel {
		err = errors.New(response.MsgPreviousError)
		return nil, err
	}
	return gameData, nil
}

func FinishTask(userId, taskId int) (err error) {

	// 判断任务是否在任务配置表中
	if _, v := taskData[taskId]; !v {
		err = errors.New(response.MsgTaskNotFoundError)
		logger.Error(err)
		return err
	}

	// 获取用户游戏数据
	gameData, err := getGameData(userId)
	if err != nil {
		return err
	}

	// 判断关卡是否在关卡表中
	if _, v := levelData[int(gameData.LevelData.CurLevel)]; !v {
		err = errors.New(response.MsgLevelNotFoundError)
		return err
	}

	// 判断该任务是否在该关卡中配置，有些关卡配置全部任务  有些不配置全部任务
	levels := viper.LevelConf.Level
	level := levels[int(gameData.LevelData.CurLevel)]
	var ok bool
	for _, v := range level.TaskList {
		if taskId == v {
			ok = true
			break
		}
	}
	if !ok {
		err = errors.New(response.MsgTaskNotFoundError)
		return err
	}

	// db更新关卡游戏数据
	gameData.LevelData.FinishTask = append(gameData.LevelData.FinishTask, uint32(taskId))
	bytes, err := proto.Marshal(gameData)
	if err != nil {
		return err
	}
	if err = data.UpdateGameData(bytes); err != nil {
		return err
	}
	// 积分更新  完成任务积分增加
	points := taskData[taskId]

	// mongo integral incr  积分增加(mongo)
	go document.AddPoints(userId, points)

	// redis integral incr  积分增加(redis)
	if err = addPoints(strconv.Itoa(userId), points); err != nil {
		return err
	}
	return nil
}

func FinishLevel(userId, levelId int) error {
	// 获取
	gameData, err := getGameData(userId)
	if err != nil {
		logger.Error(err)
		return err
	}

	// 判断任务是否全部完成
	finishTask := gameData.LevelData.FinishTask
	curLevel := gameData.LevelData.CurLevel

	taskCount := len(viper.LevelConf.Level[curLevel].TaskList)
	finishCount := len(finishTask)
	if taskCount != finishCount {
		err = errors.New(response.MsgLevelNotSuccess)
		logger.Error(err)
		return err
	}

	// 游戏数据修改当前关卡到下一关
	gameData.LevelData.CurLevel++
	gameData.LevelData.FinishTask = []uint32{}
	bytes, err := proto.Marshal(gameData)
	if err != nil {
		return err
	}
	if err := data.UpdateGameData(bytes); err != nil {
		return err
	}

	// 增加积分
	integral := levelData[levelId]

	// mongo integral incr  积分增加(mongo)
	go document.AddPoints(userId, integral)

	// 排行榜积分增加
	if err := addPoints(strconv.Itoa(userId), integral); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func CreateUserGameData() (int64, error) {
	// 存储用户游戏数据
	lastInsertId, err := data.InsertGameData(gameDataInitData)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

/*
	获取游戏数据
*/
func getGameData(userId int) (*pb.GameData, error) {
	// 数据库拿到数据
	data, err := data.GetGameData(userId)
	if err != nil {
		return nil, err
	}
	// 解析
	gameData := &pb.GameData{}
	if err := proto.Unmarshal(data.GameData, gameData); err != nil {
		return nil, err
	}
	return gameData, nil
}
