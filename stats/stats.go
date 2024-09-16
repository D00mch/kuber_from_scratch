package stats

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

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
