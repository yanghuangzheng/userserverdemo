package main
import(
	"net"
	"os"
	"os/signal"
	"syscall"
	"fmt"
      

	"serverdemo/initialize"
	"serverdemo/model"
	"serverdemo/grpcApi"
	"serverdemo/proto"
   

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	//"github.com/satori/go.uuid"

)

func main(){
	
	quit:=make(chan os.Signal,1)//a 用于优雅退出
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	initialize.InitLogger() //zap初始化
	//initialize.Initdb() //db初始化
	//var Port int
	_,err:=model.GetFreePort()
	if err!=nil{

		zap.S().Error("【获取空余端口失败】")
	}
	var all grpcApi.UserServer
    server := grpc.NewServer()//grpc服务初始化
	proto.RegisterUserServer(server, &all)//注册服务
	lis, err := net.Listen("tcp","127.0.0.1:8088")//监听tcp端口
	if err != nil {
		zap.S().Error("【监听tcp端口失败】")
	}else{
		zap.S().Debug("【监视成功:127.0.0.1:8088】")
	}
   
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())//注册服务健康检查
	
	cfg := api.DefaultConfig()//服务注册
	cfg.Address ="555" // Consul 服务器的地址
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.0.103:%d", 555), // 健康检查的 grpc 地址
		Timeout:                        "5s",                                   // 健康检查的超时时间
		Interval:                       "5s",                                   // 健康检查的时间间隔
		DeregisterCriticalServiceAfter: "10s",                                  // 服务在多长时间后未通过健康检查会被注销
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = "555"//global.ServerConfig.Name
	//serviceID:=fmt.Sprintf("%s",uuid.NewV4())
	registration.ID = "555"//global.ServerConfig.Name //随机生成id serviceID:=fmt.Sprintf("%s",uuid.NewV4())
	registration.Port = 555//*Port
	registration.Tags = []string{"bobby", "hhh", "user-srv"}
	registration.Address = "192.168.0.103"
	registration.Check = check
    
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	// 启动 gRPC 服务器
	go func() {
		if err := server.Serve(lis); err != nil {
			zap.S().Error("【server启动失败】")
		}
	}()
	//接受终止信号优雅退出
	<-quit
	if err = client.Agent().ServiceDeregister("/*serviceID*/"); err != nil {
		zap.S().Error("注销失败")
	}
	zap.S().Error("注销成功")
}