package appAction

import (
	"fmt"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

var (
	user32         = syscall.NewLazyDLL("user32.dll")
	procSetCursor  = user32.NewProc("SetCursorPos")
	procMouseEvent = user32.NewProc("mouse_event")
)

const (
	MOUSEEVENTF_LEFTDOWN = 0x0002
	MOUSEEVENTF_LEFTUP   = 0x0004
)

func moveMouse(x, y int) {
	// 鼠标移动到指定坐标
	procSetCursor.Call(uintptr(x), uintptr(y))
	fmt.Printf("鼠标移动到 (%d, %d)\n", x, y)
	time.Sleep(500 * time.Millisecond)
}

func leftClick(x, y int) {
	// 鼠标左键点击
	procMouseEvent.Call(MOUSEEVENTF_LEFTDOWN, uintptr(x), uintptr(y), 0, 0)
	procMouseEvent.Call(MOUSEEVENTF_LEFTUP, uintptr(x), uintptr(y), 0, 0)
	fmt.Println("鼠标左键点击")

}

func MouseClick(text string) {

	time.Sleep(3 * time.Second)

	moveMouse(550, 86)
	leftClick(550, 86)
	leftClick(550, 86)

	// 模拟键盘输入整段文字
	kb, _ := keybd_event.NewKeyBonding()

	clipboard.WriteAll(text)
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_V)
	kb.Launching()

	time.Sleep(1 * time.Second)
	kb.SetKeys(keybd_event.VK_ENTER)
	kb.Launching()

	time.Sleep(3 * time.Second)
	moveMouse(739, 383)
	leftClick(739, 383)

}

// 鼠标位置
//type POINT struct {
//	X, Y int32
//}
//
//func MouseClick() {
//	// 加载Windows系统库中的GetCursorPos函数
//	user32 := syscall.NewLazyDLL("user32.dll")
//	getCursorPos := user32.NewProc("GetCursorPos")
//
//	fmt.Println("开始监听鼠标位置（按Ctrl+C退出）...")
//	for {
//		// 定义一个POINT变量接收坐标
//		var pos POINT
//		// 调用GetCursorPos获取鼠标位置
//		ret, _, err := getCursorPos.Call(uintptr(unsafe.Pointer(&pos)))
//		if ret == 0 { // 返回值为0表示失败
//			fmt.Println("获取鼠标位置失败：", err)
//			return
//		}
//
//		// 打印当前坐标（X：水平坐标，Y：垂直坐标，屏幕左上角为(0,0)）
//		fmt.Printf("\r鼠标位置：X=%d, Y=%d", pos.X, pos.Y)
//
//		// 每隔100ms刷新一次（可调整频率）
//		time.Sleep(100 * time.Millisecond)
//	}
//}
