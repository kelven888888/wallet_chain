package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//
//	type User struct {
//	  gorm.Model
//	}
type Model struct {
	Id        int        `json:"id" gorm:"primary_key" form:"id"`
	CreatedAt *LocalTime `json:"created_at"`
	UpdatedAt *LocalTime `json:"updated_at"`
	DeletedAt *LocalTime `json:"deleted_at" sql:"index"`
}

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// UnmarshalJSON 反序列化处理
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	str := strings.Trim(string(data), "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", str)
	*t = LocalTime(t1)
	return err
}
func (t LocalTime) String() string {
	return time.Time(t).Format(timeFormat)
}
