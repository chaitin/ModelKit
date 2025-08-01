package report

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/chaitin/ModelKit/backend/config"
	"github.com/chaitin/ModelKit/backend/pkg/aes"
	"github.com/chaitin/ModelKit/backend/pkg/machine"
	"github.com/chaitin/ModelKit/backend/pkg/request"
	"github.com/chaitin/ModelKit/backend/pkg/version"
)

type Reporter struct {
	client    *request.Client
	version   *version.VersionInfo
	logger    *slog.Logger
	IDFile    string
	machineID string
	cfg       *config.Config
}

func NewReport(logger *slog.Logger, cfg *config.Config, version *version.VersionInfo) *Reporter {
	raw := "https://baizhi.cloud"
	u, _ := url.Parse(raw)
	client := request.NewClient(u.Scheme, u.Host, 30*time.Second)

	r := &Reporter{
		client:  client,
		logger:  logger.With("module", "reporter"),
		IDFile:  "/app/static/.machine_id",
		cfg:     cfg,
		version: version,
	}
	if _, err := r.readMachineID(); err != nil {
		r.logger.With("error", err).Warn("read machine id file failed")
	}
	return r
}

func (r *Reporter) readMachineID() (string, error) {
	data, err := os.ReadFile(r.IDFile)
	if err != nil {
		return "", err
	}
	r.machineID = string(data)
	return r.machineID, nil
}

func (r *Reporter) Report(index string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encrypt, err := aes.Encrypt([]byte(r.cfg.DataReport.Key), string(b))
	if err != nil {
		return err
	}

	req := map[string]any{
		"index": index,
		"id":    uuid.NewString(),
		"data":  encrypt,
	}

	if _, err := request.Post[map[string]any](r.client, "/api/public/data/report", req); err != nil {
		r.logger.With("error", err).Warn("report installation failed")
		return err
	}

	return nil
}

func (r *Reporter) GetMachineID() string {
	return r.machineID
}

func (r *Reporter) ReportInstallation() error {
	if r.machineID != "" {
		return nil
	}

	id, err := machine.GenerateMachineID()
	if err != nil {
		r.logger.With("error", err).Warn("generate machine id failed")
		return err
	}
	r.machineID = id

	f, err := os.Create(r.IDFile)
	if err != nil {
		r.logger.With("error", err).Warn("create machine id file failed")
		return err
	}
	defer f.Close()

	_, err = f.WriteString(id)
	if err != nil {
		r.logger.With("error", err).Warn("write machine id file failed")
		return err
	}

	return r.Report("ModelKitnstallation", InstallData{
		MachineID: id,
		Version:   r.version.Version(),
		Timestamp: time.Now().Format(time.RFC3339),
		Type:      "installation",
	})
}
