package chanel

import (
	"chanel/structs"

	"gorm.io/gorm"
)

type CronRecords struct {
	wDB *gorm.DB
	rDB *gorm.DB
}

func CronRecordsInit(w, r *gorm.DB) *CronRecords {
	return &CronRecords{
		wDB: w,
		rDB: r,
	}
}

func (m *CronRecords) Create(datas structs.ChanelModelCronRecords) error {
	if err := m.wDB.Create(&datas).Error; err != nil {
		return err
	}
	return nil
}
