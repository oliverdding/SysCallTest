/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/alitto/pond"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	syscalltest "syscalltest/internal"
)

var rootCmd = &cobra.Command{
	Use:   "syscalltest",
	Short: "A demo program for testing jaeger-client-go with HostIP function",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		sugar := logger.Sugar()

		sugar.Infoln("bench start", "task", task, "thread", thread, "interval", interval)

		if enablePprof {
			sugar.Infoln("start pprof", "port", port)
			go func() {
				sugar.Infoln(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
			}()
		}

		pool := pond.New(int(thread), int(thread)*10)

		for i := uint(0); i < task; i++ {
			n := i
			pool.Submit(func() {
				sugar.Infoln("task start", "id", n)
				time.Sleep(time.Microsecond * time.Duration(interval))
				ip, err := syscalltest.HostIP()
				if err != nil {
					sugar.Errorln("cannot get host ip", "id", n, "error", err)
				} else {
					sugar.Infoln("get host ip", "id", n, "ip", ip)
				}
				sugar.Infoln("task end", "id", n)
			})
		}

		pool.StopAndWait()

		sugar.Infoln("bench stop")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	enablePprof bool
	port        uint
	task        uint
	thread      uint
	interval    uint
)

func init() {
	rootCmd.Flags().BoolVarP(&enablePprof, "pprof", "P", true, "enable pprof")
	rootCmd.Flags().UintVarP(&port, "port", "p", 6060, "port for pprof")
	rootCmd.Flags().UintVarP(&task, "count", "c", 10000, "amount of task")
	rootCmd.Flags().UintVarP(&thread, "thread", "t", 1000, "amount of thread")
	rootCmd.Flags().UintVarP(&interval, "interval", "i", 100, "milisecond between two syscall")
}
