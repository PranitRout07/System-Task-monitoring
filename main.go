package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Command struct {
	name string
	args []string
}

func main() {

	for {

		resp := make(chan string, 1024)
		wg := &sync.WaitGroup{}
		wg.Add(9)
		go AllProcesses(resp, wg)
		go CPUDetails(resp, wg)
		go DiskDetails(resp, wg)
		go MemoryUsage(resp, wg)
		go NetworkInterfaces(resp, wg)
		go UpTime(resp, wg)
		go Top(resp, wg)
		go NetworkStats(resp, wg)
		go PipeCommand(resp,wg)
		wg.Wait()
		close(resp)

		for r := range resp {
			fmt.Println(r)
		}
		time.Sleep(1 * time.Second)
	}

}

func AllProcesses(resp chan string, wg *sync.WaitGroup) {

	cmd := Command{
		name: "ps",
		args: []string{"aux"},
	}
	CMDOutput(cmd.name, cmd.args, resp)
	wg.Done()
}

func CPUDetails(resp chan string, wg *sync.WaitGroup) {
	cmd := Command{
		name: "lscpu",
		args: []string{},
	}
	CMDOutput(cmd.name, cmd.args, resp)
	wg.Done()
}

func DiskDetails(resp chan string, wg *sync.WaitGroup) {
	cmd := Command{
		name: "free",
		args: []string{"-m"},
	}
	CMDOutput(cmd.name, cmd.args, resp)
	wg.Done()
}

func MemoryUsage(resp chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := Command{
		name: "free",
		args: []string{"-m"},
	}
	CMDOutput(cmd.name, cmd.args, resp)

}

func NetworkInterfaces(resp chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := Command{
		name: "ip",
		args: []string{"-s", "link"},
	}
	CMDOutput(cmd.name, cmd.args, resp)

}

func UpTime(resp chan string, wg *sync.WaitGroup){
	cmd := Command{
		name: "uptime",
		args: []string{},
	}
	CMDOutput(cmd.name, cmd.args, resp)
	
	defer wg.Done()
}

func Top(resp chan string, wg *sync.WaitGroup){
	cmd := Command{
		name: "Top",
		args: []string{"-b", "-n", "1"},
	}
	CMDOutput(cmd.name, cmd.args, resp)
	defer wg.Done()
}

func NetworkStats(resp chan string, wg *sync.WaitGroup){
	cmd := Command{
		name: "Netstat",
		args: []string{"-i"},
	}
	CMDOutput(cmd.name, cmd.args, resp)
	defer wg.Done()
}


func CMDOutput(name string, args []string, resp chan string) {
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
		resp <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading scanner %s: %v\n", name, err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for command %s: %v\n", name, err)
	}
}

// COMMANDS USING MULTIPLE PIPES EXAMPLE
func PipeCommand(resp chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := "ls -l | grep main | wc -l"
	Execute(cmd,resp)

}

func Execute(cmd string,resp chan string) {
	c := strings.Split(cmd, "|")
	var cmdlist []*exec.Cmd
	for _, val := range c {
		f := strings.Fields(val)
		x := exec.Command(f[0], f[1:]...)
		fmt.Println(x)
		cmdlist = append(cmdlist, x)
	}
// pipe connections
	for i := 0; i < len(cmdlist)-1; i++ {
		j := i + 1

		out, err := cmdlist[i].StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}
		cmdlist[j].Stdin = out
	}

	out, _ := cmdlist[len(cmdlist)-1].StdoutPipe()

	for i := 0 ;i<len(cmdlist);i++{
		err := cmdlist[i].Start()
		if err!=nil{
			fmt.Println("Error while starting the command",err)
		}
	}
	res,err := io.ReadAll(out)
	if err!=nil{
		log.Println(err)
	}
	for i := 0 ;i<len(cmdlist);i++{
		err :=cmdlist[i].Wait()
		if err!=nil{
			fmt.Println("Error while waiting",err)
		}
	}


	resp<- string(res)

}
