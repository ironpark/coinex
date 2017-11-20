package strategy

import (
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/ironpark/coinex/strategy/shared"
	"github.com/ironpark/coinex/strategy/proto"
)
type Strategy struct {
	rpc *plugin.Client
	strategy shared.Strategy
}

func LoadStrategy(path string) (*Strategy){
	 //We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", path), //file path (binary or script)
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	rpcClient, _ := client.Client()
	raw ,_:= rpcClient.Dispense("strategy")
	st := raw.(shared.Strategy)
	return &Strategy{client,st}
}

func LoadStrategys(path string) (*Strategy){
	//We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", path), //file path (binary or script)
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	rpcClient, _ := client.Client()
	raw ,_:= rpcClient.Dispense("strategy")
	st := raw.(shared.Strategy)
	return &Strategy{client,st}
}

func (st *Strategy) Init() {
	st.strategy.Init()
}

func (st *Strategy) Info() proto.Information {
	return st.strategy.Info()
}

func (st *Strategy) GetProperty() map[string]interface{} {
	return st.strategy.GetProperty()
}

func (st *Strategy) SetProperty(property map[string]interface{}) {
	st.strategy.SetProperty(property)
}

func (st *Strategy) SellConditions(name string) bool {
	return st.strategy.SellConditions(name)
}

func (st *Strategy) BuyConditions(name string) bool {
	return st.strategy.BuyConditions(name)
}

func (st *Strategy) RankFilter(name string) bool {
	return st.strategy.RankFilter(name)
}
func (st *Strategy) KillProcess() {
	st.rpc.Kill()
}