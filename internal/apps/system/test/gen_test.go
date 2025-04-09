package test

import (
	"sweet/internal/apps/system/cmd"
	"testing"
)

func TestGen(t *testing.T) {
	if gen, err := cmd.InitGen(&cmd.GenConfig{
		DSN:             "root:wui11413@tcp(127.0.0.1:3306)/sweet_system?charset=utf8mb4&parseTime=True&loc=Local",
		OutPath:         "/Users/aikzy/Desktop/go/sweet/internal/apps/system/types/query",
		ModelPkgPath:    "/Users/aikzy/Desktop/go/sweet/internal/apps/system/types/entity",
		WithUnitTest:    false,
		WithQueryFilter: true,
		TablePrefix:     "sys_",
		SingularTable:   false,
	}); err != nil {
		t.Fatal(err)
	} else {
		gen.SetupModelRelations()
	}
}
