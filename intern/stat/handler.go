package stat

import (
	"demo-1/configs"
	"fmt"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandler struct {
	Repo *StatRepository
}

type StatHandlerWithDeps struct {
	Repo   *StatRepository
	Config *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerWithDeps) {
	handler := &StatHandler{
		deps.Repo,
	}
	router.HandleFunc("GET /stat", handler.GetStat())
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from param", http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invvalid to param", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
		}
		fmt.Println(from, to, by)
	}
}
