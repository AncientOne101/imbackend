package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"imbackend/common/crypt"
	"imbackend/common/jwtx"
	"time"

	"imbackend/internal/svc"
	"imbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 用户登录
func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	//先校验请求数据
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "请输入数据再登录")
	}
	//判断用户邮箱是否已经注册
	user, err := l.svcCtx.UserInfoModel.FindOneByEmail(l.ctx, req.Email)
	if user == nil {
		return nil, status.Error(codes.InvalidArgument, "邮箱未注册,请先注册")
	}
	//判断用户是否被删除
	if user.IsDeleted == 1 {
		return nil, status.Error(codes.InvalidArgument, "该用户已被删除")
	}
	//将密码加密,判断密码是否正确
	if user.Password != crypt.PasswordEncrypt(l.svcCtx.Config.Salt, req.Password) {
		return nil, status.Error(codes.InvalidArgument, "密码错误,请重新输入")
	}
	//获取当前时间戳
	now := time.Now().Unix()
	//获取token过期时间
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	//生成token
	token, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, user.Id)
	if err != nil {
		return nil, err
	}
	//解析token
	claims, err := jwtx.ParseToken(l.svcCtx.Config.Auth.AccessSecret, token)
	userid := claims["uid"]
	fmt.Println("claims", claims)
	fmt.Println("userid", userid)
	//返回数据和token
	return &types.LoginResponse{
		ID:     int64(user.Id),
		Name:   user.Name,
		Email:  user.Email,
		Token:  token,
		Expire: now + accessExpire,
	}, nil
}
