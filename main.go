package main

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"golang.design/x/hotkey"
	"golang.org/x/sys/windows"
)

var (
	w32                     = windows.NewLazyDLL("user32.dll")
	procGetWindowText       = w32.NewProc("GetWindowTextW")
	procGetWindowTextLength = w32.NewProc("GetWindowTextLengthW")
	procGetForegroundWindow = w32.NewProc("GetForegroundWindow")
	procSetWindowPos        = w32.NewProc("SetWindowPos")

	windowList           = []uintptr{}
	screenWidth          = float32(win.GetSystemMetrics(win.SM_CXSCREEN))
	screenHeight         = float32(win.GetSystemMetrics(win.SM_CYSCREEN))
	multi                = .97
	relativeScreenHeight = float32(multi * float64(screenHeight))
)

func main() {
	ShowNotification("Hello!")

	wg := sync.WaitGroup{}
	wg.Add(5)
	mods := []hotkey.Modifier{hotkey.ModAlt, hotkey.ModCtrl}
	go func() {
		defer wg.Done()
		Check(listenHotkey(func() {
			addCurrentWindowToList()
		}, hotkey.KeyI, mods))
	}()
	go func() {
		defer wg.Done()
		Check(listenHotkey(func() {
			printAllWindows()
			positionAllWindows()
		}, hotkey.KeyK, mods))
	}()

	go func() {
		defer wg.Done()
		Check(listenHotkey(func() {
			moveWindow(-1)
		}, hotkey.KeyJ, mods))
	}()

	go func() {
		defer wg.Done()
		Check(listenHotkey(func() {
			moveWindow(1)
		}, hotkey.KeyL, mods))
	}()

	wg.Wait()
}

func moveWindow(i int) {
	if hwnd := getWindow(); hwnd != 0 {
		currentWindowId := GetId(windowList, hwnd)
		if currentWindowId == -1 {
			return
		}

		targetWindowId := currentWindowId + i

		if targetWindowId == -1 {
			targetWindowId = len(windowList) - 1
		} else if targetWindowId == len(windowList) {
			targetWindowId = 0
		}

		temp := windowList[targetWindowId]
		windowList[targetWindowId] = windowList[currentWindowId]
		windowList[currentWindowId] = temp

		positionAllWindows()
	}
}

func listenHotkey(onKeyDown func(), key hotkey.Key, mods []hotkey.Modifier) (err error) {
	ms := []hotkey.Modifier{}
	ms = append(ms, mods...)
	hk := hotkey.New(ms, key)
	Check(hk.Register())

	for {
		// Blocks until the hotkey is triggered.
		<-hk.Keydown()
		onKeyDown()
	}
}

func getWindowTextLength(hwnd uintptr) int {
	ret, _, _ := procGetWindowTextLength.Call(uintptr(hwnd))

	return int(ret)
}

func getWindowText(hwnd uintptr) string {
	textLen := getWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func getWindow() uintptr {
	hwnd, _, _ := procGetForegroundWindow.Call()
	return hwnd
}

func setWindowPosition(hwnd uintptr, x, y, width, height int) {
	procSetWindowPos.Call(hwnd,
		0,
		uintptr(float32(x)/100*screenWidth),
		uintptr(float32(y)/100*relativeScreenHeight),
		uintptr(float32(width)/100*screenWidth),
		uintptr(float32(height)/100*relativeScreenHeight),
		0)
}

func positionAllWindows() {
	if len(windowList) <= 1 {
		return
	}

	layout := getLayout(len(windowList) - 1)
	for i, hwnd := range windowList {
		setWindowPosition(hwnd, layout.Windows[i].x, layout.Windows[i].y, layout.Windows[i].w, layout.Windows[i].h)
	}
}

func printAllWindows() {
	for _, hwnd := range windowList {
		text := getWindowText(uintptr(hwnd))
		fmt.Println(text, "# hwnd:", hwnd)
	}
}

func addCurrentWindowToList() {
	if hwnd := getWindow(); hwnd != 0 {
		if !Contains(windowList, hwnd) {
			windowList = append(windowList, hwnd)
			text := getWindowText(uintptr(hwnd))
			ShowNotification(fmt.Sprintf("%v added to list", text))

			getLayout(len(windowList))
			return
		}

		println("window already in list")
	}
}
