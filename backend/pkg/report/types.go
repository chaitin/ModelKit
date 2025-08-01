package report

type InstallData struct {
	MachineID string `json:"machine_id"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}
