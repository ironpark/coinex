package shared

import (
	"github.com/ironpark/coinex/strategy/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of Strategy that talks over RPC.
type GRPCClient struct{ client proto.StrategyClient }

func (m *GRPCClient) Init() {
	m.client.Init(context.Background(), &proto.Empty{})
}
func (m *GRPCClient) Info() (proto.Information) {
	resp, _ := m.client.Info(context.Background(), &proto.Empty{})
	return *resp
}

func (m *GRPCClient) GetProperty() (map[string]interface{}) {
	property, _ := m.client.GetProperty(context.Background(), &proto.Empty{})
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

func (m *GRPCClient) SetProperty(property map[string]interface{}) {
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
	m.client.SetProperty(context.Background(), &dic)
}

func (m *GRPCClient) SellConditions(asset string) bool {
	resp, _ := m.client.SellConditions(context.Background(), &proto.Asset{Name:asset})
	return resp.Boolean
}
func (m *GRPCClient) BuyConditions(asset string) bool {
	resp, _ := m.client.BuyConditions(context.Background(), &proto.Asset{Name:asset})
	return resp.Boolean
}
func (m *GRPCClient) RankFilter(asset string) bool {
	resp, _ := m.client.RankFilter(context.Background(), &proto.Asset{Name:asset})
	return resp.Boolean
}


// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Strategy
}


func (m *GRPCServer) Init(context.Context, *proto.Empty) (*proto.Empty, error) {
	m.Impl.Init()
	return &proto.Empty{},nil
}

func (m *GRPCServer) Info(context.Context, *proto.Empty) (*proto.Information, error){
	information := m.Impl.Info()
	return &information,nil
}

func (m *GRPCServer) GetProperty(context.Context, *proto.Empty) (*proto.Dictionary, error) {
	property := m.Impl.GetProperty()
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
	return &dic,nil
}

func (m *GRPCServer) SetProperty(context context.Context, property *proto.Dictionary ) (*proto.Empty, error) {
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
	m.Impl.SetProperty(property_map)
	return &proto.Empty{},nil
}

func (m *GRPCServer) SellConditions(ctx context.Context, asset *proto.Asset) (*proto.Bool, error) {
	return &proto.Bool{Boolean:m.Impl.SellConditions(asset.Name)},nil
}

func (m *GRPCServer) BuyConditions(ctx context.Context,asset *proto.Asset) (*proto.Bool, error) {
	return &proto.Bool{Boolean:m.Impl.BuyConditions(asset.Name)},nil
}

func (m *GRPCServer) RankFilter(ctx context.Context,asset *proto.Asset) (*proto.Bool, error) {
	return &proto.Bool{Boolean:m.Impl.RankFilter(asset.Name)},nil
}

