package models

type Message_Mysql struct {
	Key     string
	Keyform string
	Value1  string
	Value2  string
}

type Add_Training_Time struct {
	Key string
	Value1 string
	Value2 string
}

type Prize struct {
	Uuid   string
	Name   string
	Awards string
	Time   string
}
