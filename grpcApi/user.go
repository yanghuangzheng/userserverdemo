package grpcApi

import (
	"context"
	"time"
	"errors"
	
	"serverdemo/form"
	"serverdemo/proto"
	"serverdemo/global"
	"serverdemo/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)
type UserServer struct{}


func (s *UserServer) GetUserList(con context.Context, in *proto.PageInfo) (*proto.UserListRespons, error) {
	var users[] form.User
    result:=global.DB.Find(&users)
	if result.Error !=nil{
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 没有找到记录
			zap.S().Info("【 GetUserList: 没有找到对应 mobile 的记录】")
			return nil, gorm.ErrRecordNotFound
		} else {
			// 数据库操作出错
			zap.S().Errorw("【 GetUserList的db操作出错】", "error", result.Error)
			return nil, result.Error
		}
	}
   rsp:=&proto.UserListRespons{}
   rsp.Total=int32(result.RowsAffected)
   //分页 gorm limt方法
   global.DB.Scopes(Paginate(int(in.Pn),int(in.Psize))).Find(&users)  
   for _,user:=range users{
	userInfoRsp:=modelToResponse(user)
	rsp.Data=append(rsp.Data,userInfoRsp)
   }
   zap.S().Info("【GetUserList操作ok】")
return rsp,nil
}
////////////////////////////////////////////////////////////////////////////
func (s *UserServer) SelectByMobile(con context.Context, in *proto.MobileInfo) (*proto.UserRespons, error) {
    var  user=form.User{}
	result := global.DB.Where("mobile = ?", in.Mobile).First(&user)
        if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// 没有找到记录
				zap.S().Info("【SelectByMobile: 没有找到对应 mobile 的记录】")
				return nil, gorm.ErrRecordNotFound
			} else {
				// 数据库操作出错
				zap.S().Errorw("【SelectByMobile的db操作出错】", "error", result.Error)
				return nil, result.Error
			}
        }
    // 构建响应
    rsp := &proto.UserRespons{
        Id:       user.ID,
        Moblie:   user.Mobile,
        Name:     user.Name,
        Birthday: uint64(user.Birthday.Unix()),
        Gender:   user.Gender,
        Role:     int32(user.Role),
    }
    return rsp, nil
}
//////////////////////////////////////////////////////////////////////////
func (s *UserServer) SelectById(con context.Context, in *proto.IdInfo) (*proto.UserRespons, error) {
	var  user=form.User{}
	result := global.DB.Where("Id = ?", in.Id).First(&user)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					// 没有找到记录
					zap.S().Info("【SelectById: 没有找到对应 mobile 的记录】")
					return nil, gorm.ErrRecordNotFound
				} else {
					// 数据库操作出错
					zap.S().Errorw("【SelectById的db操作出错】", "error", result.Error)
					return nil, result.Error
				}
			}
		// 构建响应
		rsp := &proto.UserRespons{
			Id:       user.ID,
			Moblie:   user.Mobile,
			Name:     user.Name,
			Birthday: uint64(user.Birthday.Unix()),
			Gender:   user.Gender,
			Role:     int32(user.Role),
		}
		return rsp, nil
	}
//////////////////////////////////////////////////////////////////////////
func (s *UserServer) UpdateUser(con context.Context, in *proto.UserInfo) (*proto.SuccessResponse, error) {
	birthday:=time.Unix(int64(in.Birthday), 0)
	  user:=form.User{}
	result := global.DB.Where("Id = ?", in.Id).First(&user)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					// 没有找到记录
					zap.S().Info("【UpdateUser: 没有找到对应 users 的记录】")
					return nil, gorm.ErrRecordNotFound
				} else {
					// 数据库操作出错
					zap.S().Errorw("【UpdateUser的find db操作出错】", "error", result.Error)
					return nil, result.Error
				}
			}
	Password:=model.Md5(user.Salt+in.Password)
	var  user1=form.User{
		Password:  Password,
		Mobile:    in.Moblie,
		Name:      in.Name,
		Birthday:  &birthday,
		Gender :   in.Gender,
	}
	result1 := global.DB.Model(&form.User{}).Where("id = ?", in.Id).Updates(user1)
    if result1.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 没有找到记录
			zap.S().Info("【UpdateUser: 没有找到对应 mobile 的记录】")
			return nil, gorm.ErrRecordNotFound
		} else {
			// 数据库操作出错
			zap.S().Errorw("【UpdateUser的db操作出错】", "error", result.Error)
			return nil, result.Error
    } 
		}
	rsp := &proto.SuccessResponse{
		Success: "yes",
	}
	return rsp, nil
}
//////////////////////////////////////////////////////////////////////////
func (s *UserServer) CreateUser(con context.Context, in *proto.UserInfo) (*proto.SuccessResponse, error) {
	birthday:=time.Unix(int64(in.Birthday), 0)
     salt,err:=model.GenerateSalt(10)
	 if err!=nil{
		zap.S().Info("【GenerateSalt(10)操作出错】")
	 }
	 Password:=model.Md5(salt+in.Password)
	var  user=form.User{
		ID      :  model.GenerateRandomNumber(),
		Password:  Password,
		Mobile:    in.Moblie,
		Name:      in.Name,
		Birthday:  &birthday,
		Gender :   in.Gender,
        Salt:      salt,
	}
	error:=global.DB.Save(&user)
	if err!=nil {
		zap.S().Info("【CreateUser的db操作出错】")
				return nil,error.Error
	}
	rsp := &proto.SuccessResponse{
		Success: "yes",
	}
	return rsp, nil
}
//////////////////////////////////////////////////////////////////////////
func (s *UserServer) Logger(con context.Context, in *proto.PasswordInfo) (*proto.IdResponse, error) {
	rsp := &proto.IdResponse{
		Success: "no",
	}
	var  user=form.User{}
	result := global.DB.Where("Id = ?", in.Id).First(&user)
			if result.Error != nil {
				zap.S().Info("【 无此用户】")
				return nil,result.Error
			}
			if user.Password!=model.Md5(user.Salt+in.Password) {
				return rsp, nil
			}
	rsp = &proto.IdResponse{
		Success: "yes",
		Id: user.ID,
		Role:int32(user.Role),
	}
	return rsp, nil
}
