package telegram


type Result struct {
	Update_id int     `json:"update_id"`
	Message   Message `json:"message"`
}

type GetUpdatesRes struct {
	Ok          bool     `json:"ok"`
	Results     []Result `json:"result"`
	Description string
}

type Chat struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type From struct {
	Id       int
	Username string
}

type Message struct {
	Id int    `json:"message_id"`
	From       From   `json:"from"`
	Chat       Chat   `json:"chat"`
	Text       string `json:"text"`
}


type Update struct {
	Id  int64   `json:"update_id"`
	Msg Message `json:"message"`
}
