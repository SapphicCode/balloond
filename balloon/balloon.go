// Package balloon provides an automatic libvirt memory balloon daemon.
package balloon

import (
	"net"
	"os"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

// Balloon describes a hypervisor memory balloon interface.
type Balloon struct {
	libvirt *libvirt.Libvirt
	cache   *cache.Cache

	Logger zerolog.Logger // Logger is the zerolog Logger. New creates a default, but feel free to override it.

	Interval time.Duration // Interval is the interval at which the daemon should refresh all domain balloons.
	DryRun   bool          // DryRun effectively pretends running actions like DomainSetMemory, but doesn't.

	FreeAllowance uint64 // FreeAllowance is the amount of memory, in kB, to allow the VM to have available.
	MemoryChunk   uint64 // MemoryChunk is the chunk size of memory (in kB) to allocate or deallocate. This is a granularity slider.
}

// New creates a Balloon from a net.Conn and populates its settings with default values.
func New(conn net.Conn) Balloon {
	l := libvirt.New(conn)
	c := cache.New(1*time.Minute, 1*time.Second)
	return Balloon{
		l,
		c,

		zerolog.New(os.Stdout),

		10 * time.Second,
		false,

		256 * 1024,
		32 * 1024,
	}
}
