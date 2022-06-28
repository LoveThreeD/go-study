package data

import (
	"errors"
	"sTest/entity"
	m "sTest/pkg/mysql"
)

/*
	更新游戏数据
*/
func UpdateGameData(bytes []byte) error {
	query := "update g_game_data set game_data = ?"
	result, err := m.DB.Exec(query, bytes)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 0 {
		return errors.New("update fail")
	}
	return nil
}

/*
	插入游戏数据
*/
func InsertGameData(bytes []byte) (int64, error) {
	insertGameDataSQL := "insert into g_game_data(game_data) values(?)"
	result, err := m.DB.Exec(insertGameDataSQL, bytes)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

/*
	获取用户的gameData数据
*/
func GetGameData(userId int) (*entity.GameData, error) {
	sqlStr := "select user_id ,game_data from g_game_data where user_id = ?"
	gameData := &entity.GameData{}
	if err := m.DB.Get(gameData, sqlStr, userId); err != nil {
		return nil, err
	}
	return gameData, nil
}
