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

// ModelAPIConfig holds the schema definition for the ModelAPIConfig entity.
type ModelAPIConfig struct {
	ent.Schema
}

func (ModelAPIConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "modelkit_model_api_configs",
		},
	}
}

// Fields of the ModelAPIConfig.
func (ModelAPIConfig) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}),
		field.UUID("model_id", uuid.UUID{}),
		field.String("api_base"),
		field.String("api_key"),
		field.String("api_version").Optional(),
		field.String("api_header").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ModelAPIConfig.
func (ModelAPIConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("model", Model.Type).
			Ref("api_config").
			Field("model_id").
			Required().
			Unique(),
	}
}
