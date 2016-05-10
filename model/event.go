package model

import (
	"time"
)

type Event struct {
	ID         uint64    `json:"id" db:"id"`
	Uid        uint64    `json:"uid" db:"uid"`
	Cid        uint64    `json:"cid" db:"cid"`
	AppName    string    `json:"appname" db:"appname"`
	Active     bool      `json:"active" db:"active"`
	CreateTime time.Time `json:"createtime" db:"createtime"`
	EndTime    time.Time `json:"endtime" db:"endtime"`
	Cpus       float64   `json:"cpus" db:"cpus"`
	Mem        float64   `json:"mem" db:"mem"`
	Instances  uint32    `json:"instances" db:"instances"`
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
