package pkg

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strings"
// )

// func CpuInfo() {
// 	file,_ := os.Open("/proc/cpuinfo")
// 	x := bufio.NewScanner(file)
// 	for x.Scan(){
// 		fmt.Println(strings.Fields(x.Text()))
// 	}
// }