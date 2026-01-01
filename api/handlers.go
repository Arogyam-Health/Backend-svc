package api

import (
	"backend-service/internal/cache"
	"encoding/json"
	"log"
	"net/http"
)

func MediaHandler(store *cache.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		w.Header().Set("Content-Type", "application/json")

		ids := q.Get("ids")
		log.Print("ids:", ids)
		if ids != "" {
			idlst := make([]string, 0, 4)
			start := 0
			for i := 0; i < len(ids); i++ {
				if ids[i] == ',' {
					idlst = append(idlst, ids[start:i])
					start = i + 1
				}
			}
			if start <= len(ids)-1 {
				idlst = append(idlst, ids[start:])
			}
			json.NewEncoder(w).Encode(store.GetByIDs(idlst))
			return
		}

		if ids := q.Get("ids"); ids != "" {
			idlst := make([]string, 0, len(ids))
			for _, id := range ids {
				idlst = append(idlst, string(id))
			}
			json.NewEncoder(w).Encode(store.GetByIDs(idlst))
			return
		}

		json.NewEncoder(w).Encode(store.GetAllMedia())
	}
}

// create an /ready hadnler to check if the service is up
func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
