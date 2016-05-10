package model

import (
	"time"
)

type Event struct {
	ID         uint64    `db:"id"`
	Uid        uint64    `db:"uid"`
	Cid        uint64    `db:"cid"`
	AppName    string    `db:"appname"`
	Active     bool      `db:"active"`
	CreateTime time.Time `db:"createtime"`
	EndTime    time.Time `db:"endtime"`
	Cpus       float64   `db:"cpus"`
	Mem        float64   `db:"mem"`
	Instances  uint32    `db:"instances"`
	TimeLen    string    `json:"timelen"`
}

type Message struct {
	Id        string
	ClusterId string
	JobType   string
	Target    string
	Timestamp int64
	Meta      string
	Method    string
	Header    map[string]string
	Path      string
	Uri       string
	ReplyTo   string
	Task      map[string]interface{}
}
