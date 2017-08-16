package shared

import (
	"github.com/IronPark/coinex/strategy/proto"
	"golang.org/x/net/context"
)

//Init(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
//Info(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Information, error)
//GetProperty(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Dictionary, error)
//SetProperty(ctx context.Context, in *Dictionary, opts ...grpc.CallOption) (*Empty, error)
//Update(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.StrategyClient }

func (m *GRPCClient) Init() error {
	_, err := m.client.Init(context.Background(), &proto.Empty{})
	return err
}
func (m *GRPCClient) Info() (*proto.Information,error) {
	resp, err := m.client.Info(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *GRPCClient) GetProperty() (*proto.Dictionary,error) {
	resp, err := m.client.GetProperty(context.Background(), &proto.Empty{})
	if err != nil {
		return &proto.Dictionary{}, err
	}

	return resp, nil
}

func (m *GRPCClient) SetProperty(dictionary *proto.Dictionary) error {
	_, err := m.client.SetProperty(context.Background(), dictionary)
	return err
}

func (m *GRPCClient) Update() (*proto.UpdateState, error) {
	return m.client.Update(context.Background(), &proto.Empty{})
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Strategy
}

func (m *GRPCServer) Init(context.Context, *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{},m.Impl.Init()
}

func (m *GRPCServer) Info(context.Context, *proto.Empty) (*proto.Information, error){
	return m.Impl.Info()
}

func (m *GRPCServer) GetProperty(context.Context, *proto.Empty) (*proto.Dictionary, error) {
	return m.Impl.GetProperty()
}

func (m *GRPCServer) SetProperty(context context.Context, property *proto.Dictionary ) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.SetProperty(property)
}

func (m *GRPCServer) Update(context.Context, *proto.Empty) (*proto.UpdateState, error) {
	return m.Impl.Update()
}

