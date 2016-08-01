package Threebits

import (
	"bufio"
	"os"
	"io"
	"strings"
	"fmt"
	"sync"
	"errors"
	"net/rpc"
	"net"
	"strconv"
	"net/http"
	"github.com/accessviolationsec/Threebits/structures"
)


type TestProvider struct {
	Channel chan structures.Test
	AuthKey string
	Done chan struct{}
}

type ResponseProvider struct {
	Channel chan structures.Response
	AuthKey string
	Done chan struct{}
}


func (t *TestProvider) GetJob(authkey string, reply *structures.Test) error{
	var input structures.Test
	if authkey != t.AuthKey{
		return errors.New("Authkey failed")
	}
	select {
		case input = <- t.Channel:
			*reply = input
			return nil
		case <- t.Done: return nil
	}
	return nil
}

func (r *ResponseProvider) SendResponse(response structures.Response, success *bool) error {
	if response.AuthKey != r.AuthKey {
		return errors.New("AuthKey Failed")
	}
	select {
		case r.Channel <- response: return nil
		case <- r.Done: return nil
	}
	return nil
}


func readWhitelist(whitelist string) (map[string]bool){
	var fileHandle io.Reader
	var err error
	var response = make(map[string]bool)

	fileHandle, err = os.Open(whitelist)
	if err != nil{
		return response
	}
	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan(){
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#"){
			continue
		}
		if len(line) == 0 {
			continue
		}
		response[line] = true
	}
	return response
}

func readTargets(targets string, whitelist string, scans string, input chan structures.Test, numJobs *sync.WaitGroup, Done chan struct{}){
	// open file, write things onto input channel
	var fileHandle io.Reader
	var err error

	skipIPs := readWhitelist(whitelist)


	scanList := strings.Split(scans, ",")
	if targets == "stdio" || targets == "-"{
		fileHandle = os.Stdin
	} else {
		fileHandle, err = os.Open(targets)
		if err != nil{
			fmt.Println("Error opening file.", err)
		}
	}
	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		// this should be IP,port in the file
		target := strings.Split(scanner.Text(), ",")
		// skip comments
		if strings.HasPrefix(target[0], "#"){
			continue
		}
		if err := scanner.Err(); err != nil{
			fmt.Println("Error reading input", err)
		}
		if len(target) < 2 {
			fmt.Println("Error splitting target:", target)
			continue
		}
		if skipIPs[target[0]]{
			continue
		}
		for _, scan := range scanList {
			port, err := strconv.Atoi(target[1])
			if err != nil{
				continue
			}
			message := structures.Test{
				Target: target[0],
				Port: port,
				Test: scan,
			}
			select {
				case input <- message:
					numJobs.Add(1)
				case <- Done:
					return
			}
		}
	}
	fmt.Println("manager, done reading input file. ending")
}

func readResults(outputpath string,  quiet bool, results chan structures.Response, wg *sync.WaitGroup, numJobs *sync.WaitGroup, Done chan struct{}){
	defer wg.Done()

	fileHandle, err := os.Create(outputpath)
	var resp structures.Response
	if err != nil {
		fmt.Println("can't open output", err)
	}
	defer fileHandle.Close()

	MAINLOOP: for {
		select {
			case resp = <- results:
				numJobs.Done()

				if quiet == true && resp.Success == false {
					continue
				}
				_, err := fileHandle.WriteString(resp.Stringify() + "\n")
				if err != nil {
					fmt.Println("error writing to output file", err)
				}
			case <-Done:
				break MAINLOOP
		}
	}
}


func RunManager(IP string, port int, targetfile string, whitelist string, scans string, output string, quiet bool, parentWg *sync.WaitGroup) {
	defer parentWg.Done()

	var wg sync.WaitGroup
	var numJobs sync.WaitGroup


	Done := make(chan struct{})

	input := make(chan structures.Test)
	results := make(chan structures.Response)

	wg.Add(1)
	go readResults(output, quiet, results, &wg, &numJobs, Done)


	t := TestProvider{Channel: input, AuthKey: "test", Done: Done}
	r := ResponseProvider{Channel: results, AuthKey: "test", Done: Done}
	rpc.Register(&t)
	rpc.Register(&r)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", IP + ":" + strconv.Itoa(port))
	if err != nil{
		fmt.Println("Error setting up listener:", err)
	}
	go http.Serve(listener, nil)

	readTargets(targetfile, whitelist, scans, input, &numJobs, Done)

	// wait until all the jobs are returned, then signal the reader to close
	numJobs.Wait()
	close(Done)
	wg.Wait()
}
