package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"imbackend/common/jwtx"
	"imbackend/internal/svc"
	"imbackend/internal/types"
	"imbackend/model"
	"strings"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteRequest) (resp *types.DeleteResponse, err error) {
	//参数校验
	if req.ID == 0 && req.Email == "" {
		return nil, err
	}
	//判断请求头中是否存在token
	if len(strings.TrimSpace(req.Token)) == 0 {
		return nil, status.Error(codes.Unauthenticated, "请登陆后操作")
	}
	//校验token
	claims, err := jwtx.ParseToken(l.svcCtx.Config.Salt, req.Token)
	if claims["id"] != req.ID {
		return nil, status.Error(codes.Unauthenticated, "请登陆后操作")
	}

	if err != nil {
		return nil, err
	}
	//判断请求对象is_deleted字段是否为1
	if req.IsDeleted == 1 {
		return nil, status.Error(codes.InvalidArgument, "该用户已被删除")
	}

	//查询数据库是否存在该用户
	var userInfo *model.UserInfo
	userInfo, err = l.svcCtx.UserInfoModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		return nil, err
	}
	//判断用户是否存在
	if userInfo == nil || userInfo.Id != req.ID {
		return nil, err
	}
	//删除用户,将用户的is_delete字段设置为1
	userInfo.IsDeleted = 1
	//将修改后的用户信息保存到数据库
	err = l.svcCtx.UserInfoModel.Update(l.ctx, userInfo)
	if err != nil {
		return nil, err
	}
	//返回数据
	return &types.DeleteResponse{
		ID:          int64(userInfo.Id),
		Name:        userInfo.Name,
		Email:       userInfo.Email,
		IsDeleted:   userInfo.IsDeleted,
		Create_time: userInfo.CreateTime.String(),
		Update_time: userInfo.UpdateTime.String(),
	}, nil
}
