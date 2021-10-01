package main

import (
	"net/http"
)

type Reply struct {
	Products []Product `json:"products"`
}

type Product struct {
	Name string `json:"Title"`
	Title string `json:"handle"`
	Variants []struct{
		ProductId int `json:"product_id"`
		Price string `json:"price"`
		Size string `json:"option1"`
		Available bool `json:"available"`
	} `json:"variants"`
}

type Task struct {
	Client http.Client
	SizeMap map[float64]bool
}

func createTask() *Task {
	return &Task{
		Client: http.Client{},
		SizeMap: make(map[float64]bool),
	}
}

type ProxyConfig struct {
	ProxyArray []string `json:"ProxyArray"`
}

