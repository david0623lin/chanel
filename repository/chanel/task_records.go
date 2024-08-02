package chanel

import (
	"chanel/structs"

	"gorm.io/gorm"
)

type TaskRecords struct {
	wDB *gorm.DB
	rDB *gorm.DB
}

func TaskRecordsInit(w, r *gorm.DB) *TaskRecords {
	return &TaskRecords{
		wDB: w,
		rDB: r,
	}
}

func (m *TaskRecords) GetTasks() ([]structs.ChanelModelTaskRecords, error) {
	var res []structs.ChanelModelTaskRecords

	query := m.rDB.Raw(
		`
			SELECT * FROM tasks
			ORDER BY id ASC
		`,
	)

	if err := query.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (m *TaskRecords) Create(datas structs.ChanelModelTaskRecords) error {
	if err := m.wDB.Create(&datas).Error; err != nil {
		return err
	}
	return nil
}
