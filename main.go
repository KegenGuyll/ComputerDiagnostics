package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {

	var diskUsage []*disk.UsageStat

	v, _ := mem.VirtualMemory()
	cores, _ := cpu.Percent(time.Second, false)
	cpuStats, _ := cpu.Info()
	parts, _ := disk.Partitions(false)

	for _, cpuStat := range cpuStats {
		printCPUStats(cpuStat)
	}

	for _, core := range cores {
		printCPUUsage(core)
	}

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		check(err)
		diskUsage = append(diskUsage, u)
		printDiskUsage(u)
	}

	printMemoryUsage(v)

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
	fmt.Println("Total: " + strconv.FormatFloat(cpu, 'f', 2, 64) + "% usage.\n")
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
