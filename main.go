package main

import (

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Command struct{
	cmdStr string 
	key string
}

var (
	data = make(map[string]string)  //In future this may be changed to a time series database
	mutex   sync.Mutex
)

func main() {
	r := chi.NewRouter()
	// r.Get("/", )
	r.Get("/metrics", serveMetrics)

	go func() {
		fmt.Println("Running on port 3000...")
		if err := http.ListenAndServe(":3000", r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		for {
			fetchMetrics()
			time.Sleep(1 * time.Second) // Adjust the sleep duration as needed
										// In future it will be updated such that it can 
										//take input from users .
		}
	}()

	select {} // Keep the main function running
}

func fetchMetrics() {

	wg := &sync.WaitGroup{}
	wg.Add(7)
	// command,args,key,sync.wait
	
	go executePipeCommand("ps aux", "all_processes",wg)
	go executePipeCommand("lscpu", "cpu_details",wg)
	go executePipeCommand("df -h", "disk_details",wg)
	go executePipeCommand("free -m", "memory_usage",wg)
	go executePipeCommand("ip -s link", "network_interfaces",wg)
	go executePipeCommand("ip address", "ip_addresses",wg)
	go executePipeCommand("ls -l | grep main | wc -l", "pipe_command",wg)

	wg.Wait()

}

func executePipeCommand(cmd string, key string,wg *sync.WaitGroup) {
	defer wg.Done()
	c := strings.Split(cmd, "|")
	var cmdlist []*exec.Cmd
	for _, val := range c {
		f := strings.Fields(val)
		x := exec.Command(f[0], f[1:]...)
		cmdlist = append(cmdlist, x)
	}

	for i := 0; i < len(cmdlist)-1; i++ {
		j := i + 1
		out, err := cmdlist[i].StdoutPipe()
		if err != nil {
			fmt.Println(err)
			return
		}
		cmdlist[j].Stdin = out
	}

	out, _ := cmdlist[len(cmdlist)-1].StdoutPipe()

	for i := 0; i < len(cmdlist); i++ {
		if err := cmdlist[i].Start(); err != nil {
			fmt.Println("Error while starting the command", err)
			return
		}
	}

	res, err := io.ReadAll(out)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(cmdlist); i++ {
		if err := cmdlist[i].Wait(); err != nil {
			fmt.Println("Error while waiting", err)
		}
	}

	mutex.Lock()
	data[key] = string(res)
	mutex.Unlock()
	// log.Printf("Pipe command output: %s", string(res)) // Debugging line
}

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "index.html")
// }

func serveMetrics(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock() //
	// defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
