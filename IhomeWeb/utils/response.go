package utils

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, errno, errmsg string, data interface{}) error {
	/*
		response := map[string]string{
					"Error":  utils.RECODE_MOBILEERR,
					"Errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
				}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	*/

	resp := map[string]interface{}{
		"errno":  errno,
		"errmsg": errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}
