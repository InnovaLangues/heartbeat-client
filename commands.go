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
    	Uid          string  `json:"uid"`
    	Timestamp    int64   `json:"timestamp"`
    	CpuCount     int     `json:"cpuCount"`
    	CpuLoadMin1  float64 `json:"cpuLoadMin1"`
    	CpuLoadMin5  float64 `json:"cpuLoadMin5"`
    	CpuLoadMin15 float64 `json:"cpuLoadMin15"`
    	Memory       int     `json:"cpuLoadMin15"`
    }

	time := time.Now().Unix()

	cmd := exec.Command("/bin/cat", "/proc/loadavg")

	// /usr/bin/free | /usr/bin/awk 'NR == 2 {printf $2 " " $3 " " $4 " "} NR == 3 {printf $3 " " $4 " "} NR ==4 {printf $2 " " $3 "  " $4}'

	output, err := cmd.Output()

	cpu := strings.Split(string(output), " ")

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
    
    s := Snapshot{strings.TrimSpace(string(key)), time, 8,  cpuLoad1Min, cpuLoad5Min, cpuLoad15Min}

    b, err := json.Marshal(s)

    fmt.Print(string(b))
}

func doConfigure(c *cli.Context) {

	fmt.Println(chalk.Bold.TextStyle("What is the server UID ?"))
}
