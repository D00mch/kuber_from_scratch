package stats

import (
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func TestMemoryKb(t *testing.T) {

	memStats, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("Can't get memory, %v", err)
	}

	stats := Stats{MemStats: memStats}

	if stats.MemTotalKb() > 65_000_000 {
		t.Fatalf("Impossible to have that much RAM: %v", stats.MemTotalKb())
	}
	if stats.MemAvailableKb() > 65_000_000 {
		t.Fatalf("Impossible to have that much available RAM: %v",
			stats.MemTotalKb())
	}
	if stats.MemUsedKb() > 65_000_000 {
		t.Fatalf("Impossible to have that much available RAM: %v",
			stats.MemTotalKb())
	}

	fmt.Printf("Stats: %v", stats)

	fmt.Printf("\nMemory used percent: %v\n", stats.MemUsedPercent())
}

func TestDisk(t *testing.T) {
	diskStats, err := disk.Usage("/")
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}
	stats := Stats{DiskStats: diskStats}
	fmt.Printf("Stats: %v\n", stats)
}

func TestLoad(t *testing.T) {
	loadStats, err := GetLoadStats()
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}
	fmt.Printf("LoadStats: %v", loadStats)
}
