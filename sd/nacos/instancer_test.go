// author: JGZ
// time:   2021-01-28 16:38
package nacos

import (
	"fmt"
	"github.com/jjggzz/kit/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
	"time"
)

func Test_getInstancer(t *testing.T) {
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId("73d8542f-4ccb-4f16-af3b-c104bee8af8f"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("D:\\tmp\\nacos\\log"),
		constant.WithCacheDir("D:\\tmp\\nacos\\cache"),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel("debug"),
	)
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			"192.168.151.109",
			8848,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	s := NewInstancer(namingClient, log.NewNopLogger(), "testNacos", "", nil)
	defer s.Stop()
	for {
		state := s.cache.State()
		for _, e := range state.Instances {
			fmt.Println(e)
		}
		fmt.Println("-----------------------------")
		time.Sleep(time.Second * 3)
	}

}
