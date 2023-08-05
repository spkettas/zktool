package tool

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	PATH   = "/data/conf/zk.conf"
)

// loadZkHost 加载zk地址
func loadZkHost() string {
	var cfg *ini.File
	err := error(nil)
	host := "argus-zk:12181"

	if cfg, err = ini.Load(PATH); err != nil {
		fmt.Printf("%v\n", err)
		return host
	}

	cfgHost := cfg.Section("ZOOKEEPER").Key("host").String()
	fmt.Printf("zk host: %v\n", host)
	if len(cfgHost) == 0 {
		return host
	}

	return cfgHost
}

// NewZk 创建zk连接
func NewZk(host string) *zk.Conn {
	if len(host) == 0 {
		host = loadZkHost()
	}

	// 配置文件是用,分隔
	hosts := strings.Split(host, ",")
	conn, _, err := zk.Connect(hosts, time.Second*5,
		zk.WithLogInfo(false))

	if err != nil {
		fmt.Printf("connect to zk err:%v\n", err)
		return nil
	}

	return conn
}

// MkDirs 创建多层目录
func MkDirs(path string, conn *zk.Conn) {
	layer := strings.Split(path, "/")
	curPath := ""

	for _, path := range layer {
		if len(path) == 0 {
			continue
		}

		curPath += "/" + path
		if flag, _, _ := conn.Exists(curPath); flag {
			continue
		}

		fmt.Printf("mk path: %s\n", curPath)
		_, err := conn.Create(curPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Printf("create path err: %v\n", err)
		}
	}
}

// SetVal 更新值
func SetVal(path, value string, conn *zk.Conn) error {
	if len(path) == 0 {
		fmt.Printf("path is null")
		return nil
	}

	// create
	if flag, _, _ := conn.Exists(path); !flag {
		MkDirs(path, conn)
	}

	// update
	_, state, _ := conn.Get(path)
	_, err := conn.Set(path, []byte(value), state.Version)
	if err != nil {
		fmt.Printf("set val err: %v\n", err)
		return err
	}

	fmt.Printf("set %v -> %v\n", path, value)
	return nil
}

// DelPath 删除路径
func DelPath(path string, conn *zk.Conn) bool {
	if len(path) == 0 {
		fmt.Printf("path is null")
		return false
	}

	if flag, _, _ := conn.Exists(path); !flag {
		fmt.Printf("not exits path: %v", path)
		return false
	}

	if err := conn.Delete(path, -1); err != nil {
		fmt.Printf("del path err: %v\n", err)
		return true
	}

	fmt.Printf("del path: %v\n", path)
	return false
}

// ShowList 遍历所有zk节点值
func ShowList(path string, conn *zk.Conn) bool {
	if flag, _, _ := conn.Exists(path); !flag {
		fmt.Printf("not exits path: %v", path)
		return false
	}

	// root
	data, stat, _ := conn.Get(path)
	fmt.Printf("%v=%v\tver=%v\n", path, string(data), stat.Version)

	// leaf
	children, _, err := conn.Children(path)
	if err != nil {
		return false
	}

	for _, p := range children {
		ShowList(path+"/"+p, conn)
	}

	return true
}
