package main

// import plugin so it will be built with Tarsier
// TODO this file can be generated with some build script
import (
	_ "github.com/sejvlond/tarsier/plugins/dummy"
	_ "github.com/sejvlond/tarsier/plugins/faults"
	_ "github.com/sejvlond/tarsier/plugins/heavy_load"
	_ "github.com/sejvlond/tarsier/plugins/net"
	_ "github.com/sejvlond/tarsier/plugins/persistent_storage"
	_ "github.com/sejvlond/tarsier/plugins/sleeping_beauty"
)
