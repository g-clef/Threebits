package Threebits


import (
	"net/rpc"
	"strconv"
	"fmt"
	"net"
	"sync"
	"github.com/accessviolationsec/Threebits/structures"
	"time"
)

func jobGetter(serverIP string, port int, authKey string, input chan structures.Test, wg *sync.WaitGroup, Done chan struct{}){
	defer wg.Done()

	client, err := rpc.DialHTTP("tcp", serverIP + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error connecting to server", err)
	}

	MAINLOOP: for {
		var test structures.Test

		err := client.Call("TestProvider.GetJob", authKey, &test)
		if err != nil{
			break
		}
		// we might get an empty test...this seems to also indicate that
		// the channel on the other end is closed, and we should stop.
		if test == (structures.Test{}){
			break
		}
		select {
		case input <- test:
		case <-Done: break MAINLOOP
		}
	}
	client.Close()
}

func responseSender(serverIP string, port int, authKey string, output chan structures.Response, Done chan struct{}){

	var resp structures.Response


	client, err := rpc.DialHTTP("tcp", serverIP + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error connecting to server", err)
	}

	MAINLOOP: for {
		var success bool

		select {
		case resp = <- output:
			resp.AuthKey = authKey
			err := client.Call("ResponseProvider.SendResponse", resp, &success)
			if err != nil{
				fmt.Println("Error sending response", err)
				break MAINLOOP
			}
		case <- Done: break MAINLOOP
		}
	}
	client.Close()
}


func RunWorkers(serverIP string, port int, numWorkers int, timeout int, wg *sync.WaitGroup){
	// In normal operation, the job getter will see that the connection
	// from the manager closed, which will cause the job getter to close.
	// At that point, the workerWG.wait() will return, Done() will be closed,
	// which should signal the workers and responseSender to shut down also.
	defer wg.Done()

	var Done = make(chan struct{})
	var workerWG sync.WaitGroup

	input := make(chan structures.Test)
	results := make(chan structures.Response)


	workerWG.Add(1)
	go jobGetter(serverIP, port, "test", input, &workerWG, Done)
	go responseSender(serverIP, port, "test", results, Done)
	for i:=0; i < numWorkers; i++{
		go worker(timeout, input, results, Done)
	}

	workerWG.Wait()
	close(Done)
}


func worker(timeout int, input chan structures.Test, output chan structures.Response, Done chan struct{}){
	var test structures.Test
	var response structures.Response
	var responseMessage string
	var tOut time.Duration = time.Duration(timeout) * time.Second

	for {
		select {
		case test = <-input:
			for pluginname, plugin := range Plugins {
				if test.Test != pluginname {
					continue
				}
				sock, err := net.DialTimeout(plugin.Protocol(), test.GetAddr(),tOut)
				if err != nil {
					response = structures.Response{Success:false, Test: test, Message:err.Error()}
				} else {
					now := time.Now()
					err = sock.SetDeadline(now.Add(tOut))
					success, message, err := plugin.Handle(sock, test)
					if err != nil {
						responseMessage = message + err.Error()
					} else {
						responseMessage = message
					}
					response = structures.Response{Success:success, Test:test, Message:responseMessage}
					sock.Close()
				}
				output <- response
			}
		case <-Done:
			return
		}
	}
}