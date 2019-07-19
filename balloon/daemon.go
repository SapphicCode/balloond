package balloon

import (
	"time"

	"github.com/digitalocean/go-libvirt"
)

// RunDaemon simply runs ProcessDomain for every domain in the libvirt connection every Interval.
func (balloon Balloon) RunDaemon() {
	balloon.Logger.Info().Msg("Balloon daemon started. Establishing connection.")
	err := balloon.libvirt.Connect()
	if err != nil {
		balloon.Logger.Fatal().Err(err).Msg("Error connecting to libvirtd.")
	}
	balloon.Logger.Info().Msg("Connection successfully established.")

	balloon.Logger.Info().Dur("interval", balloon.Interval).Msg("Running at every given interval.")
	balloon.Logger.Info().Uint64("freeAllowance", balloon.FreeAllowance).Uint64("chunkSize", balloon.MemoryChunk).Msg("Running with given memory settings.")

	if balloon.DryRun {
		balloon.Logger.Info().Msg("Running in dry run mode.")
	}

	ticker := time.NewTicker(balloon.Interval)
	for range ticker.C {
		n, err := balloon.libvirt.ConnectNumOfDomains()
		if err != nil {
			balloon.Logger.Err(err).Msg("Error getting the number of active domains.")
			continue
		}
		domains, _, err := balloon.libvirt.ConnectListAllDomains(n, libvirt.ConnectListDomainsRunning)
		if err != nil {
			balloon.Logger.Err(err).Msg("Error getting running domains.")
			continue
		}
		t1 := time.Now()
		for _, domain := range domains {
			balloon.ProcessDomain(domain)
		}
		t2 := time.Now()
		balloon.Logger.Debug().Dur("timeTaken", t2.Sub(t1)).Int("domains", len(domains)).Msg("Processed all domains.")
	}
}
