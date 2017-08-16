package strategy

import (
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/IronPark/coinex/strategy/shared"
	"github.com/IronPark/coinex/strategy/proto"
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
func (st *Strategy) Init() {
	st.strategy.Init()
}

func (st *Strategy) Info() proto.Information {
	info ,_ := st.strategy.Info()
	return *info
}

func (st *Strategy) GetProperty() map[string]interface{} {
	property ,_ := st.strategy.GetProperty()
	property_map := make(map[string]interface{})
	for k,v := range property.CustomInt{
		switch v.Type {
		case "int":
			property_map[k] = v.ValueInt
		case "float":
			property_map[k] = v.ValueFloat
		case "string":
			property_map[k] = v.ValueString
		case "bool":
			property_map[k] = v.ValueBool
		}
	}
	return property_map
}

func (st *Strategy) SetProperty(property map[string]interface{}) {
	dic := proto.Dictionary{
		make(map[string]*proto.Property),
	}
	for k,value := range property{
		switch v := value.(type) {
		case int32:
			dic.CustomInt[k] = &proto.Property{ValueInt:int32(v)}
		case float32:
			dic.CustomInt[k] = &proto.Property{ValueFloat:float32(v)}
		case string:
			dic.CustomInt[k] = &proto.Property{ValueString:string(v)}
		case bool:
			dic.CustomInt[k] = &proto.Property{ValueBool:bool(v)}
		}
	}
	st.strategy.SetProperty(&dic)
}

func (st *Strategy) Update() bool {
	b,_ := st.strategy.Update()
	return b.Boolean
}

func (st *Strategy) KillProcess() {
	st.rpc.Kill()
}