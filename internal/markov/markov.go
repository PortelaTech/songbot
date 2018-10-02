package markov

import (
	"github.com/PortelaTech/songbot/internal/telegram"
	"github.com/garyburd/redigo/redis"
	"log"
	"regexp"
	"strings"
)

type Markov struct {
	length int
}

func (m Markov) StoreUpdates(updates []telegram.Result, c redis.Conn) {
	for _, update := range updates {
		splitted := strings.Split(update.Message.Text, " ")
		for index, word := range splitted {
			if index < len(splitted)-1 {
				c.Do("SADD", word, splitted[index+1])
			}
		}
	}
}

func (m Markov) Generate(seed string, c redis.Conn) string {
	log.Printf("seed: %s\n", seed)
	s := []string{}
	s = append(s, seed)
	splitted := strings.Split(seed, " ")
	key := splitted[len(splitted)-1]
	for i := 1; i < m.length; i++ {
		next, _ := redis.String(c.Do("SRANDMEMBER", key))
		s = append(s, next)
		matched, _ := regexp.MatchString(".*[\\.;!?¿¡]$", next)
		if next == "" || matched {
			break
		}
		key = next
	}
	text := strings.Join(s, " ")
	log.Printf("Text: %s\n", text)
	return text
}
