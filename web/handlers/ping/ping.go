package ping

import (
	"net/http"

	"github.com/cmelgarejo/minesweeper-svc/web/services/common"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(common.ContentTypeKey, common.AppTypeTextHTML)
	_, _ = w.Write([]byte("OK"))
}
