package chanel

import (
	"chanel/structs"

	"gorm.io/gorm"
)

type Crons struct {
	wDB *gorm.DB
	rDB *gorm.DB
}

func CronsInit(w, r *gorm.DB) *Crons {
	return &Crons{
		wDB: w,
		rDB: r,
	}
}

func (m *Crons) GetCronByID(cronID int32) (structs.ChanelModelCrons, error) {
	var res structs.ChanelModelCrons

	query := m.rDB.Raw(
		`
			SELECT * FROM crons WHERE id = ?
		`, cronID,
	)

	if err := query.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (m *Crons) GetCrons(datas structs.GetCronsRequest) ([]structs.ChanelModelCrons, error) {
	var res []structs.ChanelModelCrons

	query := m.rDB.Table("crons")

	if datas.Path != "" {
		query.Where("path LIKE ?", "%"+datas.Path+"%")
	}

	if datas.Method != "" {
		query.Where("method = ?", datas.Method)
	}

	if datas.Status != 0 {
		query.Where("status = ?", datas.Status)
	}

	query.Order("id DESC")

	if err := query.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (m *Crons) GetEnableCrons() ([]structs.ChanelModelCrons, error) {
	var res []structs.ChanelModelCrons

	query := m.rDB.Raw(
		`
			SELECT * FROM crons WHERE status = 1 ORDER BY id ASC
		`,
	)

	if err := query.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (m *Crons) CreateCron(datas structs.ChanelModelCrons) (int32, error) {
	if err := m.wDB.Create(&datas).Error; err != nil {
		return 0, err
	}
	return datas.ID, nil
}

func (m *Crons) UpdateCron(datas structs.ChanelModelCrons, cronID int32) error {
	query := m.wDB.Model(&structs.ChanelModelCrons{})
	query.Where("id = ?", cronID)
	err := query.Updates(&datas).Error

	if err != nil {
		return err
	}
	return nil
}

func (m *Crons) DeleteCronByID(cronID int32) error {
	query := m.wDB.Where("id = ?", cronID)

	if err := query.Delete(&structs.ChanelModelCrons{}).Error; err != nil {
		return err
	}
	return nil
}
