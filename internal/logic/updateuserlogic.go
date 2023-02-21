package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"imbackend/common/crypt"
	"imbackend/common/jwtx"
	"imbackend/internal/svc"
	"imbackend/internal/types"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateUser 更新用户信息
func (l *UpdateUserLogic) UpdateUser(req *types.UpdateRequest) (resp *types.UpdateResponse, err error) {
	fmt.Println("UpdateUser", req)
	//校验请求参数
	if req.Name == "" || req.Password == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "数据不能为空")
	}
	//判断用户是否已被删除
	user, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if user.IsDeleted == 1 {
		return nil, status.Error(codes.InvalidArgument, "该用户已被删除")
	}
	//判断请求头中是否存在token
	if req.Token == "" {
		return nil, status.Error(codes.Unauthenticated, "无操作权限,请登陆后操作1")
	}
	//校验token
	claims, err := jwtx.ParseToken(l.svcCtx.Config.Auth.AccessSecret, req.Token)
	fmt.Println("claims", claims)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "无操作权限,请登陆后操作2")
	}
	id := int64(claims["uid"].(float64))
	//判断token中的id是否与请求参数中的id一致
	if id != req.ID {
		return nil, status.Error(codes.Unauthenticated, "无操作权限,请登陆后操作3")
	}
	//获取用户信息
	userInfo, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	//对用户信息进行修改
	userInfo.Name = req.Name
	userInfo.Password = crypt.PasswordEncrypt(l.svcCtx.Config.Salt, req.Password)
	userInfo.Email = req.Email
	//将修改后的用户信息保存到数据库
	err = l.svcCtx.UserInfoModel.Update(l.ctx, userInfo)
	if err != nil {
		return nil, err
	}
	//返回数据
	return &types.UpdateResponse{
		ID:          userInfo.Id,
		Name:        userInfo.Name,
		Email:       userInfo.Email,
		Create_time: userInfo.CreateTime.String(),
		Update_time: userInfo.UpdateTime.String(),
	}, nil
}
