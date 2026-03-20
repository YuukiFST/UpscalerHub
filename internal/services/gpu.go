package services

import (
	"fmt"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

// GpuVendor identifies GPU manufacturers.
type GpuVendor string

const (
	VendorUnknown GpuVendor = "Unknown"
	VendorNVIDIA  GpuVendor = "NVIDIA"
	VendorAMD     GpuVendor = "AMD"
	VendorIntel   GpuVendor = "Intel"
)

// GpuInfo holds information about a detected GPU.
type GpuInfo struct {
	Name          string    `json:"name"`
	Vendor        GpuVendor `json:"vendor"`
	DriverVersion string    `json:"driverVersion"`
	VideoMemoryMB uint64    `json:"videoMemoryMB"`
	Icon          string    `json:"icon"`
}

type win32VideoController struct {
	Name          string
	DriverVersion string
	AdapterRAM    uint32
}

// DetectGPU returns information about the discrete (or primary) GPU.
func DetectGPU() *GpuInfo {
	var controllers []win32VideoController
	if err := wmi.Query("SELECT Name, DriverVersion, AdapterRAM FROM Win32_VideoController", &controllers); err != nil {
		return nil
	}
	if len(controllers) == 0 {
		return nil
	}

	// Prefer discrete GPU (> 2 GB VRAM) with NVIDIA > AMD > Intel priority
	var best *GpuInfo
	for _, c := range controllers {
		info := &GpuInfo{
			Name:          c.Name,
			Vendor:        detectVendor(c.Name),
			DriverVersion: c.DriverVersion,
			VideoMemoryMB: uint64(c.AdapterRAM) / (1024 * 1024),
		}
		info.Icon = vendorIcon(info.Vendor)

		if best == nil {
			best = info
			continue
		}
		// Prefer discrete (more VRAM) and higher-priority vendor
		if info.VideoMemoryMB > 2048 && (best.VideoMemoryMB <= 2048 || vendorPriority(info.Vendor) > vendorPriority(best.Vendor)) {
			best = info
		}
	}
	return best
}

func detectVendor(name string) GpuVendor {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "nvidia") || strings.Contains(lower, "geforce") ||
		strings.Contains(lower, "rtx") || strings.Contains(lower, "gtx"):
		return VendorNVIDIA
	case strings.Contains(lower, "amd") || strings.Contains(lower, "radeon"):
		return VendorAMD
	case strings.Contains(lower, "intel") || strings.Contains(lower, "iris") ||
		strings.Contains(lower, "arc"):
		return VendorIntel
	default:
		return VendorUnknown
	}
}

func vendorPriority(v GpuVendor) int {
	switch v {
	case VendorNVIDIA:
		return 3
	case VendorAMD:
		return 2
	case VendorIntel:
		return 1
	default:
		return 0
	}
}

func vendorIcon(v GpuVendor) string {
	switch v {
	case VendorNVIDIA:
		return "🟢"
	case VendorAMD:
		return "🔴"
	case VendorIntel:
		return "🔵"
	default:
		return "⚪"
	}
}

// FormatGPU returns a display string like "🟢 NVIDIA GeForce RTX 4090".
func FormatGPU(info *GpuInfo) string {
	if info == nil {
		return "⚠️ No GPU detected"
	}
	return fmt.Sprintf("%s %s", info.Icon, info.Name)
}
