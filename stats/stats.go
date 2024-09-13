package stats

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// # github.com/shirou/gopsutil/disk
// iostat_darwin.c:28:2: warning: 'IOMasterPort' is deprecated: first deprecated in macOS 12.0
//  [-Wdeprecated-declarations]
// /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/System/Library/Frameworks/IOKit.framewo
// rk/Headers/IOKitLib.h:143:1: note: 'IOMasterPort' has been explicitly marked deprecated her
// e
// # github.com/shirou/gopsutil/host
// smc_darwin.c:75:41: warning: 'kIOMasterPortDefault' is deprecated: first deprecated in macO
// S 12.0 [-Wdeprecated-declarations]
// /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/System/Library/Frameworks/IOKit.framewo
// rk/Headers/IOKitLib.h:133:19: note: 'kIOMasterPortDefault' has been explicitly marked depre
// cated here
// hostStat: {"hostname":"air","uptime":353164,"bootTime":1724874571,"procs":649,"os":"darwin"
// ,"platform":"darwin","platformFamily":"Standalone Workstation","platformVersion":"14.4","ke
// rnelVersion":"23.4.0","kernelArch":"arm64","virtualizationSystem":"","virtualizationRole":"
// ","hostid":"a776b25d-7381-58cb-98b7-7c90f7a9c2f8"}
// vmState: {"total":17179869184,"available":6084263936,"used":11095605248,"usedPercent":64.58
// 492279052734,"free":882704384,"active":5428756480,"inactive":5201559552,"wired":2563555328,
// "laundry":0,"buffers":0,"cached":0,"writeback":0,"dirty":0,"writebacktmp":0,"shared":0,"sla
// b":0,"sreclaimable":0,"sunreclaim":0,"pagetables":0,"swapcached":0,"commitlimit":0,"committ
// edas":0,"hightotal":0,"highfree":0,"lowtotal":0,"lowfree":0,"swaptotal":0,"swapfree":0,"map
// ped":0,"vmalloctotal":0,"vmallocused":0,"vmallocchunk":0,"hugepagestotal":0,"hugepagesfree"
// :0,"hugepagesize":0}
// diskStat: {"path":"/","fstype":"apfs","total":245107195904,"free":76771721216,"used":168335
// 474688,"usedPercent":68.67830789999783,"inodesTotal":750127592,"inodesUsed":403752,"inodesF
// ree":749723840,"inodesUsedPercent":0.05382444324218379}

type CpuStats struct {
	Usage float64
}

type LoadStats struct {
	Avg []float64
}

type Stats struct {
	MemStats  *mem.VirtualMemoryStat
	DiskStats *disk.UsageStat
	CpuStats  *CpuStats
	LoadStats *LoadStats
	TaskCount int
}

func GetStats() *Stats {
	return &Stats{
		MemStats: GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		LoadStats: GetLoadAvg(),
	}
}

func (s *Stats) MemAvailableKb() uint64 { return s.MemStats.Available / 1024 }
func (s *Stats) MemTotalKb() uint64     { return s.MemStats.Total / 1024 }
func (s *Stats) MemUsedKb() uint64      { return s.MemStats.Used / 1024 }
func (s *Stats) MemUsedPercent() float64 {
	return float64(s.MemStats.Available) / float64(s.MemStats.Total)
}

func (s *Stats) DiskTotal() uint64 { return s.DiskStats.Total }
func (s *Stats) DiskFree() uint64  { return s.DiskStats.Free }
func (s *Stats) DiskUsed() uint64  { return s.DiskStats.Used }

func GetMemoryInfo() *mem.VirtualMemoryStat {
	memstats, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error reading from /proc/meminfo")
		return &mem.VirtualMemoryStat{}
	}
	return memstats
}

func GetDiskInfo() *disk.UsageStat {
	diskstats, err := disk.Usage("/")
	if err != nil {
		log.Printf("Error reading from /")
		return &disk.UsageStat{}
	}
	return diskstats
}

// GetCpuInfo See https://godoc.org/github.com/c9s/goprocinfo/linux#CPUStat
// func GetCpuStats() *cpu {
// 	stats, err := linux.ReadStat("/proc/stat")
// 	if err != nil {
// 		log.Printf("Error reading from /proc/stat")
// 		return &linux.CPUStat{}
// 	}
// 	return &stats.CPUStatAll
// }

// GetLoadAvg See https://godoc.org/github.com/c9s/goprocinfo/linux#LoadAvg
func GetLoadAvg() *LoadStats {
	loadavg, err := GetLoadStats()
	if err != nil {
		log.Printf("Error getting load")
		return &LoadStats{}
	}
	return loadavg
}

func GetLoadStats() (*LoadStats, error) {
	cmd := exec.Command("sysctl", "-n", "vm.loadavg")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Output is like "{0.87 1.01 1.05}", so we trim and split it
	loadStr := strings.Trim(string(output), "{} \n")
	loadParts := strings.Fields(loadStr)

	var loadavg []float64
	for _, load := range loadParts {
		value, err := strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, err
		}
		loadavg = append(loadavg, value)
	}

	return &LoadStats{Avg: loadavg}, nil
}
