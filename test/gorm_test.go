package test

import (
	"sweet/scripts"
	"testing"
)

func TestGorm(t *testing.T) {
	config := scripts.GenConfig{
		DSN:             "root:wui11413@tcp(localhost:3306)/sweet?charset=utf8mb4&parseTime=True&loc=Local",
		OutPath:         "/Users/aikzy/Desktop/sweet/internal/models/query",
		ModelPkgPath:    "/Users/aikzy/Desktop/sweet/internal/models/entity",
		WithUnitTest:    false,
		WithQueryFilter: false,
		TablePrefix:     "sw_",
		SingularTable:   false,
	}
	generator, err := scripts.NewGenerator(&config)
	if err != nil {
		t.Fatal(err)
	}

	generator.GenerateModelsWithRelations()
}
