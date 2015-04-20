// bitfmt.go
package main

import (
	"fmt"
)

//KiB and others are binary powers of byte.
const (
	KiB float32 = 1024
	MiB         = 1024 * 1024
	GiB         = 1024 * 1024 * 1024
	TiB         = 1024 * 1024 * 1024 * 1024
	PiB         = 1024 * 1024 * 1024 * 1024 * 1024
)

func fmtbytes(b float32) string {
	if b > 10*PiB {
		return fmt.Sprintf("%.0f PiB", b/PiB)
	}
	if b > 10*TiB {
		return fmt.Sprintf("%.0f TiB", b/TiB)
	}
	if b > 10*GiB {
		return fmt.Sprintf("%.0f GiB", b/GiB)
	}
	if b > 10*MiB {
		return fmt.Sprintf("%.0f MiB", b/MiB)
	}
	if b > 10*KiB {
		return fmt.Sprintf("%.0f KiB", b/KiB)
	}
	return fmt.Sprintf("%.0f B", b)
}
