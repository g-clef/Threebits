How to write a plugin:

 Plugins are go interfaces, that are expected to implement three methods: 
  * "Protocol" - takes no arguments and is expected to return just 
  "tcp" or "udp". This specifies which protocol the plugin applies to. 
  Threebits will create a different socket depending on whether this 
  setting is tcp or udp. 
  * "DefineArguments" - takes no arguments and does not return anything. 
  This is expected to call the "flag" go library to define any additional 
  command line arguments that the plugin would like to support (see the 
  GenericHexTCP plugin for an example of how this works).
  * "Handle" - takes arguments net.Conn, and structures.Test, returns 
  bool, string, error. This is where the actual test is run. Threebits 
  will call "Handle" with an open socket and a Test object and this 
  method is expected to do the handshake test over the provided socket 
  and return success/fail in the boolean, any nodes in the string, and 
  any errors in the error. Handle will be called once per target.
  
  
The Threebits framework will handle making and closing the socket to the 
target, so the plugin should not need to make its own connection, and 
should not close the socket. Plugins should only write to and read 
from the socket as necessary to determine if their test succeeds or 
fails, and return the corresponding bool, string, error conditions. 
 
For an example of a simple plugin, see the GenericHexTCP plugin in the Threebits-plugins-public collection.

Once you have the plugin written, import the go file in the framework's 
"plugins.go" file, and add a line in the CollectPlugins() method to 
tell the framework where your plugin is, and what it's name is. For 
example:
 
 	RegisterPlugin("http_banner", Threebits_plugins_public.HTTPBanner{})
