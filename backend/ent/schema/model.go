package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"

	"github.com/chaitin/ModelKit/backend/consts"
)

// Model holds the schema definition for the Model entity.
type Model struct {
	ent.Schema
}

func (Model) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "modelkit_models",
		},
	}
}

// Fields of the Model.
func (Model) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}),
		field.String("model_name"),
		field.String("model_type").GoType(consts.ModelType("")),
		field.String("api_base"),
		field.String("api_key"),
		field.String("api_version").Optional(),
		field.String("api_header").Optional(),
		field.String("provider").GoType(consts.ModelProvider("")),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Model.
func (Model) Edges() []ent.Edge {
	return nil
}
