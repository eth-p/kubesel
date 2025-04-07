package kubesel

import (
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/process"
)

type PidType = int32

// Owner identifies the shell that initialized a [ManagedKubeconfig].
// This is used for garbage collection purposes.
type Owner struct {
	Process PidType `json:"pid"`
	Epoch   uint64  `json:"epoch"`
}

func (o *Owner) fileName() string {
	pidHex := strconv.FormatInt(int64(o.Process), 16)
	bootTimeHex := strconv.FormatUint(o.Epoch, 16)
	return fmt.Sprintf("kubesel-%s-%s-kubeconfig.yaml", bootTimeHex, pidHex)
}

// OwnerForProcess creates an [Owner] using the specified process
// as the session's owner.
func OwnerForProcess(pid PidType) (*Owner, error) {
	// Get the system boot time.
	bootTime, err := host.BootTime()
	if err != nil {
		return nil, fmt.Errorf("finding epoch time: %w", err)
	}

	// Check if the owner process is alive.
	exists, err := process.PidExists(pid)
	if err != nil {
		return nil, fmt.Errorf("checking process %d: %w", pid, err)
	}

	if !exists {
		return nil, fmt.Errorf("%w: %d", ErrOwnerProcessNotExist, pid)
	}

	return &Owner{
		Process: pid,
		Epoch:   bootTime,
	}, nil
}
