package initialize
import(
	"log"
	"os"
	"time"

	"serverdemo/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"go.uber.org/zap"
)

func Initdb()  {
	dsn:="sql路径"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		  SlowThreshold:              time.Second,   // Slow SQL threshold
		  LogLevel:                   logger.Info, // Log level
		  IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
		  ParameterizedQueries:      true,           // Don't include params in the SQL log
		  Colorful:                  true,          // Disable color
		},)
		var err error
		global.DB,err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		  })
		  if err!=nil{
			zap.S().Error("【Initdb】启动失败")
		}
	}
	//启动配置界面