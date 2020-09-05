package execmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"zdora/types"
)

func generateCmd(args ...string) *exec.Cmd {
	var cmd *exec.Cmd
	args = append([]string{"-c"}, args...)
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command(os.Getenv("SHELL"), args...)
	case "windows":
		cmd = exec.Command("cmd", args...)
	default:
		cmd = exec.Command(os.Getenv("SHELL"), args...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	return cmd
}

func handleStdout(ctx context.Context, cancel context.CancelFunc, cmd *exec.Cmd, ch chan types.PID, logs chan string) {
	go func() {
		//创建获取命令输出管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
			return
		}
		//执行命令
		if err := cmd.Start(); err != nil {
			fmt.Println("Error:The command is err,", err)
			return
		}
		ch <- types.PID(cmd.Process.Pid)
		//使用带缓冲的读取器
		outputBuf := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				//一次获取一行,_ 获取当前行是否被读完
				output, _, err := outputBuf.ReadLine()
				if err != nil {
					// 判断是否到文件的结尾了否则出错
					if err != io.EOF {
						fmt.Printf("Error :%s\n", err)
					}
					cancel()
					return
				}
				logs <- string(output)
				fmt.Printf("%s\n", string(output))
			}
		}
	}()
}

func wait(cancel context.CancelFunc, cmd *exec.Cmd, logs chan string) {
	go func() {
		//wait 方法会一直阻塞到其所属的命令完全运行结束为止
		if err := cmd.Wait(); err != nil {
			cancel()
			fmt.Println("wait:", err.Error())
		}
	}()
}

func RunCmd(ctx context.Context, cancel context.CancelFunc, logs chan string, args ...string) types.PID {
	var pid types.PID
	cmd := generateCmd(args...)
	var pidChan = make(chan types.PID, 1)
	defer close(pidChan)
	handleStdout(ctx, cancel, cmd, pidChan, logs)
	select {
	case pid = <-pidChan:
	}
	wait(cancel, cmd, logs)
	return pid
}

func RunCmdCombined(args ...string) (types.PID, string) {
	cmd := generateCmd(args...)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return types.PID(cmd.Process.Pid), string(output)
	}
	return 0, ""
}
