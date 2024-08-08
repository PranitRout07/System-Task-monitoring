package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"
	// "time"
)

func main() {
	// now := time.Now()
	for{
	// cmd1 := exec.Command("ps", "aux")
	// cmd2 := exec.Command("free", "-m")
	// cmd3 := exec.Command("lscpu")
  name1 := "ps"
  args1 := []string{"aux"}
  name3 := "free"
  args3 := []string{"-m"}
  name2 := "lscpu"
  args2 := []string{}
  
	resp := make(chan string, 1024) //three messages are comming inside the channel . so its minimum buffer size is 3.
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go allProcesses(name1,args1, resp, wg)
	go cpuDetails(name2,args2, resp, wg)
	go diskSize(name3,args3, resp, wg)
	// time.Sleep(130*time.Millisecond)

	wg.Wait()
	close(resp)

	for r := range resp {
		fmt.Println(r)
	}
  time.Sleep(1*time.Second)
  }
  
	// fmt.Println(time.Since(now))
}

func allProcesses(name string, flags []string, resp chan string, wg *sync.WaitGroup) {
  
  
  cmd := exec.Command(name, flags...)
	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		// fmt.Printf("Error creating StdoutPipe for %s: %v\n", name, err)
		log.Fatal(err)

	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)

	}

	// Read the command output
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
	  resp<-scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}


	// time.Sleep(80 * time.Millisecond)
	wg.Done()
}

func cpuDetails(name string,flags []string, resp chan string, wg *sync.WaitGroup) {
  cmd := exec.Command(name, flags...)
	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		// fmt.Printf("Error creating StdoutPipe for %s: %v\n", name, err)
		log.Fatal(err)

	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)

	}

	// Read the command output
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
	  resp<-scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// time.Sleep(120 * time.Millisecond)
	wg.Done()
}

func diskSize(name string,flags []string , resp chan string, wg *sync.WaitGroup) {
  cmd := exec.Command(name, flags...)
	//
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		// fmt.Printf("Error creating StdoutPipe for %s: %v\n", name, err)
		log.Fatal(err)

	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)

	}

	// Read the command output
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
	  resp<-scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}
