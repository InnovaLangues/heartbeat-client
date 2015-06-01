package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
    "io/ioutil"
    "encoding/json"
    "time"
    "strings"
    "strconv"

	"github.com/codegangsta/cli"
	"github.com/ttacon/chalk"
	"github.com/codeskyblue/go-sh"

)

var Commands = []cli.Command{
	commandPush,
	commandConfigure,
}

var commandPush = cli.Command{
	Name:  "push",
	Usage: "Pushes json data to the hearbeat server",
	Description: `This command pushes a json string to the hearbeat server
`,
	Action: doPush,
}

var commandConfigure = cli.Command{
	Name:  "configure",
	Usage: "Configure the client UID",
	Description: `This command configures the client (this) machines UID
`,
	Action: doConfigure,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doPush(c *cli.Context) {
	key, err := ioutil.ReadFile("/etc/heartbeat-client/heartbeat-client.key")
    
    if err != nil {
		log.Fatal(key)
	}

    type Snapshot struct {
    	Uid                    string  `json:"uid"`
    	Timestamp              int64   `json:"timestamp"`
    	CpuCount               int     `json:"cpuCount"`
    	CpuLoadMin1            float64 `json:"cpuLoadMin1"`
    	CpuLoadMin5            float64 `json:"cpuLoadMin5"`
    	CpuLoadMin15           float64 `json:"cpuLoadMin15"`
    	MemoryTotal            int     `json:"memoryTotal"`
    	MemoryUsed             int     `json:"memoryUsed"`
    	MemoryFree             int     `json:"memoryFree"`
    	MemoryBuffersCacheUsed int     `json:"memoryBuffersCacheUsed"`
    	MemoryBuffersCacheFree int     `json:"memoryBuffersCacheFree"`
    	MemorySwapTotal        int     `json:"memorySwapTotal"`
    	MemorySwapUsed         int     `json:"memorySwapUsed"`
    	MemorySwapFree         int     `json:"memorySwapFree"`
    	DiskTotal              int     `json:"diskTotal"`
    	DiskUsed               int     `json:"diskUsed"`
    	DiskFree               int     `json:"diskFree"`
    }
	
	timestamp := time.Now().Unix()

	cmd1, err := exec.Command("/bin/grep", "-c", "^processor", "/proc/cpuinfo").Output()
	
	if err != nil {
		log.Fatal(err)
	}

	cpuCount, err := strconv.Atoi(strings.TrimSpace(string(cmd1)))

	if err != nil {
		log.Fatal(err)
	}

	cmd2, err := exec.Command("/bin/cat", "/proc/loadavg").Output()

	if err != nil {
		log.Fatal(err)
	}

	cpu := strings.Split(string(cmd2), " ")

	cpuLoad1Min, err := strconv.ParseFloat(cpu[0], 64)
	
	if err != nil {
		log.Fatal(err)
	}

	cpuLoad5Min, err := strconv.ParseFloat(cpu[1], 64)
	
	if err != nil {
		log.Fatal(err)
	}

	cpuLoad15Min, err := strconv.ParseFloat(cpu[2], 64)
	
	if err != nil {
		log.Fatal(err)
	}

	cmd3, err := sh.Command("/usr/bin/free").Command("/usr/bin/awk", "NR == 2 {printf $2 \" \" $3 \" \" $4 \" \"} NR == 3 {printf $3 \" \" $4 \" \"} NR ==4 {printf $2 \" \" $3 \" \" $4}").Output()

	if err != nil {
		log.Fatal(err)
	}

	memory := strings.Split(string(cmd3), " ")

	memoryTotal, err := strconv.Atoi(memory[0])
	
	if err != nil {
		log.Fatal(err)
	}

	memoryUsed, err := strconv.Atoi(memory[1])
	
	if err != nil {
		log.Fatal(err)
	}

	memoryFree, err := strconv.Atoi(memory[2])
	
	if err != nil {
		log.Fatal(err)
	}

	memoryBuffersCacheUsed, err := strconv.Atoi(memory[3])
	
	if err != nil {
		log.Fatal(err)
	}

	memoryBuffersCacheFree, err := strconv.Atoi(memory[4])
	
	if err != nil {
		log.Fatal(err)
	}

	memorySwapTotal, err := strconv.Atoi(memory[5])
	
	if err != nil {
		log.Fatal(err)
	}

	memorySwapUsed, err := strconv.Atoi(memory[6])
	
	if err != nil {
		log.Fatal(err)
	}

	memorySwapFree, err := strconv.Atoi(memory[7])
	
	if err != nil {
		log.Fatal(err)
	}

	cmd4, err := sh.Command("/bin/df", "--total").Command("/usr/bin/awk", "/total/ {printf $2 \" \" $3 \" \" $4}").Output()

	if err != nil {
		log.Fatal(err)
	}

	disk := strings.Split(string(cmd4), " ")

	diskTotal, err := strconv.Atoi(disk[0])
	
	if err != nil {
		log.Fatal(err)
	}

	diskUsed, err := strconv.Atoi(disk[1])
	
	if err != nil {
		log.Fatal(err)
	}

	diskFree, err := strconv.Atoi(disk[2])
	
	if err != nil {
		log.Fatal(err)
	}
    
    snapshot := Snapshot{
    	strings.TrimSpace(
    		string(key)), 
    		timestamp,
    		cpuCount,
    		cpuLoad1Min, 
    		cpuLoad5Min, 
    		cpuLoad15Min, 
    		memoryTotal, 
    		memoryUsed, 
    		memoryFree, 
    		memoryBuffersCacheUsed, 
    		memoryBuffersCacheFree, 
    		memorySwapTotal, 
    		memorySwapUsed, 
    		memorySwapFree,
    		diskTotal, 
    		diskUsed, 
    		diskFree, 
    	}

    output, err := json.Marshal(snapshot)

    fmt.Print(string(output))
}

func doConfigure(c *cli.Context) {

	fmt.Println(chalk.Bold.TextStyle("What is the server UID ?"))
}
