package repo

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"

	"github.com/chaitin/ModelKit/backend/consts"
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

func (r *ModelRepo) UpdateModel(ctx context.Context, id string, fn func(tx *db.Tx, old *db.Model, up *db.ModelUpdateOne) error) (*db.Model, error) {
	modelID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var m *db.Model
	err = entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		old, err := tx.Model.Get(ctx, modelID)
		if err != nil {
			return err
		}

		up := tx.Model.UpdateOneID(old.ID)
		if err := fn(tx, old, up); err != nil {
			return err
		}
		if n, err := up.Save(ctx); err != nil {
			return err
		} else {
			m = n
		}
		return nil
	})
	return m, err
}

func (r *ModelRepo) GetModel(ctx context.Context, modelName string, provider consts.ModelProvider) (*db.Model, error) {
	result, err := r.db.Model.Query().
		Where(model.ModelName(modelName), model.Provider(provider)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
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

	// 按创建时间降序排列
	query = query.Order(model.ByCreatedAt(sql.OrderDesc()))

	return query.All(ctx)
}
