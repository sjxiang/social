package config

import (
	"testing"
)


func TestNewConfig(t *testing.T) {
	
	cfg, err := New()
	if err!= nil {
		t.Fatal(err)
	}

	t.Logf("配置清单, %+v", cfg)
}

func TestDeaultEnv(t *testing.T) {

	url := defaultEnvString("URL", "http://www.baidu.com")
	if url!= "URL_ADDRESS" {
		t.Fatal("配置清单, 配置文件中没有 URL 环境变量")
	}

	t.Log("测试正确")
}