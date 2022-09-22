package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var randLetter = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = randLetter[rand.Intn(len(randLetter))]
	}
	return string(s)
}

func Eval(s string) (string, error) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	_, err := i.Eval(s)
	if err != nil {
		return "", err
	}
	ret, err := i.Eval("plugin()")
	return ret.String(), err
}

type User struct {
	Name     string
	Plugins  []*Plugin
	LogLevel logrus.Level
}

func (p User) String() string {
	return fmt.Sprintf("%v(%v): %v", p.Name, p.LogLevel, p.Plugins)
}

type Plugin struct {
	ID      int
	Enable  bool
	Name    string
	Version string
	Payload string
}

func (p Plugin) String() string {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	i.Eval(p.Version)
	r, _ := i.Eval("api()")
	return fmt.Sprintf("[%v] %v: %v", p.ID, p.Name, r.String())
}

var hackSatellite = &Plugin{
	ID:   1,
	Name: "one-tap hack satellite",
	Version: `
	import "io/ioutil"
	func api() string {
		v, _ := ioutil.ReadFile("satellite.txt")
		return string(v)
	}`,
	Payload: "satellite.go",
}

var hackSuperComputer = &Plugin{
	ID:   2,
	Name: "one-tap hack supercomputer",
	Version: `
	import "io/ioutil"
	func api() string {
		v, _ := ioutil.ReadFile("supercomputer.txt")
		return string(v)
	}`,
	Payload: "supercomputer.go",
}

var hackAllSubnet = &Plugin{
	ID:   3,
	Name: "one-tap hack all subnet computers",
	Version: `
	import "io/ioutil"
	func api() string {
		v, _ := ioutil.ReadFile("subnet.txt")
		return string(v)
	}`,
	Payload: "subnet.go",
}

var pluginList = []*Plugin{hackSatellite, hackSuperComputer, hackAllSubnet}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./static/", false)))
	r.GET("/login", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			c.AbortWithStatus(500)
			return
		} else {
			var user bytes.Buffer
			enc := gob.NewEncoder(&user)
			err := enc.Encode(User{Name: name, Plugins: pluginList, LogLevel: logrus.DebugLevel})
			if err != nil {
				logrus.Errorln(err)
				c.AbortWithStatus(500)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"config": base64.StdEncoding.EncodeToString(user.Bytes()),
			})
		}
	})
	r.GET("/list", func(c *gin.Context) {
		v, err := c.Request.Cookie("config")
		if err != nil || v.Value == "" {
			c.AbortWithStatus(500)
			return
		}
		cfg, err := base64.StdEncoding.DecodeString(v.Value)
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatus(500)
			return
		}
		logrus.Info(v.Value)
		dec := gob.NewDecoder(bytes.NewReader(cfg))
		var u User
		err = dec.Decode(&u)
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatus(500)
			return
		}
		if logrus.IsLevelEnabled(u.LogLevel) {
			logrus.Info(u)
		}
		kv := make(map[int]string)
		for _, v := range u.Plugins {
			kv[v.ID] = v.Name
		}
		c.JSON(http.StatusOK, gin.H{
			"plugin": kv,
		})
	})
	r.GET("/use", func(c *gin.Context) {
		id := c.Query("id")
		n, err := strconv.Atoi(id)
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatus(500)
			return
		}
		v, err := c.Request.Cookie("config")
		if err != nil || v.Value == "" {
			c.AbortWithStatus(500)
			return
		}
		cfg, err := base64.StdEncoding.DecodeString(v.Value)
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatus(500)
			return
		}
		logrus.Info(v.Value)
		dec := gob.NewDecoder(bytes.NewReader(cfg))
		var u User
		err = dec.Decode(&u)
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatus(500)
			return
		}
		for _, v := range u.Plugins {
			if v.ID == n {
				f, err := ioutil.ReadFile(v.Payload)
				if err != nil {
					logrus.Errorln(err)
					c.AbortWithStatus(500)
					return
				}
				ret, err := Eval(string(f))
				if err != nil {
					logrus.Errorln(err)
					c.AbortWithStatus(500)
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"data": ret,
				})
				return
			}
		}
	})
	r.GET("/less", func(c *gin.Context) {
		dst, err := base64.StdEncoding.DecodeString(c.Query("data"))
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatus(500)
			return
		}

		fname := RandString(8)
		err = ioutil.WriteFile("/tmp/"+fname, dst, 0755)
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatus(500)
			return
		}

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("lessc", "/tmp/"+fname)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Run()
		os.Remove("/tmp/" + fname)

		c.JSON(http.StatusOK, gin.H{
			"data": stdout.String(),
		})
	})
	r.Run()
}
