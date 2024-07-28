package worker

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

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
