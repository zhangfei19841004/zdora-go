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
	"zdora/constants"
	"zdora/model"
	"zdora/types"
	"zdora/util"
	"zdora/zdora/client"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute command at zdora's client",
	Args:  cobra.MinimumNArgs(1),
	Long: `在zdora客户端执行命令，如果指定了客户端ID，则在指定的客户端执行命令。
如果没有指定客户端ID，则会在空闲的客户端中随机选择并进行执行。
`,
	Run: func(cmd *cobra.Command, args []string) {
		msg := &model.Message{
			MessageType:    constants.ZDORA_EXEC,
			TargetClientId: types.ClientId(targetClientId),
			ExecMessage: &model.ExecutorMessage{
				Command:  args,
				Combined: combined,
			},
		}
		message, err := client.Cli(msg)
		if err != nil {
			fmt.Printf("error:%+v\n", err)
			return
		}
		var info model.Message
		util.Unmarshal(string(message), &info)
		fmt.Printf("command:%+v run at ip:[%s] clientId:[%s] pid:[%d]\n", info.ExecMessage.Command, info.ExecMessage.Ip, info.TargetClientId, info.ExecMessage.Pid)
	},
}

func init() {
	execCmd.Flags().StringVarP(&targetClientId, "target", "t", "", "execute command on taret client")
	execCmd.Flags().BoolVarP(&combined, "combine", "c", false, "execute command and combined output")
	rootCmd.AddCommand(execCmd)
}
