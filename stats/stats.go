package stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type diagnostic struct {
	CPU      []cpu.InfoStat         `json:"CPU"`
	CPUUsage float64                `json:"CPU Usage"`
	Memory   *mem.VirtualMemoryStat `json:"Memory"`
	Disk     []*disk.UsageStat      `json:"Disk"`
}

var cpuData []cpu.InfoStat
var cpuUsage float64
var memoryData *mem.VirtualMemoryStat
var diskData []*disk.UsageStat

// StartServer is a public function
func StartServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/stats", createDiagnostic)
	log.Fatal(http.ListenAndServe(":9000", router))
}

func runDiagnostic() {
	var diskUsage []*disk.UsageStat
	v, _ := mem.VirtualMemory()
	cores, _ := cpu.Percent(time.Second, false)
	cpuStats, _ := cpu.Info()
	parts, _ := disk.Partitions(false)

	for _, cpuStat := range cpuStats {
		printCPUStats(cpuStat)
		cpuData = append(cpuData, cpuStat)
	}

	for _, core := range cores {
		printCPUUsage(core)
		cpuUsage = core
	}

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		check(err)
		diskUsage = append(diskUsage, u)
		diskData = append(diskData, u)
		printDiskUsage(u)
	}

	printMemoryUsage(v)
	memoryData = v

}

func createDiagnostic(w http.ResponseWriter, r *http.Request) {
	runDiagnostic()
	foobar := diagnostic{
		CPU:      cpuData,
		CPUUsage: cpuUsage,
		Memory:   memoryData,
		Disk:     diskData,
	}
	json.NewEncoder(w).Encode(foobar)
}

func printCPUStats(cpu cpu.InfoStat) {
	fmt.Println("CPU:")
	fmt.Println("VendorID: " + cpu.VendorID)
	fmt.Println("Family: " + cpu.Family)
	fmt.Println("Model: " + cpu.Model)
	fmt.Println("Cores: " + strconv.FormatInt(int64(cpu.Cores), 10))
	fmt.Println("Mhz: " + strconv.FormatFloat(cpu.Mhz, 'f', 2, 64) + "Mhz")
}

func printCPUUsage(cpu float64) {
	fmt.Println("Total: " + strconv.FormatFloat(cpu, 'f', 2, 64) + "% usage.")
}

func printDiskUsage(u *disk.UsageStat) {
	fmt.Println("Disk:")
	fmt.Println(u.Path + "\t" + strconv.FormatFloat(u.UsedPercent, 'f', 2, 64) + "% full.")
	fmt.Println("Total: " + strconv.FormatUint(u.Total/1024/1024/1024, 10) + " GiB")
	fmt.Println("Free:  " + strconv.FormatUint(u.Free/1024/1024/1024, 10) + " GiB")
	fmt.Println("Used:  " + strconv.FormatUint(u.Used/1024/1024/1024, 10) + " GiB")
}

func printMemoryUsage(memory *mem.VirtualMemoryStat) {
	fmt.Println("Memory:")
	fmt.Println("Total: " + strconv.FormatUint(memory.Total/1024/1024/1024, 10) + " GiB")
	fmt.Println("Free: " + strconv.FormatUint(memory.Available/1024/1024/1024, 10) + " GiB")
	fmt.Println("Used: " + strconv.FormatUint(memory.Used/1024/1024/1024, 10) + " GiB")
	fmt.Println("Used %: " + strconv.FormatFloat(memory.UsedPercent, 'f', 2, 64) + "% used.")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
