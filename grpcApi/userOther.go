package grpcApi

import(

	"serverdemo/form"
	"serverdemo/proto"

	"gorm.io/gorm"
)
func modelToResponse(user form.User) *proto.UserRespons { //将数据库转换为proto文件//在grpc的message中字段有默认值不能随便赋值nil进去，容易出错 要搞清楚 哪些字段有默认值
	UserInfoRsp := proto.UserRespons{
		Id:     user.ID,
		Name:   user.Name,
		Gender: user.Gender,
		Moblie: user.Mobile,
		Role:   int32(user.Role),
	}
	if user.Birthday != nil {
		UserInfoRsp.Birthday = uint64(user.Birthday.Unix()) //转换
	}
	return &UserInfoRsp
}
func Paginate(page ,pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
	  if page <= 0 {
		page = 1
	  }
	  switch {
	  case pageSize > 100:
		pageSize = 100
	  case pageSize <= 0:
		pageSize = 10
	  }
  
	  offset := (page - 1) * pageSize
	  return db.Offset(offset).Limit(pageSize)
	}
  }
