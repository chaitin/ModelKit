package repo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"

	"github.com/chaitin/ModelKit/backend/db"
	"github.com/chaitin/ModelKit/backend/db/model"
	"github.com/chaitin/ModelKit/backend/domain"
	"github.com/chaitin/ModelKit/backend/pkg/entx"
)

type ModelRepo struct {
	db *db.Client
}

func NewModelRepo(db *db.Client) domain.ModelRepo {
	return &ModelRepo{db: db}
}

func (r *ModelRepo) UpdateModel(ctx context.Context, req *domain.UpdateModelReq) (*db.Model, error) {
	var m *db.Model
	err := entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		var old *db.Model
		var err error

		// 根据不同的查找条件获取模型
		if req.ID != "" {
			// 通过 ID 查找
			modelID, parseErr := uuid.Parse(req.ID)
			if parseErr != nil {
				return parseErr
			}
			old, err = tx.Model.Query().
				Where(model.ID(modelID)).
				WithAPIConfig().
				Only(ctx)
		} else if req.ModelName != "" && req.Provider != "" {
			// 通过 ModelName + Provider 查找
			old, err = tx.Model.Query().
				Where(model.ModelName(req.ModelName), model.Provider(req.Provider)).
				WithAPIConfig().
				Only(ctx)
		} else {
			return fmt.Errorf("必须提供 ID 或者 ModelName+Provider 来查找模型")
		}

		if err != nil {
			return err
		} else {
			m = old
		}

		// 更新API配置
		if old.Edges.APIConfig != nil {
			apiConfigUpdate := tx.ModelAPIConfig.UpdateOneID(old.Edges.APIConfig.ID)

			if req.APIKey != "" {
				apiConfigUpdate.SetAPIKey(req.APIKey)
			}
			if req.APIVersion != "" {
				apiConfigUpdate.SetAPIVersion(req.APIVersion)
			}
			if req.APIHeader != "" {
				apiConfigUpdate.SetAPIHeader(req.APIHeader)
			}

			if _, err := apiConfigUpdate.Save(ctx); err != nil {
				return err
			}
		}

		return nil
	})
	return m, err
}

func (r *ModelRepo) GetModel(ctx context.Context, req *domain.GetModelReq) (*db.Model, error) {
	query := r.db.Model.Query()

	// 根据不同的查找条件获取模型
	if req.ID != "" {
		// 通过 ID 查找
		modelID, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}
		return query.Where(model.ID(modelID)).
			WithAPIConfig().
			Only(ctx)
	} else if req.ModelName != "" && req.Provider != "" {
		// 通过 ModelName + Provider 查找
		return query.Where(model.ModelName(req.ModelName), model.Provider(req.Provider)).
			WithAPIConfig().
			Only(ctx)
	} else {
		return nil, fmt.Errorf("必须提供 ID 或者 ModelName+Provider 来查找模型")
	}
}

func (r *ModelRepo) ListModel(ctx context.Context, req *domain.ListModelReq) ([]*db.Model, error) {
	query := r.db.Model.Query()

	// 添加筛选条件
	if req.ID != "" {
		if id, err := uuid.Parse(req.ID); err == nil {
			query = query.Where(model.ID(id))
		}
	}
	if req.ModelName != "" {
		query = query.Where(model.ModelName(req.ModelName))
	}
	if req.Provider != "" {
		query = query.Where(model.Provider(req.Provider))
	}
	if req.ModelType != "" {
		query = query.Where(model.ModelType(req.ModelType))
	}

	// 按创建时间降序排列并加载API配置
	query = query.Order(model.ByCreatedAt(sql.OrderDesc())).
		WithAPIConfig()

	return query.All(ctx)
}
