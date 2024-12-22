// Package cli contains functions for running checks from command-line.
package cli

import (
	"fmt"

	"github.com/jreisinger/checkip/check"
)

// ExtPrintSummary add IpAddrInfo.Summary for IOCLoc check
func (rs Checks) ExtPrintSummary() string {
	res := ""
	for _, r := range rs {
		// To avoid "invalid memory address or nil pointer dereference"
		// runtime error and printing empty summary info.
		if r.IpAddrInfo == nil || r.IpAddrInfo.Summary() == "" {
			continue
		}

		if r.Type == check.Info || r.Type == check.InfoAndIsMalicious {
			fmt.Printf("%-15s %s\n", r.Description, r.IpAddrInfo.Summary())
			if r.Description == "IOCLoc" {
				res = r.IpAddrInfo.Summary()
			}

		}

	}
	return res
}
