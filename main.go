package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/wanggaolin/go_lib/w"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	flsg_arry             []string
	zookeep_config        string // 配置文件
	zookeep_dataDir       string // 数据目录
	zookeep_max_node      int    // 最大节点数量
	zookeep_stateful_name string
	zookeep_serviceName   string
	zk                    *zookeeper
)

type zookeeper struct {
}

func (z *zookeeper) get_args() (err error) {
	filter_args, _ := regexp.Compile(`^--.+=.*`)
	filter_config, _ := regexp.Compile(`^--config=.+`)
	filter_dataDir, _ := regexp.Compile(`^--dataDir=.+`)
	filter_max_node, _ := regexp.Compile(`^--max-node=.+`)
	filter_serviceName, _ := regexp.Compile(`^--server=.+`)
	filter_StatefulSet, _ := regexp.Compile(`^--stateful_name=.+`)
	for _, item := range os.Args[1:] {
		if filter_args.Match([]byte(item)) {
			match_config := filter_config.FindStringSubmatch(item)
			match_dataDir := filter_dataDir.FindStringSubmatch(item)
			match_max_node := filter_max_node.FindStringSubmatch(item)
			match_serviceName := filter_serviceName.FindStringSubmatch(item)
			match_StatefulSet := filter_StatefulSet.FindStringSubmatch(item)
			if len(match_config) > 0 {
				zookeep_config = strings.Split(match_config[0], "=")[1]
			} else if len(match_dataDir) > 0 {
				zookeep_dataDir = strings.Split(match_dataDir[0], "=")[1]
				flsg_arry = append(flsg_arry, strings.Trim(item, "--"))
			} else if len(match_max_node) > 0 {
				zookeep_max_node, _ = strconv.Atoi(strings.Split(match_max_node[0], "=")[1])
			} else if len(match_serviceName) > 0 {
				zookeep_serviceName = strings.Split(match_serviceName[0], "=")[1]
			} else if len(match_StatefulSet) > 0 {
				zookeep_stateful_name = strings.Split(match_StatefulSet[0], "=")[1]
			} else {
				flsg_arry = append(flsg_arry, strings.Trim(item, "--"))
			}
		}
	}
	if zookeep_config == "" {
		err = errors.New("use --config to specify the configuration file")
		return err
	} else if zookeep_dataDir == "" {
		err = errors.New("use --dataDir to specify the data store directory")
		return err
	} else if zookeep_max_node == 0 {
		err = errors.New("use --dataDir to specify the max node size,require int")
		return err
	} else if zookeep_serviceName == "" {
		err = errors.New("use --server to specify the data store directory")
		return err
	} else if zookeep_stateful_name == "" {
		err = errors.New("use --stateful_name to specify the data store directory")
		return err
	}
	return
}

func init() {
	zk = &zookeeper{}
}

func (z *zookeeper) create_myid() (err error) {
	myid := ""
	if w.File.PathExists(zookeep_dataDir) == false {
		if err = os.MkdirAll(zookeep_dataDir, 0755); err != nil {
			return err
		}
	}
	myid_file := path.Join(zookeep_dataDir, "myid")
	file, err := os.OpenFile(myid_file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	hostarry := strings.Split(os.Getenv("HOSTNAME"), "-")
	myid = hostarry[len(hostarry)-1]
	statefull_inex, _ := strconv.Atoi(myid)
	myid = fmt.Sprintf("%d", statefull_inex+1)
	write := bufio.NewWriter(file)
	write.WriteString(myid)
	write.Flush()
	return
}

func (z *zookeeper) create_config() (err error) {
	file, err := os.OpenFile(zookeep_config, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for i := 0; i < zookeep_max_node; i++ {
		index := i + 1
		server_node := fmt.Sprintf("server.%d=%s-%d.%s:2888:3888", index, zookeep_stateful_name, i, zookeep_serviceName)
		flsg_arry = append(flsg_arry, server_node)
	}
	write := bufio.NewWriter(file)
	write.WriteString(strings.Join(flsg_arry, "\n"))
	write.Flush()
	return
}

func run() (err error) {
	if err = zk.get_args(); err != nil {
		return err
	}
	if err = zk.create_myid(); err != nil {
		return err
	}
	if err = zk.create_config(); err != nil {
		return err
	}
	return err
}

func main() {
	if err := run(); err != nil {
		w.Log.Log_error(err)
		w.ExitError()
	}
}
