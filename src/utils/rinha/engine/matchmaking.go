package engine

import (
	"asura/src/handler"
	"asura/src/utils/rinha"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type Finder struct {
	ID        disgord.Snowflake
	Username  string
	C         chan rinha.ArenaResult
	Message   *disgord.Message
	LastFight disgord.Snowflake
	Timestamp time.Time
}

var waitingQueue = []*Finder{}
var waitingQueueMutex = sync.RWMutex{}

func isInMatchMaking(id disgord.Snowflake) int {
	for i, finder := range waitingQueue {
		if finder.ID == id {
			return i
		}
	}
	return -1
}
func AddToMatchMaking(u *disgord.User, lastFight disgord.Snowflake, message *disgord.Message) chan rinha.ArenaResult {
	waitingQueueMutex.Lock()
	defer waitingQueueMutex.Unlock()
	i := isInMatchMaking(u.ID)
	c := make(chan rinha.ArenaResult)
	if i == -1 {
		waitingQueue = append(waitingQueue, &Finder{
			ID:        u.ID,
			Username:  u.Username,
			Message:   message,
			LastFight: lastFight,
			C:         c,
			Timestamp: time.Now(),
		})
	}
	return c
}

func SpliceQueue(slice []*Finder, s []int) []*Finder {
	newArr := []*Finder{}
	for i, element := range slice {
		if !rinha.IsIntInList(i, s) {
			newArr = append(newArr, element)
		}
	}
	return newArr
}

func initBattle(first, second *Finder) {
	galo, _ := rinha.GetGaloDB(first.ID)
	advGalo, _ := rinha.GetGaloDB(second.ID)
	authorLevel := rinha.CalcLevel(galo.Xp)
	advLevel := rinha.CalcLevel(advGalo.Xp)
	winner, _ := ExecuteRinha(nil, handler.Client, RinhaOptions{
		GaloAuthor:  galo,
		GaloAdv:     advGalo,
		AdvName:     rinha.GetName(second.Username, advGalo),
		AuthorName:  rinha.GetName(first.Username, galo),
		AdvLevel:    advLevel,
		AuthorLevel: authorLevel,
		MessageID:   [2]*disgord.Message{first.Message, second.Message},
		IDs:         [2]disgord.Snowflake{first.ID, second.ID},
	}, false)
	rinha.UpdateGaloDB(first.ID, func(galo rinha.Galo) (rinha.Galo, error) {
		galo.Arena.LastFight = second.ID
		return galo, nil
	})
	rinha.UpdateGaloDB(second.ID, func(galo rinha.Galo) (rinha.Galo, error) {
		galo.Arena.LastFight = first.ID
		return galo, nil
	})
	if winner == -1 {
		first.C <- rinha.ArenaTie
		second.C <- rinha.ArenaTie
	}
	if winner == 0 {
		first.C <- rinha.ArenaWin
		second.C <- rinha.ArenaLose
	}
	if winner == 1 {
		first.C <- rinha.ArenaLose
		second.C <- rinha.ArenaWin
	}
}

func FindFight() []int {
	arr := []int{0}
	firstFighter := waitingQueue[0]
	for i, finder := range waitingQueue {
		if i != 0 {
			if firstFighter.LastFight != finder.ID {
				arr = append(arr, i)
				break
			}
		}
	}
	return arr
}
func matchmaking() {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		waitingQueueMutex.Lock()
		if len(waitingQueue) > 1 {
			fighters := FindFight()
			if len(fighters) == 2 {
				first := waitingQueue[fighters[0]]
				second := waitingQueue[fighters[1]]
				waitingQueue = SpliceQueue(waitingQueue, fighters)
				go initBattle(first, second)

			}
		}
		waitingQueueMutex.Unlock()
	}

}
func init() {
	go matchmaking()
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			waitingQueueMutex.Lock()
			toFilter := []int{}
			for i, finder := range waitingQueue {
				if time.Since(finder.Timestamp).Minutes() >= 1 {
					finder.C <- rinha.TimeExceeded
					toFilter = append(toFilter, i)
				}
			}
			waitingQueue = SpliceQueue(waitingQueue, toFilter)
			waitingQueueMutex.Unlock()
		}
	}()
}
