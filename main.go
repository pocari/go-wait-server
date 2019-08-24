package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type requestParams struct {
	waitTime int
}

type responseObject struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func parseParams(r *http.Request) (*requestParams, error) {
	rp := requestParams{}
	values := r.URL.Query()
	if values != nil {
		if _, ok := values["time"]; ok {
			waitTime, err := strconv.Atoi(values["time"][0])
			if err != nil {
				return nil, fmt.Errorf("time is not valid: %s", values["time"])
			}

			if waitTime < 0 || waitTime > 30 {
				return nil, fmt.Errorf("time range is 0 <= time <= 30")
			}

			rp.waitTime = waitTime
		}

	}

	return &rp, nil
}

func wait(ctx context.Context, waitTime int) {
	log.Printf(fmt.Sprintf("waitTime: %v\n", waitTime))
	waitCh := time.After(time.Duration(waitTime) * time.Second)
	for {
		select {
		case <-waitCh:
			return
		case <-ctx.Done():
			log.Println("request canceled!")
			return
		default:
			log.Println("wait ...")
			time.Sleep(1 * time.Second)
		}
	}
}

// https://fisproject.jp/2017/01/golang-tips/
func JSONSafeMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)
	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func waitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("start handler", r)
	w.Header().Set("Content-Type", "application/json")
	resParam := responseObject{}

	reqParam, err := parseParams(r)
	if err != nil {
		resParam.Status = http.StatusBadRequest
		resParam.Message = err.Error()
		w.WriteHeader(resParam.Status)
		data, err := JSONSafeMarshal(resParam, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	wait(r.Context(), reqParam.waitTime)

	resParam.Status = http.StatusOK
	resParam.Message = "success"
	data, err := JSONSafeMarshal(resParam, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	log.Println("end handler")
	return
}

func main() {
	http.HandleFunc("/wait", waitHandler)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		panic(err)
	}
}
