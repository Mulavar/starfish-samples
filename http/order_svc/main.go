package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotrx/mysql"
	"github.com/gotrx/starfish/pkg/client"
	"github.com/gotrx/starfish/pkg/client/config"

	"github.com/gotrx/starfish-samples/order_svc/dao"
)

const configPath = "conf/client.yml"

func main() {
	r := gin.Default()
	config.InitConf(configPath)
	client.NewRpcClient()
	mysql.InitDataResourceManager()
	mysql.RegisterResource(config.GetATConfig().DSN)

	sqlDB, err := sql.Open("mysql", config.GetATConfig().DSN)
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(4 * time.Hour)

	d := &dao.Dao{
		DB: sqlDB,
	}

	r.POST("/createSo", func(c *gin.Context) {
		type req struct {
			Req []*dao.SoMaster
		}
		var q req
		if err := c.ShouldBindJSON(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := d.CreateSO(
			context.WithValue(
				context.Background(),
				mysql.XID,
				c.Request.Header.Get("XID")),
			q.Req)

		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "fail",
			})
		} else {
			c.JSON(200, gin.H{
				"success": true,
				"message": "success",
			})
		}
	})
	r.Run(":8002")
}
