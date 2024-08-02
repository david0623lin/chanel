package chanel

import (
	"chanel/structs"

	"gorm.io/gorm"
)

type Tasks struct {
	wDB *gorm.DB
	rDB *gorm.DB
}

func TasksInit(w, r *gorm.DB) *Tasks {
	return &Tasks{
		wDB: w,
		rDB: r,
	}
}

func (m *Tasks) GetTaskByID(taskID int32) (structs.ChanelModelTasks, error) {
	var res structs.ChanelModelTasks

	query := m.rDB.Raw(
		`
			SELECT * FROM tasks WHERE id = ?
		`, taskID,
	)

	if err := query.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (m *Tasks) GetTasks(datas structs.GetTasksRequest) ([]structs.ChanelModelTasks, error) {
	var res []structs.ChanelModelTasks

	query := m.rDB.Table("tasks")

	if datas.StartTime != 0 && datas.EndTime != 0 {
		query.Where("create_time >= ? AND create_time <= ?", datas.StartTime, datas.EndTime)
	}

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

func (m *Tasks) GetNotCompletedTasks() ([]structs.ChanelModelTasks, error) {
	var res []structs.ChanelModelTasks

	query := m.rDB.Raw(
		`
			SELECT * FROM tasks WHERE status = 1 ORDER BY id ASC
		`,
	)

	if err := query.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (m *Tasks) CreateTask(datas structs.ChanelModelTasks) (int32, error) {
	if err := m.wDB.Create(&datas).Error; err != nil {
		return 0, err
	}
	return datas.ID, nil
}

func (m *Tasks) UpdateTask(datas structs.ChanelModelTasks, taskID int32) error {
	query := m.wDB.Model(&structs.ChanelModelTasks{})
	query.Where("id = ?", taskID)
	err := query.Updates(&datas).Error

	if err != nil {
		return err
	}
	return nil
}

func (m *Tasks) DeleteTaskByID(taskID int32) error {
	query := m.wDB.Where("id = ?", taskID)

	if err := query.Delete(&structs.ChanelModelTasks{}).Error; err != nil {
		return err
	}
	return nil
}
