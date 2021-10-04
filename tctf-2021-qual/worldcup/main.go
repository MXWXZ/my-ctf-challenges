package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/unrolled/secure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const flag = "flag{fake_flag}"

const SESSION = "GOSESSIONID"

var store = sessions.NewCookieStore([]byte("10u_CanT_Gu1s_mE_b1cauSe_im_l0nG"))

var level1 = "NoQWeCy70QekDB5b"
var level2 = "Autx5F53FmmSFayM"
var f1, f2 []byte
var money = make(map[int]int)
var codeStore = make(map[int]string)

var db *gorm.DB
var r = gin.Default()

type Track struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID       int32  `gorm:"primaryKey;not null"`
	Username string `gorm:"uniqueIndex;type:varchar(32);not null"`
	Password string `gorm:"type:char(32);not null" json:"-"`
	Nickname string `gorm:"type:varchar(5);not null"`
	Msg      string `gorm:"type:varchar(128)"`
	Track    Track  `gorm:"embedded"`
}

func RandHex(n int) string {
	var letter = []byte("abcdef0123456789")
	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func needLogin(f func(c *gin.Context, s *sessions.Session)) func(c *gin.Context) {
	return func(c *gin.Context) {
		s, _ := store.Get(c.Request, SESSION)
		if v, ok := s.Values["login"]; !ok || v != 1 {
			c.AbortWithStatus(403)
			return
		}
		f(c, s)
	}
}

func needSession(f func(c *gin.Context, s *sessions.Session)) func(c *gin.Context) {
	return func(c *gin.Context) {
		s, _ := store.Get(c.Request, SESSION)
		f(c, s)
	}
}

func setting(c *gin.Context, s *sessions.Session) {
	c.HTML(http.StatusOK, "setting.tmpl", gin.H{
		"name":   s.Values["uname"],
		"nick":   s.Values["unick"],
		"_nonce": c.Keys["nonce"],
	})
}

func dashboard(c *gin.Context, s *sessions.Session) {
	lv := 1
	t1, err := c.Cookie("level1")
	if err == nil && t1 == level1 {
		lv += 1
	}
	code := RandHex(6)
	codeStore[s.Values["uid"].(int)] = code
	if lv == 1 {
		t, err := template.New(fmt.Sprintf("dashboard_%v", s.Values["uid"])).Parse(strings.ReplaceAll(string(f1), "[REPLACE]", s.Values["unick"].(string)))
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		t.Execute(c.Writer, gin.H{
			"name":   s.Values["uname"],
			"lv":     lv,
			"msg":    s.Values["msg"],
			"code":   code,
			"_nonce": c.Keys["nonce"],
		})
	} else {
		c.HTML(http.StatusOK, "dashboard2.tmpl", gin.H{
			"name":   s.Values["uname"],
			"lv":     lv,
			"msg":    s.Values["msg"],
			"code":   code,
			"_nonce": c.Keys["nonce"],
		})
	}
}

func apiFlag(c *gin.Context, s *sessions.Session) {
	lv := 1
	t1, err := c.Cookie("level1")
	if err == nil && t1 == level1 {
		lv += 1
	}
	t2, err := c.Cookie("level2")
	if err == nil && t2 == level2 {
		lv += 1
	}
	if lv != 3 {
		c.String(403, "You need admin to view this")
		return
	}

	if _, ok := money[s.Values["uid"].(int)]; !ok {
		money[s.Values["uid"].(int)] = 0
	}

	m := money[s.Values["uid"].(int)]
	if m >= 200 {
		money[s.Values["uid"].(int)] = m - 200
		c.JSON(200, gin.H{"code": 0, "msg": flag})
	} else {
		c.JSON(200, gin.H{"code": 1, "msg": "Poor man!"})
	}
}

func casino(c *gin.Context, s *sessions.Session) {
	lv := 1
	t1, err := c.Cookie("level1")
	if err == nil && t1 == level1 {
		lv += 1
	}
	t2, err := c.Cookie("level2")
	if err == nil && t2 == level2 {
		lv += 1
	}
	if lv != 3 {
		c.String(403, "You need two level admin secrets to view this, get at dashboard")
		return
	}

	rep := c.DefaultQuery("bet", "1")
	exp := regexp.MustCompile(`\.[^a-zA-Z0-9]`)
	if exp.Find([]byte(rep)) != nil {
		c.String(400, "Param error")
		return
	}

	if _, ok := money[s.Values["uid"].(int)]; !ok {
		money[s.Values["uid"].(int)] = 0
	}

	r1 := rand.Intn(100)
	r2 := rand.Intn(100)
	m := money[s.Values["uid"].(int)]
	var prepare int
	if r1 >= r2 {
		if string(rep[0]) == "1" {
			prepare = m + 10
		} else {
			prepare = 0
		}
	} else {
		if string(rep[0]) == "0" {
			prepare = m + 10
		} else {
			prepare = 0
		}
	}
	t, err := template.New(fmt.Sprintf("casino_%v", s.Values["uid"])).Parse(strings.ReplaceAll(string(f2), "[REPLACE]", rep))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	err = t.Execute(c.Writer, gin.H{
		"o0ps_u_Do1nt_no_t1": r1,
		"o0ps_u_Do1nt_no_t2": r2,
		"money":              prepare,
		"_nonce":             c.Keys["nonce"],
	})
	if err != nil {
		return
	}
	money[s.Values["uid"].(int)] = prepare
}

func login(c *gin.Context, s *sessions.Session) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"_nonce": c.Keys["nonce"],
	})
}

func signup(c *gin.Context, s *sessions.Session) {
	c.HTML(http.StatusOK, "signup.tmpl", gin.H{
		"_nonce": c.Keys["nonce"],
	})
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type checkParam struct {
	Code string `form:"code" binding:"required"`
}

func apiCheck(c *gin.Context, s *sessions.Session) {
	var param checkParam
	err := c.ShouldBind(&param)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	lv := 1
	t1, err := c.Cookie("level1")
	if err == nil && t1 == level1 {
		lv += 1
	}

	if _, ok := codeStore[s.Values["uid"].(int)]; !ok {
		c.AbortWithStatus(400)
		return
	}
	if MD5(param.Code)[0:6] != codeStore[s.Values["uid"].(int)] {
		c.JSON(200, gin.H{"code": 1, "msg": "Invalid code"})
		return
	}

	delete(codeStore, s.Values["uid"].(int))
	go xssrun(s.Values["uid"].(int), lv)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  fmt.Sprintf("admin level %v will check this room soon", lv),
	})
}

type chatParam struct {
	Msg string `form:"msg" binding:"required,max=128"`
}

func apiChat(c *gin.Context, s *sessions.Session) {
	var param chatParam
	err := c.ShouldBind(&param)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	var rec User
	err = db.First(&rec, s.Values["uid"]).Error
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	rec.Msg = param.Msg
	err = db.Save(&rec).Error
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	s.Values["msg"] = rec.Msg
	s.Save(c.Request, c.Writer)
	c.JSON(200, gin.H{"code": 0, "msg": "OK"})
}

type settingParam struct {
	Nickname string `form:"nick" binding:"required,max=5"`
	Password string `form:"password"`
}

func apiSetting(c *gin.Context, s *sessions.Session) {
	var param settingParam
	err := c.ShouldBind(&param)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	var rec User
	err = db.First(&rec, s.Values["uid"]).Error
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	rec.Nickname = param.Nickname
	if param.Password != "" {
		rec.Password = MD5(param.Password)
	}
	err = db.Save(&rec).Error
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	s.Values["unick"] = rec.Nickname
	s.Save(c.Request, c.Writer)
	c.JSON(200, gin.H{"code": 0, "msg": "OK"})
}

type signParam struct {
	Username string `form:"username" binding:"required,max=32"`
	Password string `form:"password" binding:"required"`
}

func apiSignin(c *gin.Context, s *sessions.Session) {
	var param signParam
	err := c.ShouldBind(&param)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	var rec User
	err = db.First(&rec, "username = ?", param.Username).Error
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": "Invalid username/password"})
		return
	}
	if rec.Password != MD5(param.Password) {
		c.JSON(200, gin.H{"code": 1, "msg": "Invalid username/password"})
		return
	}

	money[int(rec.ID)] = 0
	s.Values["uid"] = int(rec.ID)
	s.Values["uname"] = rec.Username
	s.Values["unick"] = rec.Nickname
	s.Values["login"] = 1
	s.Values["msg"] = rec.Msg
	s.Save(c.Request, c.Writer)
	c.JSON(200, gin.H{"code": 0, "msg": "OK"})
}

func apiSignup(c *gin.Context, s *sessions.Session) {
	var param signParam
	err := c.ShouldBind(&param)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	var rec User
	if len(param.Username) > 5 {
		rec = User{
			Username: param.Username,
			Password: MD5(param.Password),
			Nickname: param.Username[0:5],
		}
	} else {
		rec = User{
			Username: param.Username,
			Password: MD5(param.Password),
			Nickname: param.Username,
		}
	}

	err = db.Create(&rec).Error
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": "Exist username"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "OK"})
}

type botParam struct {
	ID int `uri:"id" binding:"required"`
	LV int `uri:"lv" binding:"required"`
}

func bot(c *gin.Context) {
	var param botParam
	err := c.ShouldBindUri(&param)
	if err != nil {
		c.String(404, "404 page not found")
		return
	}

	var rec User
	err = db.First(&rec, "id = ?", param.ID).Error
	if err != nil {
		c.String(404, "404 page not found")
		return
	}
	if param.LV == 1 {
		t, err := template.New(fmt.Sprintf("dashboard_%v", param.ID)).Parse(strings.ReplaceAll(string(f1), "[REPLACE]", rec.Nickname))
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		t.Execute(c.Writer, gin.H{
			"name":   rec.Username,
			"lv":     param.LV,
			"msg":    rec.Msg,
			"code":   "",
			"_nonce": c.Keys["nonce"],
		})
	} else {
		c.HTML(http.StatusOK, "dashboard2.tmpl", gin.H{
			"name":   rec.Username,
			"lv":     param.LV,
			"msg":    rec.Msg,
			"code":   "",
			"_nonce": c.Keys["nonce"],
		})
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 12,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	var err error
	db, err = gorm.Open(sqlite.Open("worldcup.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		ContentSecurityPolicy: "default-src 'none'; script-src $NONCE; connect-src 'self'; img-src 'self' data:; style-src 'self'; base-uri 'self'; form-action 'self'; font-src 'self'",
		ReferrerPolicy:        "same-origin",
		IsDevelopment:         false,
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			nonce, err := secureMiddleware.ProcessAndReturnNonce(c.Writer, c.Request)
			if err != nil {
				c.Abort()
				return
			}

			c.Set("nonce", nonce)

			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()
	r.Use(secureFunc)

	f1, err = ioutil.ReadFile("templates/dashboard.tmpl")
	if err != nil {
		panic(err)
	}
	f2, err = ioutil.ReadFile("templates/casino.tmpl")
	if err != nil {
		panic(err)
	}

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.GET("/", needSession(login))
	r.GET("/signup", needSession(signup))
	r.POST("/api/signin", needSession(apiSignin))
	r.POST("/api/signup", needSession(apiSignup))

	r.GET("/dashboard", needLogin(dashboard))
	r.GET("/setting", needLogin(setting))
	r.GET("/casino", needLogin(casino))
	r.POST("/api/setting", needLogin(apiSetting))
	r.POST("/api/chat", needLogin(apiChat))
	r.POST("/api/check", needLogin(apiCheck))
	r.POST("/api/flag", needLogin(apiFlag))

	r.GET("/wTf_Pa1h_5604m/:id/:lv", bot)
	r.Run()
}
