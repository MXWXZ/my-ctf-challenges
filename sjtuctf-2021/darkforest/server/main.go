package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var engine Engine

type User struct {
	ID    int
	Name  string
	Token string
	Score int
	Conn  net.Conn
}

var user sync.Map

func handle(conn net.Conn) {
	buf := make([]byte, 512)
	l, err := conn.Read(buf)
	if err != nil {
		return
	}
	res := rdb.Get(string(buf[:l]))
	if res.Err() != nil {
		return
	}
	usr := strings.Split(res.Val(), "\n")
	if len(usr) != 3 {
		return
	}
	id, err := strconv.Atoi(usr[0])
	if err != nil {
		return
	}
	score, err := strconv.Atoi(usr[2])
	if err != nil {
		return
	}
	defer user.Delete(id)
	user.Store(id, &User{
		ID:    id,
		Name:  usr[1],
		Token: string(buf[:l]),
		Score: score,
		Conn:  conn,
	})
	fmt.Printf("New user %v login\n", id)
	conn.Write([]byte("Wait for next game!"))
	for {
		l, err = conn.Read(buf)
		if err != nil {
			return
		}
		op := strings.Split(strings.TrimSpace(string(buf[:l])), ";")
		var opchain OP
		opchain.ID = id
		last := &opchain
		iswalk := false
		isshoot := false
		if len(op) > 8 {
			return
		}
		for _, v := range op {
			switch v {
			case "tl":
				last = last.NextTurn(TLeft)
			case "tr":
				last = last.NextTurn(TRight)
			case "tu":
				last = last.NextTurn(TUp)
			case "td":
				last = last.NextTurn(TDown)
			case "w":
				if iswalk { // only one walk!
					return
				}
				last = last.NextWalk()
				iswalk = true
			default:
				if len(v) > 1 && v[0] == 's' { // shoot with item, -1 for fist
					sid, err := strconv.Atoi(v[1:])
					if err != nil {
						return
					}
					if isshoot { // only one shoot!
						return
					}
					last = last.NextShoot(sid)
					isshoot = true
				}
			}
		}
		if opchain.Next != nil {
			engine.LockOP.Lock()
			flag := false
			for _, v := range engine.OPQueue {
				if v.ID == id {
					flag = true
					break
				}
			}
			if !flag {
				engine.OPQueue = append(engine.OPQueue, &opchain)
			}
			engine.LockOP.Unlock()
		}
	}
}

var rdb *redis.Client
var rnd int = 0

func main() {
	rand.Seed(time.Now().UnixNano())
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	if err := rdb.Ping().Err(); err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(5 * time.Second)
		for {
			// bufio.NewReader(os.Stdin).ReadString('\n')
			fmt.Println("[LOG] Current round", rnd)
			p := make(map[int]string)
			user.Range(func(key, value interface{}) bool {
				v := value.(*User)
				p[v.ID] = v.Name
				return true
			})
			if len(p) > 200 {
				panic("too many!")
			}
			engine.Init(rnd, p)
			engine.Run()
			rnd++
		}
	}()

	fmt.Println("Starting the server ...")
	listener, err := net.Listen("tcp", "0.0.0.0:44444")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}
