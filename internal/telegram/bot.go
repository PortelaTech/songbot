package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/PortelaTech/songbot/internal/markov"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	Token  string
	Redis  redis.Conn
	Chance int
}

func apiTelegram(method, token string, params url.Values) ([]byte, error) {
	// 	link := "https://api.telegram.org/bot{botId}:{apiKey}/sendMessage?chat_id={chatId}&text={text}&parse_mode=Markdown"
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?%s", token, method, params.Encode())
	timeout := 30 * time.Second
	client := http.Client{ Timeout: timeout, }
	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}
	resp.Close = true
	defer resp.Body.Close()
	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return json, nil
}

func (bot *Bot) Commands(command string, chat int) {
	markov := markov.Markov{20}
	words := strings.Split(command, " ")
	seed := strings.Join(words[1:], " ") // Removes the initial command
	if words[0] == "/cmd" && len(words) >= 2 {
		text := markov.Generate(seed, bot.Redis)
		bot.SendMessage(text, chat)
	}
}


func (bot Bot) GetUpdates() ([]Result,error) {
	offset, _ := redis.String(bot.Redis.Do("GET", "update_id"))
	params := url.Values{}
	params.Set("offset", offset)
	params.Set("timeout", strconv.Itoa(30))

	res, err := apiTelegram("getUpdates", token, params)
	var getUpdatesRes GetUpdatesRes
	json.Unmarshal(res, &getUpdatesRes)
	if !getUpdatesRes.Ok {
		err = fmt.Errorf("cmd: %s\n", getUpdatesRes.Description)
		return nil,err
	}
	if len(getUpdatesRes.Results) != 0 {
		updateID := getUpdatesRes.Results[len(getUpdatesRes.Results)-1].Update_id + 1
		bot.Redis.Do("SET", "update_id", updateID)
		return getUpdatesRes.Results,nil
	}
	return nil,nil
}

func (bot Bot) SendMessage(text string, chat int) (bool, error) {
	var responseRecieved struct {
		Ok          bool
		Description string
	}
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(chat))
	params.Set("text", text)
	resp, err := apiTelegram("sendMessage", token, params)

	err = json.Unmarshal(resp, &responseRecieved)
	if err != nil {
		return false, err
	}
	if !responseRecieved.Ok {
		return false, fmt.Errorf("chobot: %s\n", responseRecieved.Description)
	}
	return responseRecieved.Ok, nil
}


func (bot Bot) Listen() {
	var err error

	rand.Seed(time.Now().UnixNano())
	bot.Chance = chance
	bot.Redis, err = redis.Dial(connection, ":" + strconv.Itoa(port))
	if err != nil {
		panic(err);
	}
	fmt.Printf("redis connection: %v, port: %v, chance: %v\n", connection, port, chance)
	bot.Poll()
}

func (bot Bot) Poll() {
	markov := Markov{10}
	for {
		updates, err := bot.GetUpdates()
		if (err != nil) {
			panic(err)
		}
		if updates != nil {
			markov.StoreUpdates(updates, bot.Redis)
			if strings.HasPrefix(updates[0].Message.Text, "/cho") {
				bot.Commands(
					updates[0].Message.Text,
					updates[0].Message.Chat.Id,
				)
			} else if rand.Intn(100) <= bot.Chance {
				in_text := updates[len(updates)-1].Message.Text
				parts := strings.Split(in_text, " ")
				seed := parts[0] // Seed the chain with the first word only

				chatId := updates[len(updates)-1].Message.Chat.Id
				out_text := markov.Generate(seed, bot.Redis)
				bot.SendMessage(out_text, chatId)

			}
		}
	}
}
