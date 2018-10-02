package db


type Song struct {
	Id int
	Title string
	Author string
	Genre string
	KeyTone string
	KeyMode string
	Content []byte
}
