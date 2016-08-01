package Threebits

import (
	"flag"
	"fmt"
	"sync"
	"os"
)

func Run(){
	CollectPlugins()
	var IP = flag.String("IP", "0.0.0.0", "IP of Server")
	var port = flag.Int("port", 8000, "Port of Server")
	var mode = flag.String("mode", "l", "mode to run in: m=Manager, w=worker, l=list plugins")
	var targetlist = flag.String("targets", "targetlist.txt", "file containing the targets to scan, one per line (only used in manager mode)")
	var whitelist = flag.String("whitelist", "whitelist.txt", "file containing a list of IPs to skip even if they are in the target list (only used in manager mode)")
	var scansToRun = flag.String("scans", "", "comma-separated list of commands to run (only used in manager mode)")
	var output = flag.String("output", "output.csv", "where to send output. if specified, will write to that filename (only used in manager mode)")
	var quiet = flag.Bool("quiet", true, "quiet mode: if false, will record falses as well as trues (only used in manager mode)")
	var numWorkers = flag.Int("numworkers", 5000, "number of workers to run (Only used in worker mode)")
	var timeout = flag.Int("timeout", 3, "socket timeout in seconds (only used in worker mode)")
	var wg sync.WaitGroup

	flag.Parse()

	switch *mode {
		case "m":
			// bail out if no actual scans are given.
			if len(*scansToRun) == 0{
				fmt.Println("No scans given to run")
				os.Exit(1)
			}
			wg.Add(1)
			RunManager(*IP, *port, *targetlist, *whitelist, *scansToRun, *output, *quiet, &wg)
		case "w":
			wg.Add(1)
			RunWorkers(*IP, *port, *numWorkers, *timeout, &wg)
		case "l":
			fmt.Println("Plugins:")
			for plugin := range Plugins{
				fmt.Println("\t" + plugin)
			}
		default: fmt.Println("unknown mode")
	}
	wg.Wait()
}
