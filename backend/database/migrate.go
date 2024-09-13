package database

import (
	"backend/pkg/common"
)

type MigrationTask struct {
	Model           common.Entity
	DropOnFlush     bool
	TruncateOnFlush bool
}

type JoinTableMigrationTask struct {
	Model           common.Entity
	Property        string
	JoinTableStruct interface{}
}

var migrationTasks []MigrationTask

var joinTableMigrationTasks []JoinTableMigrationTask

func RegisterModel(task *MigrationTask) {
	migrationTasks = append(migrationTasks, *task)
}

func RegisterJoinTable(task *JoinTableMigrationTask) {
	joinTableMigrationTasks = append(joinTableMigrationTasks, *task)
}

func Migrate() error {

	models := []interface{}{}
	for _, task := range migrationTasks {
		models = append(models, task.Model)
	}

	for _, task := range joinTableMigrationTasks {
		models = append(models, task.JoinTableStruct)
	}

	err := DB.AutoMigrate(models...)
	if err != nil {
		return err
	}

	for _, task := range joinTableMigrationTasks {
		DB.SetupJoinTable(task.Model, task.Property, task.JoinTableStruct)
	}

	return nil
}
