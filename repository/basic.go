package repository

import (
	"chanel/repository/chanel"

	"gorm.io/gorm"
)

type Chanel struct {
	Tasks       *chanel.Tasks
	TaskRecords *chanel.TaskRecords
	Crons       *chanel.Crons
	CronRecords *chanel.CronRecords
	Admins      *chanel.Admins
}

func ChanelInit(w, r *gorm.DB) *Chanel {
	return &Chanel{
		Tasks:       chanel.TasksInit(w, r),
		TaskRecords: chanel.TaskRecordsInit(w, r),
		Crons:       chanel.CronsInit(w, r),
		CronRecords: chanel.CronRecordsInit(w, r),
		Admins:      chanel.AdminsInit(w, r),
	}
}
