package Threebits

/*
* Edit this file to include your plugins into the full binary.
* To do this, first include the plugin path to the plugins
* in the "import" section.
*
* Then, add the plugins (as either TCP or UDP) to the "CollectPlugins"
* function with a name. Please make sure that the plugin has a unique
* name. The name you give it will also be given to users on the command
* line if they ask the system to list al plugins.
*
* For more details on writing plugins, see the documentation.
*
 */

import (
	// add other private repos of plugins here
	"errors"
	"github.com/accessviolationsec/Threebits/structures"
	"github.com/accessviolationsec/Threebits-plugins-public"
	"net"
)

func CollectPlugins(){
	RegisterPlugin("http_banner", Threebits_plugins_public.HTTPBanner{})
	RegisterPlugin("https_banner", Threebits_plugins_public.HTTPSBanner{})
	RegisterPlugin("ssh_banner", Threebits_plugins_public.SSHBanner{})
	// add other RegisterPlugin lines here to register your
	// plugins with the system.
}

func RegisterPlugin(pluginName string, target Plugin) error {
	if _, ok := Plugins[pluginName]; ok{
		return errors.New("Plugin already exists")
	}
	Plugins[pluginName] = target
	return nil
}


type Plugin interface {
	Handle(socket net.Conn, test structures.Test) (bool, string, error)
	Protocol()(string)
}


var Plugins = make(map[string]Plugin)
