package plugin

import (
"github.com/IronPark/coinex/strategy/proto"
)

type StrategyBase struct{
	st Strategy
}
type Information proto.Information
type Property map[string]interface{}
type Strategy interface {
	Init()
	Info() Information
	GetProperty()Property
	SetProperty(Property)
	Update()
}

func Create(stra Strategy) StrategyBase{
	return StrategyBase{
		st:stra,
	}
}

func (bs StrategyBase) Init() error{
	bs.st.Init()
	return nil
}

func (bs StrategyBase) Info() (*proto.Information,error) {
	info := proto.Information(bs.st.Info())
	return &info,nil
}

func (bs StrategyBase) GetProperty() (*proto.Dictionary, error) {
	dic := proto.Dictionary{
		make(map[string]*proto.Property),
	}
	for k,value := range bs.st.GetProperty(){
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

func (bs StrategyBase) SetProperty(property *proto.Dictionary ) error{
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
	bs.st.SetProperty(Property(property_map))
	return nil
}

func (bs StrategyBase) Update() error{
	bs.st.Update()
	return nil
}
