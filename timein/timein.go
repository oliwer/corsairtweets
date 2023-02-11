// A command to get the time in a specific timezone.
package timein

import (
	"fmt"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/tkuchiki/go-timezone"
)

const timeFormat = "15:04 (or 3:04pm)"

// TODO: load these in a binary tree for faster lookups.
var tzInfos = timezone.New().TzInfos()

// Uppercase an ascii character.
func uc(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c
}

// Uppercase the first letter of each word.
func ucFirst(s string) string {
	b := []byte(s)

	b[0] = uc(b[0])

	for i := 1; i < len(b); i++ {
		if b[i] == ' ' {
			b[i+1] = uc(b[i+1])
		}
	}

	return string(b)
}

func findTimezone(s string) (bool, string) {
	if s == "" {
		return true, "UTC"
	}

	s = ucFirst(s)
	s = strings.ReplaceAll(s, " ", "_")

	for tz := range tzInfos {
		if strings.Contains(tz, s) {
			return true, tz
		}
	}

	return false, ""
}

func timeIn(e *bot.Event) {
	city := strings.TrimSpace(e.Args)

	found, tz := findTimezone(city)
	if !found {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("I do not know any timezone which matches '%s' :(", city))
		return
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error with timezone '%s': %s", tz, err.Error()))
		return
	}

	e.Bot.Privmsg(e.Target, fmt.Sprintf("Its is now %s in timezone %s",
		time.Now().In(loc).Format(timeFormat), tz))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	//nolint:errcheck
	b.Commands().Add("timein", bot.Command{
		Help:    "show time in the given timezone",
		Handler: func(e *bot.Event) { timeIn(e) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
