package main

// operation
type OPType int

const (
	OPWalk OPType = iota
	OPTurn
	OPShoot
)

type TurnType int

const (
	TUp TurnType = iota
	TRight
	TDown
	TLeft
)

type OP struct {
	ID   int
	Item int
	Oper OPType
	Turn TurnType
	Next *OP
}

func (o *OP) NextTurn(t TurnType) *OP {
	o.Next = &OP{
		Oper: OPTurn,
		Turn: t,
	}
	return o.Next
}

func (o *OP) NextWalk() *OP {
	o.Next = &OP{
		Oper: OPWalk,
	}
	return o.Next
}

func (o *OP) NextShoot(i int) *OP {
	o.Next = &OP{
		Oper: OPShoot,
		Item: i,
	}
	return o.Next
}

type CellType int

const (
	CGround CellType = iota
	CPlayer
	CItem
)

type Cell interface {
	GetName(ob bool) string
	GetType() CellType
}

// Item

type ItemType int

const (
	ITFlag ItemType = iota
	ITUSP
	ITAK47
	ITAWP
	ITArmer     // 甲
	ITBlood     // 血包
	ITMine      // 地雷
	ITSophon    // 智子
	ITWaterDrop // 水滴
	ITDVF       // 二向箔
	ITPack      // 包
	//////
	ITTot
)

type DVF struct {
	PosX  int
	PosY  int
	Round int
	Size  int
}

type Item struct {
	Name      string
	Damage    int  // 伤害
	Vision    int  // 射程
	Penetrate bool // 多目标伤害
	Pack      []*Item
	Ammo      int
	Invisible bool
	Type      ItemType
	StepOn    StepOnFunc
}

func (e *Item) GetName(ob bool) string {
	if e.Invisible && !ob {
		return ""
	}
	return e.Name
}

func (e *Item) GetType() CellType {
	return CItem
}

var defaultItem map[int]*Item

func init() {
	wrapper := func(n ItemType) StepOnFunc {
		return func(p *Player) {
			p.PickUp(defaultItem[int(n)])
		}
	}
	defaultItem = make(map[int]*Item)
	defaultItem[int(ITFlag)] = &Item{
		Name:   "F",
		Type:   ITFlag,
		StepOn: wrapper(ITFlag),
	}
	defaultItem[int(ITUSP)] = &Item{
		Name:   "U",
		Damage: 20,
		Vision: 6,
		Type:   ITUSP,
		StepOn: wrapper(ITUSP),
	}
	defaultItem[int(ITAK47)] = &Item{
		Name:   "K",
		Damage: 40,
		Vision: 8,
		Type:   ITAK47,
		StepOn: wrapper(ITAK47),
	}
	defaultItem[int(ITAWP)] = &Item{
		Name:      "A",
		Damage:    80,
		Vision:    12,
		Penetrate: true,
		Type:      ITAWP,
		StepOn:    wrapper(ITAWP),
	}
	defaultItem[int(ITArmer)] = &Item{
		Name: "R",
		Type: ITArmer,
		StepOn: func(p *Player) {
			p.PickUp(defaultItem[int(ITArmer)])
			p.Armer += 10
		},
	}
	defaultItem[int(ITBlood)] = &Item{
		Name: "B",
		Type: ITBlood,
		StepOn: func(p *Player) {
			p.HP += 20
		},
	}
	defaultItem[int(ITMine)] = &Item{
		Name:      "M",
		Type:      ITMine,
		Invisible: true,
		StepOn: func(p *Player) {
			p.Kill(p.Name + " killed by Mine")
		},
	}
	defaultItem[int(ITSophon)] = &Item{
		Name: "S",
		Type: ITSophon,
		StepOn: func(p *Player) {
			p.PickUp(defaultItem[int(ITSophon)])
			p.Vision = 999
		},
	}
	defaultItem[int(ITWaterDrop)] = &Item{
		Name:      "D",
		Damage:    999,
		Vision:    999,
		Penetrate: true,
		Ammo:      1,
		Type:      ITWaterDrop,
		StepOn:    wrapper(ITWaterDrop),
	}
	defaultItem[int(ITDVF)] = &Item{
		Name: "V",
		Type: ITDVF,
		StepOn: func(p *Player) {
			// start vector attack
			engine.DVFPos = &DVF{
				PosX:  p.X,
				PosY:  p.Y,
				Round: 0,
				Size:  1,
			}
			engine.MsgQueue = append(engine.MsgQueue, "WARNING! "+p.Name+" throws an Dual Vector Foil!!!")
		},
	}
}

// Player

type Player struct {
	ID     int
	Name   string
	HP     int
	Armer  int
	Killed int
	Vision int
	X      int
	Y      int
	Death  string
	Turn   TurnType
	Bag    []*Item
}

func (p *Player) Kill(r string) {
	p.Death = r
	p.HP = 0
}

func (p *Player) KillNoDrop(r string) {
	p.Death = r
	p.HP = 0
	p.Bag = []*Item{}
}

func (p *Player) GetName(ob bool) string {
	return p.Name
}

func (p *Player) GetType() CellType {
	return CPlayer
}

func (p *Player) PickUp(e *Item) {
	p.Bag = append(p.Bag, e)
}

// Ground

type GroundType int

const (
	GNone GroundType = iota
	GGround
	GTree
	GWater
)

type StepOnFunc func(*Player)

type Ground struct {
	Name     string
	Contain  Cell
	AllowSee bool
	AllowRun bool
	Type     GroundType
	StepOn   StepOnFunc
}

func (c *Ground) GetType() CellType {
	if c.Contain != nil {
		return c.Contain.GetType()
	}
	return CGround
}

func (c *Ground) GetName(ob bool) string {
	if c.Contain != nil {
		name := c.Contain.GetName(ob)
		if name == "" {
			return c.Name
		}
		return name
	}
	return c.Name
}

var defaultGround = map[string]*Ground{
	// space
	"0": {
		Name:     "0",
		AllowSee: false,
		AllowRun: true,
		Type:     GNone,
		StepOn: func(p *Player) {
			p.KillNoDrop(p.Name + " suicide")
		},
	},
	// ground
	"G": {
		Name:     "G",
		AllowSee: true,
		AllowRun: true,
		Type:     GGround,
	},
	// Tree
	"T": {
		Name:     "T",
		AllowSee: false,
		AllowRun: false,
		Type:     GTree,
	},
	// Water
	"W": {
		Name:     "W",
		AllowSee: true,
		AllowRun: true,
		Type:     GWater,
		StepOn: func(p *Player) {
			p.KillNoDrop(p.Name + " suicide")
		},
	},
}

// helper function
func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
