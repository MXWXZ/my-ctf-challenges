package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/copier"
)

type Engine struct {
	Size     int
	Forest   [][]*Ground
	Player   map[int]*Player
	Tick     int
	Round    int
	MsgQueue []string
	OPQueue  []*OP
	DVFPos   *DVF
	LockOP   sync.Mutex
}

var ItemNum [ITTot]int

func init() {
	// resoureces
	ItemNum[ITFlag] = 20
	ItemNum[ITUSP] = 10
	ItemNum[ITAK47] = 5
	ItemNum[ITAWP] = 1
	ItemNum[ITArmer] = 4
	ItemNum[ITBlood] = 5
	ItemNum[ITMine] = 1
	ItemNum[ITSophon] = 1
	ItemNum[ITWaterDrop] = 2
}

func (e *Engine) Init(round int, player map[int]string) {
	// return new available position
	getPos := func() (int, int) {
		n := 0
		for {
			x := rand.Intn(e.Size-2) + 1
			y := rand.Intn(e.Size-2) + 1
			if e.Forest[x][y].Type == GGround && e.Forest[x][y].Contain == nil {
				return x, y
			}
			n++
			if n >= 3 {
				for p := 1; p < e.Size-1; p++ {
					for q := 1; q < e.Size-1; q++ {
						if e.Forest[p][q].Type == GGround && e.Forest[p][q].Contain == nil {
							return p, q
						}
					}
				}
			}
		}
	}
	e.Tick = 0
	e.Round = round
	e.DVFPos = nil

	// read map file
	f, err := os.ReadFile("map.txt")
	if err != nil {
		panic(err)
	}
	fl := strings.Split(strings.TrimSpace(string(f)), "\n")
	e.Size = len(fl)
	e.Forest = make([][]*Ground, e.Size)
	for i := range e.Forest {
		e.Forest[i] = make([]*Ground, e.Size)
		for j := range e.Forest[i] {
			e.Forest[i][j] = &Ground{}
			copier.CopyWithOption(&e.Forest[i][j], defaultGround[string(fl[i][j])], copier.Option{DeepCopy: true})
		}
	}

	// init item
	for i := 0; i < int(ITTot); i++ {
		for j := 0; j < ItemNum[i]; j++ {
			x, y := getPos()
			e.Forest[x][y].Contain = defaultItem[i]
		}
	}
	if rand.Intn(100) < 30 {
		x, y := getPos()
		e.Forest[x][y].Contain = defaultItem[int(ITDVF)]
	}

	// init player
	e.Player = make(map[int]*Player)
	for k, v := range player {
		x, y := getPos()
		// default player
		e.Player[k] = &Player{
			ID:     k,
			Name:   v,
			HP:     100,
			X:      x,
			Y:      y,
			Vision: 6,
			Killed: 0,
			Armer:  0,
			Turn:   TRight,
			Bag:    []*Item{},
		}
		e.Forest[x][y].Contain = e.Player[k]
	}
}

type MsgSend struct {
	Tick   int
	HP     int
	Armer  int
	Turn   TurnType
	Vision string
	Bag    []string
	Death  string
}

func ConnSend(id int, s string) {
	val, ok := user.Load(id)
	if ok {
		go val.(*User).Conn.Write([]byte(s))
	}
}

// calculate shoot damage
func (e *Engine) fight(id int, vision int, damage int, penetrate bool) {
	p := e.Player[id]
	target := []int{}
	if p.Turn == TRight {
		for i := p.Y + 1; i < Min(p.Y+vision+1, e.Size); i++ {
			if e.Forest[p.X][i].GetType() == CPlayer {
				target = append(target, e.Forest[p.X][i].Contain.(*Player).ID)
				if !penetrate {
					break
				}
			}
		}
	}
	if p.Turn == TLeft {
		for i := p.Y - 1; i >= Max(p.Y-vision, 0); i-- {
			if e.Forest[p.X][i].GetType() == CPlayer {
				target = append(target, e.Forest[p.X][i].Contain.(*Player).ID)
				if !penetrate {
					break
				}
			}
		}
	}
	if p.Turn == TDown {
		for i := p.X + 1; i < Min(p.X+vision+1, e.Size); i++ {
			if e.Forest[i][p.Y].GetType() == CPlayer {
				target = append(target, e.Forest[i][p.Y].Contain.(*Player).ID)
				if !penetrate {
					break
				}
			}
		}
	}
	if p.Turn == TUp {
		for i := p.X - 1; i >= Max(p.X-vision, 0); i-- {
			if e.Forest[i][p.Y].GetType() == CPlayer {
				target = append(target, e.Forest[i][p.Y].Contain.(*Player).ID)
				if !penetrate {
					break
				}
			}
		}
	}
	for _, i := range target {
		e.Player[i].HP -= Max(damage-e.Player[i].Armer, 0)
		if e.checkAlive(i, fmt.Sprintf("%v(%v killed) kill %v(%v killed)", p.Name, p.Killed, e.Player[i].Name, e.Player[i].Killed)) {
			p.Killed++
		}
	}
}

// perform alive judge
func (e *Engine) checkAlive(id int, s string) bool {
	p := e.Player[id]
	if p.HP <= 0 && e.Forest[p.X][p.Y].Contain == p {
		e.Forest[p.X][p.Y].Contain = nil
		if len(p.Bag) != 0 {
			var tmp = &Item{
				Name: "P",
				Type: ITPack,
			}
			copier.CopyWithOption(&tmp.Pack, p.Bag, copier.Option{DeepCopy: true})
			tmp.StepOn = func(p *Player) {
				for _, k := range tmp.Pack {
					k.StepOn(p)
				}
			}
			e.Forest[p.X][p.Y].Contain = tmp
		}
		if p.Death == "" {
			p.Death = s
		}
		s, _ := json.Marshal(&MsgSend{
			Death: p.Death,
			Tick:  e.Tick,
		})
		ConnSend(id, string(s))
		e.MsgQueue = append(e.MsgQueue, p.Death)
		return true
	}
	return false
}

func (e *Engine) Run() {
	for i := range e.Player {
		ConnSend(i, "Game start")
	}
	time.Sleep(time.Second * 3)
	for tick := 0; tick < 1000; tick++ {
		e.Tick = tick

		// write to tmp file
		f, err := os.OpenFile(fmt.Sprintf("round/%v.tmp", rnd), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println("[ERROR]", err)
		}
		f.WriteString(strconv.Itoa(len(e.MsgQueue)))
		f.WriteString("\n")
		for _, i := range e.MsgQueue {
			f.WriteString(i + "\n")
		}
		for i := range e.Forest {
			for j := range e.Forest[i] {
				f.WriteString(e.Forest[i][j].GetName(true) + "\t")
			}
			f.WriteString("\n")
		}
		f.WriteString("\n")
		f.Close()

		e.MsgQueue = []string{}

		var ids []int
		for i := range e.Player {
			ids = append(ids, i)
		}
		rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
		// send user vision
		for _, i := range ids {
			var vision string
			p := e.Player[i]
			if p.HP <= 0 {
				continue
			}
			if p.Turn == TRight {
				for m := p.X - 1; m < p.X+2; m++ {
					for n := p.Y; n < Min(p.Y+p.Vision+1, e.Size); n++ {
						if m == p.X && n == p.Y {
							vision += "X"
							continue
						}
						if e.Forest[m][n].GetType() == CPlayer {
							vision += "E" // enermy
							continue
						}
						vision += e.Forest[m][n].GetName(p.Vision == 999)
						if p.Vision != 999 && !e.Forest[m][n].AllowSee {
							break
						}
					}
					vision += "\n"
				}
			}
			if p.Turn == TLeft {
				for m := p.X - 1; m < p.X+2; m++ {
					tmp := ""
					for n := p.Y; n >= Max(p.Y-p.Vision, 0); n-- {
						if m == p.X && n == p.Y {
							tmp = "X" + tmp
							continue
						}
						if e.Forest[m][n].GetType() == CPlayer {
							tmp = "E" + tmp // enermy
							continue
						}
						tmp = e.Forest[m][n].GetName(p.Vision == 999) + tmp
						if p.Vision != 999 && !e.Forest[m][n].AllowSee {
							break
						}
					}
					vision += tmp + "\n"
				}
			}
			if p.Turn == TDown {
				for n := p.Y - 1; n < p.Y+2; n++ {
					for m := p.X; m < Min(p.X+p.Vision+1, e.Size); m++ {
						if m == p.X && n == p.Y {
							vision += "X"
							continue
						}
						if e.Forest[m][n].GetType() == CPlayer {
							vision += "E" // enermy
							continue
						}
						vision += e.Forest[m][n].GetName(p.Vision == 999)
						if p.Vision != 999 && !e.Forest[m][n].AllowSee { // 智子有挂
							break
						}
					}
					vision += "\n"
				}
			}
			if p.Turn == TUp {
				for n := p.Y - 1; n < p.Y+2; n++ {
					tmp := ""
					for m := p.X; m >= Max(p.X-p.Vision, 0); m-- {
						if m == p.X && n == p.Y {
							tmp = "X" + tmp
							continue
						}
						if e.Forest[m][n].GetType() == CPlayer {
							tmp = "E" + tmp // enermy
							continue
						}
						tmp = e.Forest[m][n].GetName(p.Vision == 999) + tmp
						if p.Vision != 999 && !e.Forest[m][n].AllowSee {
							break
						}
					}
					vision += tmp + "\n"
				}
			}
			bag := []string{}
			for _, m := range p.Bag {
				bag = append(bag, m.GetName(false))
			}
			msg, _ := json.Marshal(&MsgSend{
				Tick:   e.Tick,
				HP:     p.HP,
				Armer:  p.Armer,
				Turn:   p.Turn,
				Vision: vision,
				Bag:    bag,
			})
			ConnSend(i, string(msg))
		}

		// for i := range e.Forest {
		// 	for j := range e.Forest[i] {
		// 		fmt.Printf("%v", e.Forest[i][j].GetName())
		// 	}
		// 	fmt.Printf("\n")
		// }

		time.Sleep(time.Second) // wait for response
		// bufio.NewReader(os.Stdin).ReadString('\n')

		// perform operation
		e.LockOP.Lock()
		for _, i := range e.OPQueue {
			p, ok := e.Player[i.ID]
			if !ok {
				continue
			}
			if p.HP <= 0 {
				continue
			}
			cur := i.Next
			x := p.X
			y := p.Y
			dx := 0
			dy := 0
			for ; cur != nil; cur = cur.Next {
				switch cur.Oper {
				case OPWalk:
					switch p.Turn {
					case TRight:
						dx = 0
						dy = 1
					case TLeft:
						dx = 0
						dy = -1
					case TUp:
						dx = -1
						dy = 0
					case TDown:
						dx = 1
						dy = 0
					}
					newpos := e.Forest[x+dx][y+dy]
					if newpos.AllowRun && newpos.GetType() != CPlayer {
						if newpos.GetType() == CItem {
							if newpos.Contain.(*Item).StepOn != nil {
								newpos.Contain.(*Item).StepOn(p)
							}
							newpos.Contain = nil
						}
						e.Forest[x][y].Contain = nil
						newpos.Contain = p
						p.X = x + dx
						p.Y = y + dy
						if newpos.StepOn != nil {
							newpos.StepOn(p)
						}
						// test kill
						e.checkAlive(i.ID, "Unknown")
					}
				case OPTurn:
					p.Turn = cur.Turn
				case OPShoot:
					if cur.Item == -1 { // use fist
						e.fight(p.ID, 1, 10, false)
					}
					if len(p.Bag) > cur.Item && cur.Item >= 0 && p.Bag[cur.Item].Damage > 0 {
						e.fight(p.ID, p.Bag[cur.Item].Vision, p.Bag[cur.Item].Damage, p.Bag[cur.Item].Penetrate)
						if p.Bag[cur.Item].Ammo > 0 { // currently not support more weapon...
							p.Bag = append(p.Bag[:cur.Item], p.Bag[cur.Item+1:]...)
						}
					}
				default:
				}
			}
		}
		e.OPQueue = []*OP{}
		e.LockOP.Unlock()

		// perform DVF
		if e.DVFPos != nil {
			e.DVFPos.Round++
			if e.DVFPos.Round%167 == 0 {
				e.DVFPos.PosX = Max(0, e.DVFPos.PosX-1)
				e.DVFPos.PosY = Max(0, e.DVFPos.PosY-1)
				e.DVFPos.Size += 2
				// World belongs to ZERO
				for m := e.DVFPos.PosX; m < Min(e.DVFPos.PosX+e.DVFPos.Size, e.Size); m++ {
					for n := e.DVFPos.PosY; n < Min(e.DVFPos.PosY+e.DVFPos.Size, e.Size); n++ {
						if e.Forest[m][n].GetType() == CPlayer {
							p := e.Forest[m][n].Contain.(*Player)
							p.KillNoDrop(p.Name + " killed by Dual Vector Foil")
							e.checkAlive(p.ID, "Unknown")
						}
						copier.CopyWithOption(&e.Forest[m][n], defaultGround["0"], copier.Option{DeepCopy: true})
					}
				}
			}
		}
	}
	for i, p := range e.Player {
		if p.HP > 0 {
			cnt := 0
			for _, j := range p.Bag {
				if j.Type == ITFlag {
					cnt++
				}
			}
			val, ok := user.Load(i)
			if ok {
				u := val.(*User)
				u.Score += cnt
				fmt.Printf("[LOG]%v score %v\n", i, u.Score)
				if err := rdb.Set(u.Token, fmt.Sprintf("%v\n%v\n%v", u.ID, u.Name, u.Score), 0).Err(); err != nil {
					fmt.Println("[ERROR]", err)
				}
			}
		}
	}
	os.Rename(fmt.Sprintf("round/%v.tmp", rnd), fmt.Sprintf("round/%v.txt", rnd))
	for i := range e.Player {
		ConnSend(i, "Game end")
	}
	time.Sleep(time.Second * 3)
}
