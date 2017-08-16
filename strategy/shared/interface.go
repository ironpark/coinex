package shared

import (
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/IronPark/coinex/strategy/proto"
	"net/rpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"strategy": &StrategyPlugin{},
}

//rpc Init(Empty) returns (Empty);
//rpc Info(Empty) returns (Information);
//
//rpc GetProperty(Empty) returns (Dictionary);
//rpc SetProperty(Dictionary) returns (Empty);
//
//rpc Update(Empty) returns (Empty);

// KV is the interface that we're exposing as a plugin.
type Strategy interface {
	Init() error
	Info() (*proto.Information,error)
	GetProperty() (*proto.Dictionary, error)
	SetProperty(*proto.Dictionary) error
	Update() (*proto.UpdateState, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type StrategyPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Strategy
}


func (p *StrategyPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return nil, nil
}

func (*StrategyPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return nil, nil
}

func (p *StrategyPlugin) GRPCServer(s *grpc.Server) error {
	proto.RegisterStrategyServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *StrategyPlugin) GRPCClient(c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewStrategyClient(c)}, nil
}