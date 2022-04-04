package models

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Node struct {
	mgm.DefaultModel `bson:",inline"` //MongoDbModel
	Name             string           `json:"name"`
	Number           int              `json:"number"`
	CPUAvg           float64          `json:"cpuAvg"`
	CPU              float64          `json:"cpu"`
	Boot_UTCTime     uint64           `json:"boot_UTCTime"`
	Update_UTCTime   int64            `json:"update_UTCTime"`
	Memory           Resources        `json:"memory"`
	Disk             Resources        `json:"disk"`
}

type Resources struct {
	Total     uint64 `json:"total" validate:"required"`
	Available uint64 `json:"available" validate:"required"`
}

func (node *Node) setMyNode() error {
	// Read
	cpu,cpuAVG,err := getCpuAvg()
	if err != nil {
		return err
	}

	node.CPU = cpu
	node.CPUAvg = cpuAVG
	boottime, err := host.BootTime()
	if err != nil {
		return err
	}
	node.Boot_UTCTime = boottime
	node.Update_UTCTime = time.Now().UTC().Unix()
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	node.Memory.Total = v.Total
	node.Memory.Available = v.Available
	d, err := disk.Usage("/")
	if err != nil {
		return err
	}
	node.Disk.Total = d.Total
	node.Disk.Available = d.Free
	fmt.Println("Get Hardware info success!")
	return nil
}

func getCpuAvg() (float64,float64, error) {
	var cpuAvg float64
	var result string
	var lines []string
	if isFileExist("CPUAVG.info") {
		content, err := os.ReadFile("CPUAVG.info")
		if err != nil {
			return 0,0, err
		}
		Filelines = strings.Split(string(content), "\n")
		for _, value := range Filelines {
			if value == "" {
				continue
			}
			if s, err := strconv.ParseFloat(value, 64); err == nil {
				cpuAvg += s
				lines = append(lines, value)
			}
		}
	}
	array, err := cpu.Percent(0, false)
	if err != nil {
		return 0,0, err
	}
	cpuAvg += array[0]
	cpuAvg = cpuAvg / float64(len(lines)+1)
	lines = append(lines, fmt.Sprint(array[0]))

	if len(lines) > 60 {
		lines = lines[1:]
	}

	for _, value := range lines {
		temp := value + "\n"
		result += temp
	}
	os.WriteFile("CPUAVG.info", []byte(result), 0644)
	return array[0],cpuAvg, nil

}

func NodeReportStart(name string, number int) {

	err := mgm.SetDefaultConfig(nil, "closure", options.Client().ApplyURI(ConnectStr))
	if err != nil {
		os.Exit(-1)
	}
	var MyNode = &Node{}
	var NodeInfoColl = mgm.Coll(MyNode)
	MyNode.Name = name
	MyNode.Number = number
	err = NodeInfoColl.First(bson.M{"number": MyNode.Number}, MyNode)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			NodeInfoColl.Create(MyNode)
			fmt.Println("Create a new node in mongodb")
		}
	}
	MyNode.setMyNode()
	err = NodeInfoColl.Update(MyNode)
	if err != nil {
		return
	}
	fmt.Println("Update myNode to mongodb success!")
}
