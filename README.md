Threebits

 Threebits is a scanner designed to look for handshakes across many IPs in parallel. 

 Threebits runs in two modes: Manager and Worker. For a given scan, you need one Manager, but can run multiple Workers. 
 Workers can be on remote systems, as long as they can reach the manager on the specified IP and port. The manager handles 
 reading the list of targets, writing the results, and which actual tests are run against each target. The workers just
 receive tasks from the manager, and run them.

 The target file is expected to be a simple text file, one target IP and port combination per line, comma-separated. So,
 the target file should look like this:
 
  192.168.1.1,80
  
  192.168.1.2,80
  
  176.16.1.1,443

 etc. 

 The whitelist file is a file of IPs that the system should skip, even if it's listed in the target file. The file is
 a text file with one IP per line.

 The "scans" option is expected to be a comma-separated list of the plugins to run against the targets. To list the 
 plugins available, run threebits with the "-mode l" option.

How to run it:
 If you run the threebits binary with the "-h" option, you'll get the following hints:
 
   -IP string
 
    	IP of Server (default "0.0.0.0")
 
  -mode string
 
    	mode to run in: m=Manager, w=worker, l=list plugins (default "l")
 
  -numworkers int
 
    	number of workers to run (Only used in worker mode) (default 5000)
 
  -output string
 
    	where to send output. if specified, will write to that filename (only used in manager mode) (default "output.csv")
 
  -port int
 
    	Port of Server (default 8000)
 
  -quiet
 
    	quiet mode: if false, will record falses as well as trues (only used in manager mode) (default true)
 
  -scans string
 
    	comma-separated list of commands to run (only used in manager mode)
 
  -targets string
 
    	file containing the targets to scan, one per line (only used in manager mode) (default "targetlist.txt")
 
  -timeout int
 
    	socket timeout in seconds (only used in worker mode) (default 3)
 
  -whitelist string
 
      	file containing a list of IPs to skip even if they are in the target list (only used in manager mode) (default 
      	"whitelist.txt")

So, for example, to run the manager with default settings, you would type:

 ./threebits -mode m -scans HTTPBanner
 
To run workers that connect to the manager, you'll need to know the IP of the manager (for example 192.168.1.100):

 ./threebits -mode w -IP 192.168.1.100
  

