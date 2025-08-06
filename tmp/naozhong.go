package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
)

var (
	activeAlarm  *Alarm
	quitChan     = make(chan struct{})
	systrayReady = make(chan struct{})
)

// Alarm 表示一个闹钟提醒
type Alarm struct {
	triggerTime time.Time
	notifyChan  chan struct{}
	stopChan    chan struct{}
}

func main() {
	// 启动系统托盘
	go startSystray()
	<-systrayReady // 等待系统托盘初始化完成

	fmt.Println("闹钟程序已启动! 按Ctrl+C退出或在系统托盘选择退出")

	// 设置中断信号捕获
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// 首次提醒
	go triggerAlarm()

	// 创建定时器（每1小时）
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// 主循环
loop:
	for {
		select {
		case <-ticker.C:
			go triggerAlarm()
		case <-interrupt:
			fmt.Println("\n程序已退出")
			systray.Quit()
			break loop
		case <-quitChan:
			fmt.Println("\n用户退出程序")
			systray.Quit()
			break loop
		}
	}
}

// triggerAlarm 触发闹钟提醒
func triggerAlarm() {
	if activeAlarm != nil {
		return // 如果已有激活的提醒，则跳过
	}

	alarm := &Alarm{
		triggerTime: time.Now(),
		notifyChan:  make(chan struct{}, 1),
		stopChan:    make(chan struct{}),
	}
	activeAlarm = alarm

	fmt.Printf("\n提醒时间: %s\n", alarm.triggerTime.Format("15:04:05"))

	go playNotificationSound(alarm.stopChan)
	go showNotification(alarm.notifyChan)

	// 等待用户关闭或手动停止
	select {
	case <-alarm.notifyChan:
		fmt.Println("用户点击关闭了通知")
	case <-time.After(5 * time.Minute):
		fmt.Println("提醒超时自动停止")
	}

	close(alarm.stopChan) // 停止声音播放
	activeAlarm = nil
}

// playNotificationSound 播放通知声音
func playNotificationSound(stopChan chan struct{}) {
	if !isAudioEnabled() {
		fmt.Println("系统音频已禁用，不播放提示音")
		return
	}

	// 尝试Windows系统的Beep API
	dll := syscall.NewLazyDLL("kernel32.dll")
	proc := dll.NewProc("Beep")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 播放提示音 (523Hz, 持续200ms)
			_, _, _ = proc.Call(523, 200)
		case <-stopChan:
			return
		}
	}
}

// showNotification 显示通知
func showNotification(notifyChan chan struct{}) {
	// 发送带提示音的桌面通知
	err := beeep.Notify("闹钟提醒", "已过去1小时！点击关闭本次提醒。", "")
	if err != nil {
		fmt.Printf("无法发送通知: %v\n", err)
		return
	}

	// 等待用户关闭通知窗口
	go func() {
		// 对于非阻塞的通知，我们模拟通知结束
		// 实际应用中，可以使用更高级的通知库
		fmt.Println("通知已显示，等待用户操作...")
		time.Sleep(1 * time.Second)
		notifyChan <- struct{}{}
	}()
}

// startSystray 创建系统托盘图标
func startSystray() {
	// 正确的systray.Run调用方式
	systray.Run(
		// onReady函数
		func() {
			systray.SetIcon(getIcon())
			systray.SetTitle("闹钟程序")
			systray.SetTooltip("每小时提醒程序")

			// 添加菜单项
			mNow := systray.AddMenuItem("立即提醒", "立即触发一次提醒")
			systray.AddSeparator()
			mStop := systray.AddMenuItem("停止当前提醒", "停止正在播放的提示音")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("退出程序", "完全退出程序")

			close(systrayReady) // 通知主程序系统托盘已准备好

			for {
				select {
				case <-mNow.ClickedCh:
					go triggerAlarm()
				case <-mStop.ClickedCh:
					if activeAlarm != nil {
						fmt.Println("用户手动停止当前提醒")
						activeAlarm.stopChan <- struct{}{}
					}
				case <-mQuit.ClickedCh:
					close(quitChan)
					return
				}
			}
		},
		// onExit函数（清理资源）
		func() {
			fmt.Println("系统托盘已退出")
		},
	)
}

// getIcon 获取系统托盘图标 (以字节数组形式嵌入)
func getIcon() []byte {
	// 这是一个占位符图标，实际使用时应替换为真实图标
	return []byte{}
}

// isAudioEnabled 检查系统是否启用了音频
func isAudioEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Sound`, registry.READ)
	if err != nil {
		return true // 如果无法检查，默认为启用
	}
	defer k.Close()

	s, _, err := k.GetStringValue("Beep")
	if err != nil || s != "no" {
		return true
	}
	return false
}
