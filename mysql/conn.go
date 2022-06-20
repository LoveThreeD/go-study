package mysql

import (
	"fmt"
	// mysql必须导入
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func InitDB() {
	dsn := "root:123456@tcp(192.168.20.132:3306)/test?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	// 设置最大连接数
	DB.SetMaxOpenConns(20)
	// 设置空闲的最大连接数
	DB.SetMaxIdleConns(10)

	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
