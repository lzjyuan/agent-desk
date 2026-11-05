package models

import (
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

func TestAllModelsHaveDatabaseComments(t *testing.T) {
	namer := schema.NamingStrategy{TablePrefix: "t_", SingularTable: true}
	for _, model := range Models {
		if TableComment(model) == "" {
			t.Errorf("%T has no table comment", model)
		}
		parsed, err := schema.Parse(model, &sync.Map{}, namer)
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range parsed.DBNames {
			if parsed.FieldsByDBName[name].Comment == "" {
				t.Errorf("%s.%s has no column comment", parsed.Table, name)
			}
		}
	}
}
