package worker

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// updateStats updates the worker's statistics by retrieving CPU and memory usage.
//
// Returns:
// - error: nil if the statistics were successfully updated, otherwise an error indicating the cause of the failure.
func (w *Worker) updateStats() error {

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil
	}
	vmm, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}
	w.WorkerStats.CPUUsage = cpuPercent[0]
	w.WorkerStats.MemAvail = vmm.Available
	w.WorkerStats.MemTotal = vmm.Total
	w.WorkerStats.MemUsed = vmm.Available * 100 / vmm.Total
	w.WorkerStats.Timestamp = time.Now()

	return nil
}
