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

func (m *Crons) GetCrons() ([]structs.ChanelModelCrons, error) {
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
