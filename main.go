package main

//Register alicloud packer builder plugin
import (
	"github.com/alibaba/packer-provider/alicloud"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(alicloud.Builder))
	server.Serve()
}
