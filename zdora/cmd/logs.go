/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"zdora/constants"
	"zdora/model"
	"zdora/types"
	"zdora/zdora/client"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show zdora's client log",
	Long:  `展示客户端运行的log信息。`,
	Run: func(cmd *cobra.Command, args []string) {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		msg := &model.Message{
			MessageType:    constants.LOGS,
			TargetClientId: types.ClientId(targetClientId),
			ExecMessage: &model.ExecutorMessage{
				Pid: types.PID(pid),
			},
		}
		err := client.CliWithInterrupt(msg, interrupt)
		if err != nil {
			fmt.Printf("error:%+v\n", err)
			return
		}
	},
}

func init() {
	logsCmd.Flags().StringVarP(&targetClientId, "target", "t", "", "execute command on taret client")
	logsCmd.Flags().IntVarP(&pid, "pid", "p", 0, "execute command and combined output")
	rootCmd.AddCommand(logsCmd)
}
