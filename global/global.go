package global
//前置条件获取数据库的DB
import(
	"gorm.io/gorm"
)
var(
	DB*gorm.DB
)