//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"

	"github.com/chaitin/ModelKit/backend/pkg/entx"
)

func main() {
	if err := entc.Generate(
		"./schema",
		&gen.Config{
			Target:  "../db",
			Package: "github.com/chaitin/ModelKit/backend/db",
			Features: []gen.Feature{
				gen.FeatureUpsert,
				gen.FeatureModifier,
				gen.FeatureExecQuery,
				gen.FeatureIntercept,
				gen.FeatureLock,
			},
		},
		entc.Extensions(
			&entx.Page{},
		),
	); err != nil {
		log.Fatal("running ent codegen:", err)
	}
}
