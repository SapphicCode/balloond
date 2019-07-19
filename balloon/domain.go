package balloon

import (
	"github.com/digitalocean/go-libvirt"
)

const libvirtMemStatNr = uint32(libvirt.DomainMemoryStatNr)
const libvirtMemStatTagActualBalloon = int32(libvirt.DomainMemoryStatActualBalloon)
const libvirtMemStatTagUsable = int32(libvirt.DomainMemoryStatUsable)

// ProcessDomain will process a domain's memory balloon and log to Balloon.Logger.
func (balloon Balloon) ProcessDomain(dom libvirt.Domain) {
	// get stats
	stats, err := balloon.libvirt.DomainMemoryStats(dom, libvirtMemStatNr, 0)
	if err != nil {
		balloon.Logger.Err(err).Msg("Error fetching memory statistics.")
		return
	}

	// get max allowance
	maximum, err := balloon.libvirt.DomainGetMaxMemory(dom)
	if err != nil {
		balloon.Logger.Err(err).Msg("Error fetching maximum memory allowance.")
	}

	// compile statistics
	var current uint64   // current memory in use
	var available uint64 // "unused" enum in libvirt, memory available to VM in current balloon
	for _, stat := range stats {
		if stat.Tag == libvirtMemStatTagActualBalloon {
			current = stat.Val
		}
		if stat.Tag == libvirtMemStatTagUsable {
			available = stat.Val
		}
	}

	if available == 0 {
		// available == 0 can occur when the VM balloon driver isn't installed, so this is expected on some VMs.
		// Thus, our job is done here and we return immediately.
		balloon.Logger.Debug().Str("domain", dom.Name).Msg("Domain isn't returning useful memory statistics, ignoring.")
		return
	}
	if current == 0 {
		balloon.Logger.Error().Str("domain", dom.Name).Msg("Invalid state! Memory balloon is not returning current allocation.")
		// ...what happened here, libvirt?
		// Ideally this condition should never occur, but we're putting a check in place just to be sure.
		return
	}
	if current < available {
		balloon.Logger.Error().Str("domain", dom.Name).Msg("Libvirt is returning garbage statistics. Protecting ourselves.")
		return
	}

	// first check, is available memory less than the VM is supposed to have?
	if available < balloon.FreeAllowance {
		toAssign := current
		// loop in adding chunks until the FreeAllowance threshold is met again
		for available+toAssign-current < balloon.FreeAllowance {
			toAssign += balloon.MemoryChunk
		}

		if toAssign > maximum {
			toAssign = maximum // just following orders, sorry poor memory-starved VM :/
		}

		if current == toAssign {
			return // edge-case for VMs already at their highest memory pressure and presumably about to experience an OOM
		}

		if !balloon.DryRun {
			err := balloon.libvirt.DomainSetMemory(dom, toAssign)
			if err != nil {
				balloon.Logger.Err(err).Str("domain", dom.Name).Msg("Error setting new memory amount.")
				return
			}
		}
		balloon.Logger.Debug().Str("domain", dom.Name).Uint64("newAllocation", toAssign).Msg("Memory successfully allocated.")
		return
	}

	// second check, is available memory more than FreeAllowance plus a MemoryChunk
	if available > balloon.FreeAllowance+balloon.MemoryChunk {
		toAssign := current
		for available+toAssign-current > balloon.FreeAllowance+balloon.MemoryChunk {
			toAssign -= balloon.MemoryChunk
		}

		if !balloon.DryRun {
			err := balloon.libvirt.DomainSetMemory(dom, toAssign)
			if err != nil {
				balloon.Logger.Err(err).Str("domain", dom.Name).Msg("Error setting new memory amount.")
				return
			}
		}
		balloon.Logger.Debug().Str("domain", dom.Name).Uint64("newAllocation", toAssign).Msg("Memory successfully reaped.")
		return
	}
}
