package link

import (
	"demo-1/configs"
	"demo-1/pkg/event"
	"demo-1/pkg/handle"
	"demo-1/pkg/middleware"
	"demo-1/pkg/res"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

type LinkHandlerWithDeps struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
	Config         *configs.Config
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerWithDeps) {
	handler := LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	//указываем в адрес изменяемый id, который выташем и в методе будем работать с нужной ссылкой
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.HandleFunc("GET /link", handler.GetAll())

}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//заносим пришедший запрос в переменную если он подходит по полям и проходит пр валидации
		payload, err := handle.HandleBody[LinkCreateRequest](w, r)
		if err != nil {
			return
		}
		//service
		//доделываем дто генерируя случайный хэш
		link := NewLink(payload.Url)
		//если вдруг сгенерированный хэш повторяктся генерим новый
		for {
			find, _ := handler.LinkRepository.GetByHash(link.Hash)
			if find == nil {
				break
			}
			link.GenerateHash()
		}
		//заносим в бд
		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdLink, 201)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//получаем информацию о емэил переданый через контекст в middleware
		email, ok := r.Context().Value("ContextEmailKey").(string)
		if ok {
			fmt.Println(email)
		}

		body, err := handle.HandleBody[LinkUpdateRequest](w, r)
		if err != nil {
			return
		}
		//достаем из строки запроса id чтоб работать с нужной ссылкрй
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//работаем с бд заменяя данные по id
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, link, 201)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, nil, 201)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
		}
		links := handler.LinkRepository.GetLinks(offset, limit)
		count := handler.LinkRepository.GetCount()
		data := GetAllLinksResponse{
			Links: links,
			Count: count,
		}
		res.Json(w, data, 201)
	}
}
