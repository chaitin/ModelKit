package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ModelProviderModel holds the schema definition for the ModelProviderModel entity.
type ModelProviderModel struct {
	ent.Schema
}

func (ModelProviderModel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "model_provider_models",
		},
	}
}

// Fields of the ModelProviderModel.
func (ModelProviderModel) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").NotEmpty(),
		field.String("provider_id").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ModelProviderModel.
func (ModelProviderModel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("provider", ModelProvider.Type).Ref("models").Unique().Field("provider_id"),
	}
}
