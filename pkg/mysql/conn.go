package mysql

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"sTest/pkg/viper"

	// mysql必须导入
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func init() {
	//dsn := "root:123456@tcp(192.168.20.132:3306)/test?charset=utf8mb4&parseTime=True"
	c := viper.Conf.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.Username, c.Password, c.Address, c.Port, c.DbName, c.URL)
	// 也可以使用MustConnect连接不成功就panic
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Fatalf("connect DB failed, err:%v\n", err)
	}
	// 设置最大连接数
	DB.SetMaxOpenConns(20)
	// 设置空闲的最大连接数
	DB.SetMaxIdleConns(10)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
