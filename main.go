package main

import (
	"flag"
	"zktool/tool"
)

var (
	host  string
	path  string
	value string
	cmd   string
)

func init() {
	flag.StringVar(&host, "z", "", "zookeeper host, multiple split by comma")
	flag.StringVar(&path, "path", "", "zk path")
	flag.StringVar(&value, "val", "", "zk val")
	flag.StringVar(&cmd, "c", "help", "cmd options: help, list, create, del")
}

// / main
func main() {
	flag.Parse()

	conn := tool.NewZk(host)
	defer conn.Close()

	switch cmd {
	case "create":
		tool.SetVal(path, value, conn)
	case "del":
		tool.DelPath(path, conn)
	case "list":
		tool.ShowList(path, conn)
	case "help":
		flag.Usage()
	default:
		flag.Usage()
	}
}
