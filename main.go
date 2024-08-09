package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
	"time"

)
type Command struct{
	name string 
	args []string
}

func main() {

	for{
  
	resp := make(chan string, 1024) 
	wg := &sync.WaitGroup{}
	wg.Add(5)
	go AllProcesses(resp, wg)
	go CPUDetails(resp, wg)
	go DiskDetails(resp, wg)
	go MemoryUsage(resp,wg)
	go NetworkInterfaces(resp,wg)	

	wg.Wait()
	close(resp)

	for r := range resp {
		fmt.Println(r)
	}
  time.Sleep(1*time.Second)
  }
  
}

func AllProcesses(resp chan string, wg *sync.WaitGroup) {
  
	cmd := Command{
		name:"ps",
		args:[]string{"aux"},
	}
	CMDOutput(cmd.name,cmd.args,resp)
	wg.Done()
}

func CPUDetails(resp chan string, wg *sync.WaitGroup) {
	cmd := Command{
		name:"lscpu",
		args:[]string{},
	}
	CMDOutput(cmd.name,cmd.args,resp)
	wg.Done()
}

func DiskDetails(resp chan string, wg *sync.WaitGroup) {
	cmd := Command{
		name:"free",
		args:[]string{"-m"},
	}
	CMDOutput(cmd.name,cmd.args,resp)
	wg.Done()
}

func MemoryUsage(resp chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := Command{
		name: "free",
		args: []string{"-m"},
	}
	CMDOutput(cmd.name, cmd.args, resp);
		
	
}

func NetworkInterfaces(resp chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := Command{
		name: "ip",
		args: []string{"-s", "link"},
	}
	CMDOutput(cmd.name, cmd.args, resp)
		
}

func CMDOutput(name string,args []string,resp chan string){
	cmd := exec.Command(name, args...)
	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating StdoutPipe for %s: %v\n", name, err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error while starting the command %s: %v\n", name, err)

	}

	// Read the command output
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
	  resp<-scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading scanner %s: %v\n", name, err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for command %s: %v\n", name, err)
	}
}
