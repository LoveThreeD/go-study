package service

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/golang/protobuf/proto"
	"sTest/pkg/response"
	"sTest/pkg/viper"
	pb "sTest/proto"
	"sTest/repository/data"
	"sTest/repository/document"
	"strconv"
)

var (
	// gameData的默认初始数据,用来初始化用户的游戏数据
	gameDataInitData []byte
)

func init() {
	var err error
	// 1.初始化数据
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
}

func EnterLevel(levelID int, userID int) (gameData *pb.GameData, err error) {
	// 获取与判断
	gameData, err = getGameData(userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// 配置表查验关卡信息
	// IO操作应该使用结构体接收一次，而不是多次
	// 大于则提示敬请期待             2-> 第一关没通关
	levelLen := len(viper.LevelConf.Level)
	if uint32(levelID) > gameData.LevelData.CurLevel {
		err = errors.New(response.MsgPreviousError)
		logger.Error(err)
		return nil, err
	}
	if levelID > levelLen {
		err = errors.New(response.MsgNotSubsequentError)
		logger.Error(err)
		return nil, err
	}
	return
}

func MissionAccomplished(userID, taskID int) (err error) {
	// 判断任务是否在配置表中
	var ok bool
	for _, val := range viper.TaskConf.Task {
		if val.Id == taskID {
			ok = true
			break
		}
	}
	if !ok {
		err = errors.New(response.MsgTaskNotFoundError)
		logger.Error(err)
		return err
	}

	// 获取与判断
	gameData, err := getGameData(userID)
	if err != nil {
		logger.Error(err)
		return err
	}
	// 判断是否超出
	if int(gameData.LevelData.CurLevel) > len(viper.LevelConf.Level) {
		err = errors.New("stay tuned")
		logger.Error(err)
		return err
	}

	// 判断任务列表
	ok = false
	for _, val := range viper.LevelConf.Level[int(gameData.LevelData.CurLevel)-1].TaskList {
		if val == taskID {
			ok = true
		}
	}
	if !ok {
		err = errors.New(response.MsgTaskNotFoundError)
		logger.Error(err)
		return err
	}

	for _, val := range gameData.LevelData.FinishTask {
		if int(val) == taskID {
			err = errors.New(response.MsgTaskRepeatError)
			logger.Error(err)
			return err
		}
	}

	// 更新
	gameData.LevelData.FinishTask = append(gameData.LevelData.FinishTask, uint32(taskID))
	bytes, err := proto.Marshal(gameData)
	if err != nil {
		return err
	}
	if err = data.UpdateGameData(bytes); err != nil {
		return err
	}
	// 积分更新  完成任务积分增加
	var integral int
	for _, val := range viper.TaskConf.Task {
		if val.Id == taskID {
			integral = val.RewardScore
			break
		}
	}

	// mongo integral incr  积分增加(mongo)
	go document.IncrIntegral(userID, integral)

	// redis integral incr  积分增加(redis)
	if err = AddIntegral(strconv.Itoa(userID), integral); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func Leave(userID, levelID int) (err error) {
	// 获取
	gameData, err := getGameData(userID)
	if err != nil {
		logger.Error(err)
		return err
	}

	// 判断任务是否全部完成
	finishTask := gameData.LevelData.FinishTask
	levels := viper.LevelConf.Level
	if len(levels) < levelID || levelID < 1 || gameData.LevelData.CurLevel != uint32(levelID) {
		err = errors.New(response.MsgLevelChooseError)
		logger.Error(err)
		return err
	}

	var levelTaskTotal int
	var completed int
	for _, val := range viper.LevelConf.Level[levelID-1].TaskList {
		levelTaskTotal += val
	}
	for _, val := range finishTask {
		completed += int(val)
	}
	if levelTaskTotal != completed {
		err = errors.New(response.MsgLevelNotSuccess)
		logger.Error(err)
		return err
	}

	// 修改当前关卡到下一关
	gameData.LevelData.CurLevel++
	gameData.LevelData.FinishTask = []uint32{}
	bytes, err := proto.Marshal(gameData)
	if err != nil {
		return err
	}
	// 更新
	if err := data.UpdateGameData(bytes); err != nil {
		return err
	}

	// 增加积分
	var integral int
	for _, val := range viper.LevelConf.Level {
		if val.Id == levelID {
			integral = val.FinishReward
		}
	}

	// mongo integral incr  积分增加(mongo)
	go document.IncrIntegral(userID, integral)

	if err = AddIntegral(strconv.Itoa(userID), integral); err != nil {
		logger.Error(err)
		return err
	}

	return
}

// InitUserGameData 初始化用户游戏数据
func InitUserGameData() (userId int64, err error) {
	// 存储用户游戏数据
	lastInsertId, err := data.InsertGameData(gameDataInitData)
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
