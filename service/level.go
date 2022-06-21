package service

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/golang/protobuf/proto"
	"sTest/entity"
	m "sTest/pkg/mysql"
	"sTest/pkg/viper"
	pb "sTest/proto"
)

func EnterLevel(levelID int, userID int) (gameData *pb.GameData, err error) {
	// 数据库拿到数据
	sqlStr := "select user_id ,game_data from g_game_data where user_id = ?"
	data := entity.GameData{}
	if err = m.DB.Get(&data, sqlStr, userID); err != nil {
		logger.Error(err)
		return nil, err
	}
	// 解析
	gameData = &pb.GameData{}
	if err = proto.Unmarshal(data.GameData, gameData); err != nil {
		logger.Error(err)
		return nil, err
	}

	// 配置表查验关卡信息
	// IO操作应该使用结构体接收一次，而不是多次
	// 大于则提示敬请期待             2-> 第一关没通关
	levelLen := len(viper.LevelConf.Level)
	if uint32(levelID) > gameData.LevelData.CurLevel {
		err = errors.New("先通过前面的关卡,才能继续挑战")
		logger.Error(err)
		return nil, err
	}
	if levelID > levelLen {
		err = errors.New("后续关卡暂未开放,敬请期待")
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
		}
	}
	if !ok {
		err = errors.New("任务不在任务列表中")
		logger.Error(err)
		return err
	}

	// 获取与判断
	gameData, err := getGameDataAndConvent(userID)
	if err != nil {
		logger.Error(err)
		return err
	}
	// 判断是否超出
	if int(gameData.LevelData.CurLevel) > len(viper.LevelConf.Level) {
		err = errors.New("敬请期待")
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
		err = errors.New("任务不在关卡任务列表中")
		logger.Error(err)
		return err
	}

	for _, val := range gameData.LevelData.FinishTask {
		logger.Info(val)
		logger.Info(taskID)
		if int(val) == taskID {
			err = errors.New("任务已完成,无需重复完成")
			logger.Error(err)
			return err
		}
	}

	// 更新
	gameData.LevelData.FinishTask = append(gameData.LevelData.FinishTask, uint32(taskID))
	bytes, err := proto.Marshal(gameData)
	if err != nil {
		logger.Error(err)
		return err
	}
	updateLevelSQL := "update g_game_data set game_data = ?"
	if _, err := m.DB.Exec(updateLevelSQL, bytes); err != nil {
		logger.Error(err)
		return err
	}

	// 积分更新  完成任务积分增加
	var account string
	findAccountSQL := "select account from t_account_data where user_id = ?"
	if err = m.DB.Get(&account, findAccountSQL, userID); err != nil {
		logger.Error(err)
		return err
	}
	var avatarURL string
	findAvatarSQL := "select avatar_url from t_base_data where user_id = ?"
	if err = m.DB.Get(&avatarURL, findAvatarSQL, userID); err != nil {
		logger.Error(err)
		return err
	}

	var integral int
	for _, val := range viper.TaskConf.Task {
		if val.Id == taskID {
			integral = val.RewardScore
		}
	}

	if err = AddIntegral(account+":"+avatarURL, integral); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func Leave(userID, levelID int) (err error) {
	// 获取
	gameData, err := getGameDataAndConvent(userID)
	if err != nil {
		logger.Error(err)
		return err
	}

	// 判断任务是否全部完成
	finishTask := gameData.LevelData.FinishTask
	levels := viper.LevelConf.Level
	if len(levels) < levelID || levelID < 1 || gameData.LevelData.CurLevel != uint32(levelID) {
		err = errors.New("关卡不正确")
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
		err = errors.New("请先完成关卡任务")
		logger.Error(err)
		return err
	}

	// 修改当前关卡到下一关
	gameData.LevelData.CurLevel++
	gameData.LevelData.FinishTask = []uint32{}
	bytes, err := proto.Marshal(gameData)
	if err != nil {
		logger.Error(err)
		return err
	}
	// 更新
	updateLevelSQL := "update g_game_data set game_data = ?"
	if _, err := m.DB.Exec(updateLevelSQL, bytes); err != nil {
		logger.Error(err)
		return err
	}

	// 增加积分
	var account string
	findAccountSQL := "select account from t_account_data where user_id = ?"
	if err = m.DB.Get(&account, findAccountSQL, userID); err != nil {
		logger.Error(err)
		return err
	}
	var avatarURL string
	findAvatarSQL := "select avatar_url from t_base_data where user_id = ?"
	if err = m.DB.Get(&avatarURL, findAvatarSQL, userID); err != nil {
		logger.Error(err)
		return err
	}

	var integral int
	for _, val := range viper.LevelConf.Level {
		if val.Id == levelID {
			integral = val.FinishReward
		}
	}

	if err = AddIntegral(account+":"+avatarURL, integral); err != nil {
		logger.Error(err)
		return err
	}

	return
}

// InitUserGameData 初始化用户游戏数据
func InitUserGameData(userID int) (ok bool, err error) {
	// 1.初始化数据
	levelInit := pb.LevelData{
		CurLevel: 1,
	}
	materialInit := pb.MaterialData{
		Warehouse: make(map[uint32]uint32),
	}
	data := pb.GameData{
		LevelData:  &levelInit,
		ShopData:   make(map[uint32]uint32),
		Statistics: &pb.StatisticsData{},
		Setting:    &pb.GameSetting{},
		Material:   &materialInit,
	}

	bytes, err := proto.Marshal(&data)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	// 2.保存到MySql
	insertGameDataSQL := "insert into g_game_data(user_id,game_data) values(?,?)"
	if _, err := m.DB.Exec(insertGameDataSQL, userID, bytes); err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}

func getGameDataAndConvent(userID int) (*pb.GameData, error) {
	// 数据库拿到数据
	sqlStr := "select user_id ,game_data from g_game_data where user_id = ?"
	data := entity.GameData{}
	if err := m.DB.Get(&data, sqlStr, userID); err != nil {
		logger.Error(err)
		return nil, err
	}
	// 解析
	gameData := &pb.GameData{}
	if err := proto.Unmarshal(data.GameData, gameData); err != nil {
		logger.Error(err)
		return nil, err
	}
	return gameData, nil
}
