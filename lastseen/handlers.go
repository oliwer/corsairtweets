// Tells you when someone was last seen in the channel.
package lastseen

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

var (
	channels     *channelsCache
	ignoredNicks map[string]struct{}
)

type userStat struct {
	FirstSeen    int64  `json:"first"`
	LastSeen     int64  `json:"last"`
	LastMessage  string `json:"msg"`
	MessageCount int64  `json:"cnt"`
	nick         string // only used when sorting
}

type channelStat map[string]*userStat

func (cs channelStat) getUser(name string, create bool) *userStat {
	user, exists := cs[name]
	if !exists && create {
		user = new(userStat)
		user.FirstSeen = time.Now().Unix()
		cs[name] = user
	}
	return user
}

type channelsCache struct {
	sync.RWMutex
	cache map[string]channelStat
}

func (cc *channelsCache) getChannel(name string) channelStat {
	cc.RLock()
	channel, exists := cc.cache[name]
	cc.RUnlock()
	if !exists {
		cc.Lock()
		channel = make(channelStat)
		cc.cache[name] = channel
		cc.Unlock()
	}
	return channel
}

func (cc *channelsCache) dump(name string) {
	filename := fmt.Sprintf("lastSeen_%s.json", strings.TrimPrefix(name, "#"))
	fh, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Println("lastseen.dump:", err)
		return
	}
	defer fh.Close()
	json.NewEncoder(fh).Encode(cc.getChannel(name))
}

func (cc *channelsCache) restore() {
	cc.Lock()
	defer cc.Unlock()

	entries, err := os.ReadDir(".")
	if err != nil {
		log.Println("lastseen.restore:", err)
		return
	}
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		elements := strings.SplitN(ent.Name(), "_", 2)
		if elements[0] != "lastseen" {
			continue
		}
		channel := "#" + strings.TrimSuffix(elements[1], ".json")
		fh, err := os.Open(ent.Name())
		if err != nil {
			log.Println("lastseen.restore:", err)
			continue
		}
		defer fh.Close()
		var cs channelStat
		if err = json.NewDecoder(fh).Decode(&cs); err != nil {
			log.Printf("lastseen.restore: %s: %s\n", ent.Name(), err.Error())
			continue
		}
		cc.cache[channel] = cs
		log.Println("lastseen: restored cache for ", channel)
	}
}

func (cc *channelsCache) saveRegularly(freq time.Duration) {
	for {
		time.Sleep(freq)
		for channel := range cc.cache {
			cc.dump(channel)
		}
	}
}

func loadIgnoredNicks(ignored []string) {
	ignoredNicks = make(map[string]struct{}, len(ignored))
	for _, nick := range ignored {
		ignoredNicks[nick] = struct{}{}
	}
}

func onMessage(channel, nick, text string, date time.Time) {
	if _, found := ignoredNicks[nick]; found {
		return
	}
	// Only monitor public messages.
	if channel[0] != '#' {
		return
	}
	user := channels.getChannel(channel).getUser(nick, true)
	user.LastSeen = date.Unix()
	if len(text) > 0 {
		user.LastMessage = text
	}
	user.MessageCount += 1
}

func seen(e *bot.Event) {
	channel := e.Target
	nick := strings.TrimSpace(e.Args)
	if len(nick) == 0 {
		e.Bot.Privmsg(channel, "seen who???")
		return
	}
	if strings.Contains(nick, " ") {
		return
	}

	user := channels.getChannel(channel).getUser(nick, false)
	if user == nil {
		e.Bot.Privmsg(channel, fmt.Sprintf("I have never seen any «%s» ¯\\_(ツ)_/¯", nick))
		return
	}

	laps := time.Now().Sub(time.Unix(user.LastSeen, 0)).String()
	firstSeen := time.Unix(user.FirstSeen, 0).Format("2006/01/02 15:04 MST")
	lastMessage := ""
	if len(user.LastMessage) > 0 {
		lastMessage = fmt.Sprintf(" Last message: <%s> %s", nick, user.LastMessage)
	}

	pronoun := "He"
	// This is good enough for #Corsair 8)
	if strings.Contains(nick, "chick") {
		pronoun = "She"
	}

	e.Bot.Privmsg(channel, fmt.Sprintf("%s was last seen %s ago. %s wrote %d messages since %s.%s",
		nick, laps, pronoun, user.MessageCount, firstSeen, lastMessage))
}

func seenTop(e *bot.Event) {
	channel := channels.getChannel(e.Target)

	// Build a sorted list of the channel users.
	users := make([]*userStat, 0, len(channel))
	for nick, stat := range channel {
		stat.nick = nick
		users = append(users, stat)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].LastSeen < users[j].LastSeen
	})

	// How many entries to display.
	n, err := strconv.ParseUint(strings.TrimSpace(e.Args), 10, 8)
	if err != nil {
		// *sigh* Let's just use a default value.
		n = 5
	}

	// Print the top n.
	var msg strings.Builder
	fmt.Fprintf(&msg, "Older %d seen users:", n)
	for i, u := range users {
		if i >= int(n) {
			break
		}
		lastSeen := time.Unix(u.LastSeen, 0).Format("2006/01/02 15:04")
		fmt.Fprintf(&msg, "  %d. %s (%s)", i+1, u.nick, lastSeen)
	}
	e.Bot.Privmsg(e.Target, msg.String())
}

// Register registers the plugin with a bot.
// Use ignore to provide a list of nicks to ignore.
func Register(b bot.Bot, ignore []string) {
	channels = &channelsCache{cache: make(map[string]channelStat)}
	channels.restore()
	go channels.saveRegularly(5 * time.Minute)
	loadIgnoredNicks(ignore)
	b.Conn().HandleFunc("privmsg", func(conn *client.Conn, line *client.Line) {
		onMessage(line.Args[0], line.Nick, line.Args[1], line.Time)
	})
	b.Conn().HandleFunc("join", func(conn *client.Conn, line *client.Line) {
		onMessage(line.Args[0], line.Nick, "", line.Time)
	})
	b.Commands().Add("seen", bot.Command{
		Help:    "find when a user was last seen",
		Handler: func(e *bot.Event) { seen(e) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
	b.Commands().Add("seentop", bot.Command{
		Help:    "show top N of the oldest seen users",
		Handler: func(e *bot.Event) { seenTop(e) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
