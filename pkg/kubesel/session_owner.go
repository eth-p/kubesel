package kubesel

import (
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/process"
)

// SessionOwner identifies the shell that initialized a [Session].
// This is used for garbage collection purposes.
type SessionOwner struct {
	Process SessionOwnerPID `json:"pid"`
	Epoch   uint64          `json:"epoch"`
}

type SessionOwnerPID = int32

func (o *SessionOwner) fileName() string {
	pidHex := strconv.FormatInt(int64(o.Process), 16)
	bootTimeHex := strconv.FormatUint(o.Epoch, 16)
	return fmt.Sprintf("kubesel-%s-%s-kubeconfig.yaml", bootTimeHex, pidHex)
}

// SessionOwnerForProcess creates a [SessionOwner] using the specified process
// as the session's owner.
func SessionOwnerForProcess(pid SessionOwnerPID) (*SessionOwner, error) {
	// Get the system boot time.
	bootTime, err := host.BootTime()
	if err != nil {
		return nil, fmt.Errorf("finding boot time: %w", err)
	}

	// Check if the owner process is alive.
	exists, err := process.PidExists(pid)
	if err != nil {
		return nil, fmt.Errorf("checking process %d: %w", pid, err)
	}

	if !exists {
		return nil, fmt.Errorf("%w: %d", ErrOwnerProcessNotExist, pid)
	}

	return &SessionOwner{
		Process: pid,
		Epoch:   bootTime,
	}, nil
}
