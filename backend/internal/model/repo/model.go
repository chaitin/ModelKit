package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	"github.com/chaitin/ModelKit/backend/db"
	"github.com/chaitin/ModelKit/backend/domain"
	"github.com/chaitin/ModelKit/backend/pkg/entx"
)

type ModelRepo struct {
	db    *db.Client
	cache *cache.Cache
}

func NewModelRepo(db *db.Client) domain.ModelRepo {
	cache := cache.New(24*time.Hour, 10*time.Minute)
	return &ModelRepo{db: db, cache: cache}
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
