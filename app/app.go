package app

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"pmon/config"
	"pmon/db"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kardianos/service"
)

// Program ...
type Program struct {
	IsRunning bool
	Pname     string
	RedisHost string
	RedisKeys map[string]string
	RedisConn redis.Conn
}

// Start ...
func (p *Program) Start(s service.Service) error {
	config.Init()
	db.Init()

	p.Pname = config.GetConfig().Target

	p.RedisKeys = config.GetConfig().DB.Redis.Keys

	p.IsRunning = true
	go p.run()
	return nil
}

// Stop ...
func (p *Program) Stop(s service.Service) error {
	p.RedisConn.Do("SET", p.RedisKeys["status"], 0)
	p.RedisConn.Close()

	p.IsRunning = false

	return nil
}

func (p *Program) run() {
	var wg sync.WaitGroup

	statusCh := make(chan int)
	qlenCh := make(chan int)

	wg.Add(3)

	go func() {
		redisNeedReconnect := make(chan int)

	touchRedis:
		conn := db.Redis()
		if conn == nil {
			time.Sleep(3 * time.Second)
			goto touchRedis
		}

		defer func() {
			conn.Close()
			wg.Done()
		}()

		for p.IsRunning {
			select {
			case <-redisNeedReconnect:
				conn = db.Redis()
				if conn == nil {
					time.Sleep(3 * time.Second)
					redisNeedReconnect <- 1
					break
				}
				p.RedisConn = conn
			case v := <-statusCh:
				_, err := conn.Do("SET", p.RedisKeys["status"], v)

				if err != nil {
					redisNeedReconnect <- 1
				}
			case v := <-qlenCh:
				_, err := conn.Do("SET", p.RedisKeys["qlen"], v)
				if err != nil {
					redisNeedReconnect <- 1
				}
			}
		}
	}()

	go func() {
		defer wg.Done()

		for p.IsRunning {
			qlen := -1
			db.Mssql("store").Table("U9.ProcingVoucher").Count(&qlen)
			if qlen != -1 {
				qlenCh <- qlen
			}
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		defer wg.Done()

		for p.IsRunning {
			pid, err := findProcessID(p.Pname)
			if err == nil {
				process, err := os.FindProcess(pid)
				if err == nil {
					statusCh <- 1
					process.Wait()
				}
			}

			statusCh <- 0

			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

func findProcessID(processName string) (int, error) {
	buf := bytes.Buffer{}
	cmd := exec.Command("wmic", "process", "get", "processid,name")
	cmd.Stdout = &buf
	cmd.Run()

	cmd2 := exec.Command("findstr", processName)
	cmd2.Stdin = &buf
	data, _ := cmd2.CombinedOutput()

	if len(data) == 0 {
		return -1, errors.New("not found")
	}

	info := string(data)
	reg := regexp.MustCompile(`[0-9]+`)
	pid := reg.FindString(info)

	return strconv.Atoi(pid)
}
