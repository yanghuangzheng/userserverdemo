package nacos

import (
	"time" // 导入 time 包，用于处理时间相关的功能
	"fmt"  // 导入 fmt 包，用于格式化 I/O 操作
	//"encoding/json"

	// "github.com/nacos-group/nacos-sdk-go/v2" // 这个导入被注释掉了，可能是备用或者未来的版本
	"github.com/nacos-group/nacos-sdk-go/clients" // 导入 Nacos SDK 的客户端包
	"github.com/nacos-group/nacos-sdk-go/common/constant" // 导入 Nacos SDK 的常量包
	"github.com/nacos-group/nacos-sdk-go/vo" // 导入 Nacos SDK 的值对象 (Value Object) 包
)

func nacos() {
	// 创建 Nacos 服务器配置
	sc := []constant.ServerConfig{
		{
			IpAddr: "192.168.1.103", // Nacos 服务器的 IP 地址
			Port:   8848,            // Nacos 服务器的端口号
		},
	}

	// 创建 Nacos 客户端配置
	cc := constant.ClientConfig{
		NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // Nacos 命名空间 ID
		TimeoutMs:           5000,                                   // 请求超时时间（毫秒）
		NotLoadCacheAtStart: true,                                   // 启动时不加载缓存
		LogDir:              "tmp/nacos/log",                         // 日志目录
		CacheDir:            "tmp/nacos/cache",                       // 缓存目录
		LogLevel:            "debug",                                 // 日志级别
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,  // 传入服务器配置
		"clientConfig":  cc,  // 传入客户端配置
	})
	if err != nil {
		panic(err) // 如果创建客户端失败，抛出异常
	}

	// 获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-web.yaml", // 配置的数据 ID
		Group:  "dev",           // 配置的组名
	})
	if err != nil {
		panic(err) // 如果获取配置失败，抛出异常
	}
	fmt.Println(content) // 打印获取到的配置内容

	// 使程序暂停 3000 秒（50 分钟），以便让监听器有足够的时间来检测配置变化
	time.Sleep(3000 * time.Second)


	//json.Unmarshal()//将json转换
	//serverConfigs:=config.ServerConfig{}
	//想要将一个字符串装换成struct需要去设置这个struct的tag
	//json.Unmarshal([]byte(content),&serverConfigs)
	//fmt.Println(serverConfigs)



	// 监听配置变化
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web.yaml", // 配置的数据 ID
		Group:  "dev",           // 配置的组名
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data) // 当配置发生变化时，打印新的配置内容
		},
	})
	if err != nil {
		panic(err) // 如果监听配置失败，抛出异常
	}
}
