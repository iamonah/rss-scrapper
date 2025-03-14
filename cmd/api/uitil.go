package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any, statusCode int, header http.Header) {
	jsonByte, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)

	}
	jsonByte = append(jsonByte, '\n')

	maps.Copy(w.Header(), header)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonByte)
}

func ReadJSON(r *http.Request, location any) error {
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(location)
	if err != nil {
		return fmt.Errorf("error parsing JSON, %s", err)
	}
	return nil
}
