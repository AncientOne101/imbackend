package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"imbackend/internal/config"
	"imbackend/model"
)

type ServiceContext struct {
	// 服务配置
	Config config.Config
	// 数据库连接
	Db sqlx.SqlConn
	// 用户信息模块
	UserInfoModel model.UserInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	// 初始化用户信息模块
	return &ServiceContext{
		Config:        c,
		Db:            conn,
		UserInfoModel: model.NewUserInfoModel(conn, c.CacheRedis),
	}
}
