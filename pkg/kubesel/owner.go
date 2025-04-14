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
	ownerData
}

type ownerData struct {
	Process PidType `json:"pid"`
	Epoch   uint64  `json:"epoch"`
}

func (o *Owner) fileName() string {
	pidHex := strconv.FormatInt(int64(o.Process), 16)
	bootTimeHex := strconv.FormatUint(o.Epoch, 16)
	return fmt.Sprintf("kubesel-%s-%s-kubeconfig.yaml", bootTimeHex, pidHex)
}

// IsAlive returns true if the owner is still alive.
//
// An owner is considered alive if the [Owner]'s process is not dead,
// and if the system hasn't rebooted since the [Owner] was first created.
func (o *Owner) IsAlive() (bool, error) {
	// Get the system boot time.
	bootTime, err := host.BootTime()
	if err != nil {
		return false, fmt.Errorf("finding epoch time: %w", err)
	}

	if o.Epoch != bootTime {
		return false, nil
	}

	// Check if the owner process is alive.
	exists, err := process.PidExists(o.Process)
	if err != nil {
		return false, fmt.Errorf("checking process %d: %w", o.Process, err)
	}

	return exists, nil
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
		ownerData: ownerData{
			Process: pid,
			Epoch:   bootTime,
		},
	}, nil
}
