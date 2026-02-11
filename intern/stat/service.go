package stat

import (
	"demo-1/pkg/event"
	"log"
)

// сервис - кухня, комбинирует канал передающий инфу о событиях и репозиторий (для взаимодействия с бд)
type StatServiceWithDeps struct {
	Repo  *StatRepository
	Event *event.EventBus
}

type StatService struct {
	Repo  *StatRepository
	Event *event.EventBus
}

func NewStatService(deps StatServiceWithDeps) *StatService {
	return &StatService{
		Repo:  deps.Repo,
		Event: deps.Event,
	}
}

func (serv *StatService) AddClick() {
	//проходимся циклам по приходящим в канал моделям (функция - работает как горутина)
	for msg := range serv.Event.Subscribe() {
		if msg.Type == event.EventLinkVisited {
			id, ok := msg.Data.(uint)
			if !ok {
				log.Fatalln("Bad LinkVisitedData: ", msg.Data)
				continue
			}
			//заносим айди переданный из event в бд через репозиторий
			serv.Repo.AddTicket(id)
		}
	}
}
