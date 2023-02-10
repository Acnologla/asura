package engine

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type ArenaResult = int

const (
	TimeExceeded ArenaResult = iota
	ArenaWin
	ArenaLose
	ArenaTie
)

type Finder struct {
	ID        disgord.Snowflake
	Username  string
	C         chan ArenaResult
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
func AddToMatchMaking(u *disgord.User, lastFight disgord.Snowflake, message *disgord.Message) chan ArenaResult {
	waitingQueueMutex.Lock()
	defer waitingQueueMutex.Unlock()
	i := isInMatchMaking(u.ID)
	c := make(chan ArenaResult)
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
	ctx := context.Background()
	galo := database.User.GetUser(ctx, first.ID, "Galos")
	advGalo := database.User.GetUser(ctx, second.ID, "Galos")
	rooster := rinha.GetEquippedGalo(&galo)
	advRooster := rinha.GetEquippedGalo(&advGalo)
	authorLevel := rinha.CalcLevel(rooster.Xp)
	advLevel := rinha.CalcLevel(advRooster.Xp)
	winner, _ := ExecuteRinha(nil, handler.Client, RinhaOptions{
		GaloAuthor:  &galo,
		GaloAdv:     &advGalo,
		AdvName:     rinha.GetName(second.Username, *advRooster),
		AuthorName:  rinha.GetName(first.Username, *rooster),
		AdvLevel:    advLevel,
		AuthorLevel: authorLevel,
		MessageID:   [2]*disgord.Message{first.Message, second.Message},
		IDs:         [2]disgord.Snowflake{first.ID, second.ID},
	}, false)
	database.User.UpdateUser(ctx, first.ID, func(u entities.User) entities.User {
		u.ArenaLastFight = second.ID
		return u
	})
	database.User.UpdateUser(ctx, second.ID, func(u entities.User) entities.User {
		u.ArenaLastFight = first.ID
		return u
	})
	if winner == -1 {
		first.C <- ArenaTie
		second.C <- ArenaTie
	}
	if winner == 0 {
		first.C <- ArenaWin
		second.C <- ArenaLose
	}
	if winner == 1 {
		first.C <- ArenaLose
		second.C <- ArenaWin
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
		ticker := time.NewTicker(time.Minute * 3)
		for range ticker.C {
			waitingQueueMutex.Lock()
			toFilter := []int{}
			for i, finder := range waitingQueue {
				if time.Since(finder.Timestamp).Minutes() >= 1 {
					finder.C <- TimeExceeded
					toFilter = append(toFilter, i)
				}
			}
			waitingQueue = SpliceQueue(waitingQueue, toFilter)
			waitingQueueMutex.Unlock()
		}
	}()
}
