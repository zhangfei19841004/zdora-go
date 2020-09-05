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
	"text/tabwriter"
	"zdora/constants"
	"zdora/model"
	"zdora/util"
	"zdora/zdora/client"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show zdora's client infos",
	Long:  `List zdora's client infos`,
	Run: func(cmd *cobra.Command, args []string) {
		msg := &model.Message{
			MessageType: constants.PS,
		}
		message, err := client.Cli(msg)
		if err != nil {
			fmt.Printf("error:%+v\n", err)
			return
		}
		var infos []*model.PsCommandInfo
		util.Unmarshal(string(message), &infos)
		ps(infos)
	},
}

func ps(infos []*model.PsCommandInfo) {
	w := tabwriter.NewWriter(os.Stderr, 3, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", "ID", "ClientId", "IP", "PID", "Command", "Status")
	for i, info := range infos {
		status := "Running"
		if info.Pid == 0 {
			status = "Ready"
		}
		fmt.Fprintf(w, "%d\t%v\t%v\t%v\t%v\t%v\n", i+1, info.ClientId, info.Ip, info.Pid, info.Commond, status)
	}
	w.Flush()
}

func init() {
	rootCmd.AddCommand(psCmd)
}
