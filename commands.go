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
    	CpuCount               int64   `json:"cpuCount"`
    	CpuLoadMin1            float64 `json:"cpuLoadMin1"`
    	CpuLoadMin5            float64 `json:"cpuLoadMin5"`
    	CpuLoadMin15           float64 `json:"cpuLoadMin15"`
    	MemoryTotal            int64   `json:"memoryTotal"`
    	MemoryUsed             int64   `json:"memoryUsed"`
    	MemoryFree             int64   `json:"memoryFree"`
    	MemoryBuffersCacheUsed int64   `json:"memoryBuffersCacheUsed"`
    	MemoryBuffersCacheFree int64   `json:"memoryBuffersCacheFree"`
    	MemorySwapTotal        int64   `json:"memorySwapTotal"`
    	MemorySwapUsed         int64   `json:"memorySwapUsed"`
    	MemorySwapFree         int64   `json:"memorySwapFree"`

    }

	time := time.Now().Unix()

	cmd1, err := exec.Command("/bin/cat", "/proc/loadavg").Output()

	if err != nil {
		log.Fatal(err)
	}

	cpu := strings.Split(string(cmd1), " ")

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

	cmd2, err := sh.Command("/usr/bin/free").Command("/usr/bin/awk", "NR == 2 {printf $2 \" \" $3 \" \" $4 \" \"} NR == 3 {printf $3 \" \" $4 \" \"} NR ==4 {printf $2 \" \" $3 \" \" $4}").Output()

	if err != nil {
		log.Fatal(err)
	}

	memory := strings.Split(string(cmd2), " ")

	MemoryTotal, err := strconv.ParseInt(memory[0], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}

	MemoryUsed, err := strconv.ParseInt(memory[1], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}

	MemoryFree, err := strconv.ParseInt(memory[2], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}

	MemoryBuffersCacheUsed, err := strconv.ParseInt(memory[3], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}

	MemoryBuffersCacheFree, err := strconv.ParseInt(memory[4], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}

	MemorySwapTotal, err := strconv.ParseInt(memory[5], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}

	MemorySwapUsed, err := strconv.ParseInt(memory[6], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}

	MemorySwapFree, err := strconv.ParseInt(memory[7], 10, 64)
	
	if err != nil {
		log.Fatal(err)
	}
    
    snapshot := Snapshot{
    	strings.TrimSpace(
    		string(key)), 
    		time, 8,  //TODO
    		cpuLoad1Min, 
    		cpuLoad5Min, 
    		cpuLoad15Min, 
    		MemoryTotal, 
    		MemoryUsed, 
    		MemoryFree, 
    		MemoryBuffersCacheUsed, 
    		MemoryBuffersCacheFree, 
    		MemorySwapTotal, 
    		MemorySwapUsed, 
    		MemorySwapFree,
    	}

    output, err := json.Marshal(snapshot)

    fmt.Print(string(output))
}

func doConfigure(c *cli.Context) {

	fmt.Println(chalk.Bold.TextStyle("What is the server UID ?"))
}
