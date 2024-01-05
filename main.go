/*
 * This file was last modified at 2024.01.04 15:07 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * main.go
 * $Id$
 * Задача: необходимо реализовать очень быстрый математический микросервис на языке программирования Go
 *
 * Данный сервис должен:
 *
 * 1) Работать без внешних зависимостей (СУБД, кэширование, очереди, API).
 * 2) Неограниченно масштабироваться (даже на уровне DNS).
 * 3) Очень быстро работать.
 * 4) Поддерживать очень лёгкий деплой.
 * 5) Иметь простой и компактный код, который можно написать за 20 минут на собеседовании.
 *
 */
//!+
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Robot struct {
	mass     float64
	velocity float64
}

type ResponseOk struct {
	KineticEnergy float64 `json:"kinetic_energy"`
}

func (r *Robot) getKineticEnergy() float64 {
	return (r.mass * math.Pow(r.velocity, 2)) * 0.5
}

func jsonMarshalError(err error) []byte {
	masqErrorMessage := strings.Replace(err.Error(), `"`, `\"`, -1)
	return []byte(fmt.Sprintf(`{"error":"%s"}`, masqErrorMessage))
}

func writeError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonMarshalError(err))
}

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	mass, err := strconv.ParseFloat(vars["mass"], 64)
	if err != nil {
		writeError(w, err)
		return
	}
	velocity, err := strconv.ParseFloat(vars["velocity"], 64)
	if err != nil {
		writeError(w, err)
		return
	}
	bot := Robot{mass: mass, velocity: velocity}
	energy := bot.getKineticEnergy()
	response := ResponseOk{KineticEnergy: energy}
	responseJson, err := json.Marshal(response)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseJson); err != nil {
		log.Println(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/{mass:[0-9]?\\.?[0-9]+}/{velocity:[0-9]?\\.?[0-9]+}", handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
