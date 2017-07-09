package main

import (
	//"github.com/hashicorp/go-plugin"
	//"os/exec"
	//"os"
	//"github.com/hashicorp/go-plugin/examples/grpc/shared"
)

func LoadPlugin()  {
	// We're a host. Start by launching the plugin process.
	//client := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig: nil,
	//	Plugins:         shared.PluginMap,
	//	Cmd:             exec.Command("sh", "-c", os.Getenv("KV_PLUGIN")),
	//	AllowedProtocols: []plugin.Protocol{
	//		plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	//})
	//defer client.Kill()

	//rpcClient, _ := client.Client()
	//raw ,_:= rpcClient.Dispense("strage")
	//kv := raw.(shared.KV)
}