package helper

import (
	"encoding/json"
	"go-axiata/model"
	"net/http"
	"strconv"
)

func RespondJSON(w http.ResponseWriter, code int, status bool, msg string, info interface{}, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := model.Response{
		Status:  status,
		Message: msg,
		Info:    info,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func AddConditions(query string, req model.Request) (string, []any) {
	args := []any{}

	if req.Tag != "" {
		query += ` join post_tags pt on pt.post_id = p.id 
		join tags t on t.id = pt.tag_id 
		where LOWER(t.label) = LOWER($` + strconv.Itoa(len(args)+1) + `)`
		args = append(args, req.Tag)
	}

	page, _ := strconv.Atoi(req.Page)
	limit, _ := strconv.Atoi(req.Limit)

	query += ` order by p.created_at desc`

	offset := (page - 1) * limit
	query += ` LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	return query, args
}

func ParsePaginationSearch(filter model.Request) model.Request {
	if filter.Limit == "" {
		filter.Limit = "10"
	}

	if filter.Page == "" {
		filter.Page = "1"
	}

	return filter
}
