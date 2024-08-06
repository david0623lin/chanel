package chanel

import (
	"chanel/structs"
	"errors"

	"gorm.io/gorm"
)

type Admins struct {
	wDB *gorm.DB
	rDB *gorm.DB
}

func AdminsInit(w, r *gorm.DB) *Admins {
	return &Admins{
		wDB: w,
		rDB: r,
	}
}

func (m *Admins) GetAdmin(datas structs.AdminLoginRequest) (structs.ChanelModelAdmins, error) {
	var res structs.ChanelModelAdmins

	query := m.rDB.Table("admins")
	query.Where("account = ?", datas.Account)

	if err := query.First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, errors.New("RecordNotFound")
		}
		return res, err
	}
	return res, nil
}

func (m *Admins) CreateAdmin(datas structs.AdminRegisterRequest) error {
	if err := m.wDB.Table("admins").Create(&datas).Error; err != nil {
		return err
	}
	return nil
}
