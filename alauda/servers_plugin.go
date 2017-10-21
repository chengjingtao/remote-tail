package alauda

import "github.com/mylxsw/remote-tail/command"
import "fmt"
import "os"
import "strings"
import "strconv"

var ServersPluginName = "alauda"

type ServersPlugin struct {
	metaData map[string]interface{}

	APIServer       string
	Token           string
	NameSpace       string
	ExecServer      string
	execServerHost  string
	execServerPort  int
	Service         string
	serviceName     string
	applicationName string
	User            string
	Password        string
}

func (sp *ServersPlugin) init(metaData map[string]interface{}) {
	sp.metaData = metaData
	value, ok := sp.metaData["api_server"]
	if !ok {
		fmt.Println("api_server required")
		os.Exit(1)
	}
	sp.APIServer = fmt.Sprint(value)

	value, ok = sp.metaData["token"]
	if !ok {
		fmt.Println("token required")
		os.Exit(1)
	}
	sp.Token = fmt.Sprint(value)

	value, ok = sp.metaData["namespace"]
	if !ok {
		fmt.Println("namespace required")
		os.Exit(1)
	}
	sp.NameSpace = fmt.Sprint(value)

	value, ok = sp.metaData["exec_server"]
	if !ok {
		fmt.Println("exec_server required")
		os.Exit(1)
	}
	sp.ExecServer = fmt.Sprint(value)
	sp.execServerHost = strings.Split(sp.ExecServer, ":")[0]
	sp.execServerPort, _ = strconv.Atoi(strings.Split(sp.ExecServer, ":")[1])

	value, ok = sp.metaData["service"]
	if !ok {
		fmt.Println("service required")
		os.Exit(1)
	}
	sp.Service = fmt.Sprint(value)
	if strings.Contains(sp.Service, "/") {
		segments := strings.Split(sp.Service, "/")
		sp.serviceName = segments[1]
		sp.applicationName = segments[0]
	} else {
		sp.serviceName = sp.Service
	}

	value, ok = sp.metaData["user"]
	if !ok {
		fmt.Println("user required")
		os.Exit(1)
	}
	sp.User = fmt.Sprint(value)

	value, ok = sp.metaData["password"]
	if !ok {
		fmt.Println("passworder required")
		os.Exit(1)
	}
	sp.Password = fmt.Sprint(value)
}

func (sp *ServersPlugin) buildServer(index int) command.Server {
	return command.Server{
		ServerName: sp.Service + "." + fmt.Sprint(index),
		Hostname:   sp.execServerHost,
		Port:       sp.execServerPort,
		User:       sp.NameSpace + "/" + sp.User,
		Password:   sp.Password,
		Script:     sp.NameSpace + "/" + sp.Service + "." + fmt.Sprint(index) + " " + " tail -f",
	}
}
func (sp *ServersPlugin) LoadServers(metaData map[string]interface{}) map[string]command.Server {
	sp.init(metaData)

	client := NewAlaudaClient(sp.APIServer, sp.Token)
	count := client.GetInstanceCount(sp.NameSpace, sp.applicationName, sp.serviceName)
	fmt.Printf("service %s count is :%d \n", sp.Service, count)

	servers := make(map[string]command.Server)
	for i := 0; i < count; i++ {
		s := sp.buildServer(i)
		key := sp.Service + "." + fmt.Sprint(i)
		servers[key] = s
	}
	return servers
}
