package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"imbackend/common/crypt"
	"imbackend/internal/svc"
	"imbackend/internal/types"
	"imbackend/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	//日志组件
	logx.Logger
	//上下文组件
	ctx context.Context
	//服务组件
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 用户注册
func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	//先校验请求数据
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "用户名或邮箱或密码不能为空")
	}
	//判断用户邮箱是否已经注册
	user, err := l.svcCtx.UserInfoModel.FindOneByEmail(l.ctx, req.Email)
	if user != nil {
		return nil, status.Error(codes.InvalidArgument, "邮箱已经注册,请直接登录")
	}
	//判断用户名是否已经存在,如果存在,则需要重新输入用户名
	user, err = l.svcCtx.UserInfoModel.FindOneByName(l.ctx, req.Name)
	if user != nil && user.Name != "" {
		return nil, status.Error(codes.InvalidArgument, "用户名已被使用,请重新输入用户名")
	}

	//将参数封装成newUser
	userInfo := model.UserInfo{
		Name:  req.Name,
		Email: req.Email,
		//密码加密
		Password: crypt.PasswordEncrypt(l.svcCtx.Config.Salt, req.Password),
	}
	//将用户信息保存到数据库
	_, err = l.svcCtx.UserInfoModel.Insert(l.ctx, &userInfo)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	//从数据库中查询用户信息
	user, err = l.svcCtx.UserInfoModel.FindOneByName(l.ctx, req.Name)
	//返回响应
	return &types.RegisterResponse{
		ID:    int64(user.Id),
		Name:  user.Name,
		Email: user.Email,
	}, nil

}
