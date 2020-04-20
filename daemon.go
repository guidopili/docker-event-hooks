package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var BackgroundDaemonInfoDir = ".docker-event-hooks"
var Timestamp = time.Now().Unix()

var PIDFile = fmt.Sprintf("%s/%d.pid", BackgroundDaemonInfoDir, Timestamp)
var LogFile = fmt.Sprintf("%s/%d.log", BackgroundDaemonInfoDir, Timestamp)

func savePID(pid int) {
	file, err := os.Create(PIDFile)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	file.Sync()
}

func startBg(cmd *exec.Cmd) {
	_, err := os.Stat(BackgroundDaemonInfoDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(BackgroundDaemonInfoDir, 0777)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	filePtr, _ := os.Create(LogFile)
	cmd.Stdout = filePtr
	cmd.Stderr = filePtr

	_ = cmd.Start()
	savePID(cmd.Process.Pid)

	go func() {
		_ = <-ch
		signal.Stop(ch)
		os.Remove(PIDFile)
		os.Exit(0)
	}()

	os.Exit(0)
}

func stopBg() {
	_ = filepath.Walk(BackgroundDaemonInfoDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".pid" {
			dat, err := ioutil.ReadFile(path)

			if err != nil {
				return err
			}

			pid, err := strconv.Atoi(string(dat))

			if err != nil {
				return err
			}

			log.Println(fmt.Sprintf("Sending SIGTERM to process with PID %v", pid))

			syscall.Kill(pid, syscall.SIGTERM)
			os.Remove(path)
		}

		return nil
	})

	os.Exit(0)
}
