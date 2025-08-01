package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ModelProvider holds the schema definition for the ModelProvider entity.
type ModelProvider struct {
	ent.Schema
}

func (ModelProvider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "model_providers",
		},
	}
}

// Fields of the ModelProvider.
func (ModelProvider) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.String("api_base").NotEmpty(),
		field.Int("priority").Default(0),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ModelProvider.
func (ModelProvider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("models", ModelProviderModel.Type),
	}
}
