package main

//Register alicloud packer builder plugin
import (
	"github.com/alibaba/packer-provider/alicloud-import"
	"github.com/alibaba/packer-provider/ecs"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(ecs.Builder))
	server.RegisterPostProcessor(new(alicloudimport.PostProcessor))
	server.Serve()
}
