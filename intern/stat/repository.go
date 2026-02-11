package stat

import (
	"demo-1/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	Db *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddTicket(linkId uint) {
	//если нет статистики за сегодня - создаем
	var stat Stat
	currentData := datatypes.Date(time.Now())
	repo.Db.Find(&stat, "link_id = ? and date = ?", linkId, currentData)
	if stat.ID == 0 {
		repo.Db.Create(&Stat{
			Clicks: 1,
			Date:   currentData,
			LinkId: linkId,
		})
	} else {
		stat.Clicks += 1
		repo.Db.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(from time.Time, to time.Time, by string) []GetStatResponse {
	var response []GetStatResponse
	var selectSql string

	switch by {
	case GroupByDay:
		selectSql = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectSql = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}

	repo.Db.Table("stats").
		Select(selectSql).
		Where("date BETWEEN ? and ?", from, to).
		Group("period").
		Order("period").
		Scan(&response)
	return response
}
