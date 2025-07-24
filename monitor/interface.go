package monitor

// MonitorDetector defines the interface for monitor detection implementations
type MonitorDetector interface {
	DetectMonitors() (*DetectedMonitors, error)
}

// Ensure both detector types implement the interface
var _ MonitorDetector = (*Detector)(nil)       // Shell-based detector
var _ MonitorDetector = (*NativeDetector)(nil) // Native X11 detector
