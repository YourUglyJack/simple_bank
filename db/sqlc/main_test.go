package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"simple_bank/util"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// 在一个go包里包含所有单元测试
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config, err:", err)
	}
	
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("connoy connect to db:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := testDB.PingContext(ctx); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	// 把 db 连接包装成 sqlc 的查询对象
	testQueries = New(testDB)
	os.Exit(m.Run())
}
