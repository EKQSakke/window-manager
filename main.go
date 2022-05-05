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

	windowList = []uintptr{}
	screenWidth = int(win.GetSystemMetrics(win.SM_CXSCREEN))
    
	screenHeight = int(win.GetSystemMetrics(win.SM_CYSCREEN))
	multi = .97
	relativeScreenHeight = int(multi * float64(screenHeight))
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	wg := sync.WaitGroup{}
	wg.Add(2)
	mods := []hotkey.Modifier{hotkey.ModShift, hotkey.ModCtrl}
	go func() {
		defer wg.Done()
		check(listenHotkey(func() {
			addCurrentWindowToList()
			println("I")
		}, hotkey.KeyI, mods))
	}()
	go func() {
		defer wg.Done()
		check(listenHotkey(func() {
			printAllWindows()
			positionAllWindows()
			println("K")
		}, hotkey.KeyK, mods))
	}()
	wg.Wait()
}

func listenHotkey(onKeyDown func(), key hotkey.Key, mods []hotkey.Modifier) (err error) {
	ms := []hotkey.Modifier{}
	ms = append(ms, mods...)
	hk := hotkey.New(ms, key)

	check(hk.Register())

	for {
		// Blocks until the hotkey is triggered.
		<-hk.Keydown()
		onKeyDown()
	}
}

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

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
	procSetWindowPos.Call(hwnd, 0, uintptr(x), uintptr(y), uintptr(width), uintptr(height), 0)
}

func positionAllWindows() {	
	for i, hwnd := range windowList {
		if i == 0 { 
			setWindowPosition(hwnd, 0, 0, screenWidth / 2, relativeScreenHeight)
		} else {
			setWindowPosition(hwnd, screenWidth / 2, 0, screenWidth / 2, relativeScreenHeight)
		}
	}
}

func printAllWindows() {
	for _, hwnd := range windowList {
		text := GetWindowText(HWND(hwnd))
		fmt.Println(text, "# hwnd:", hwnd)
	}
}

func addCurrentWindowToList() {
	if hwnd := getWindow(); hwnd != 0 {
		if !Contains(windowList, hwnd) {
			windowList = append(windowList, hwnd)
			println("window already in list")
			return
		}

		text := GetWindowText(HWND(hwnd))
		fmt.Println("added window :", text)
		windowList = append(windowList, hwnd)
	}	
}