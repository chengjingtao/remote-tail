package command

type Server struct {
	ServerName     string `toml:"server_name"`
	Hostname       string `toml:"hostname"`
	Port           int    `toml:"port"`
	User           string `toml:"user"`
	Password       string `toml:"password"`
	PrivateKeyPath string `toml:"private_key_path"`
	TailFile       string `toml:"tail_file"`
	Script         string
}

type Config struct {
	TailFile      string            `toml:"tail_file"`
	Servers       map[string]Server `toml:"servers"`
	Slient        bool              `toml:"slient"`
	ServersPlugin ServersPlugin     `toml:"servers_plugin"`
}

type ServersPlugin struct {
	Name     string                 `toml:"name"`
	MetaData map[string]interface{} `toml:"meta_data"`
	action   ServersPluginAction
}

type ServersPluginAction interface {
	LoadServers(metaData map[string]interface{}) map[string]Server
}

func (sp *ServersPlugin) Init() {
	action := ServersPlugins[sp.Name]
	sp.action = action
}
func (sp *ServersPlugin) LoadServers() map[string]Server {
	return sp.action.LoadServers(sp.MetaData)
}

var ServersPlugins map[string]ServersPluginAction = make(map[string]ServersPluginAction)
