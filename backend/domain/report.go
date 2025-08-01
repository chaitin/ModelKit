package domain

import (
	"context"

	"github.com/chaitin/ModelKit/backend/consts"
)

type ReportUsecase interface {
}

type ReportRepo interface {
	GetAdminCount(ctx context.Context) (int64, error)
	GetMemberCount(ctx context.Context) (int64, error)
	GetLast24HoursStats(ctx context.Context) (*ActivityStats, error)
	GetCurrentModels(ctx context.Context) ([]*ReportModelUsage, error)
}

type ReportData struct {
	Timestamp     string              `json:"timestamp"`      // 上报时间戳
	Version       string              `json:"version"`        // 系统版本号
	MachineID     string              `json:"machine_id"`     // 机器ID
	AdminCount    int64               `json:"admin_count"`    // 管理员数量
	MemberCount   int64               `json:"member_count"`   // 成员数量
	Last24Hours   *ActivityStats      `json:"last_24_hours"`  // 最近24小时统计
	CurrentModels []*ReportModelUsage `json:"current_models"` // 当前使用的模型列表
}

type ActivityStats struct {
	ChatCount       int64   `json:"chat_count"`       // 对话次数
	CompletionCount int64   `json:"completion_count"` // 补全生成次数
	AcceptedCount   int64   `json:"accepted_count"`   // 补全采纳次数
	TotalCodeLines  int64   `json:"total_code_lines"` // 总代码量（行数）
	AcceptanceRate  float64 `json:"acceptance_rate"`  // 采纳率 (accepted_count / completion_count)
}

type ReportModelUsage struct {
	ModelID   string           `json:"model_id"`   // 模型ID
	ModelName string           `json:"model_name"` // 模型名称
	ModelType consts.ModelType `json:"model_type"` // 模型类型 (llm/coder)
}
