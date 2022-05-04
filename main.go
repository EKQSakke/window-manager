package main

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"golang.design/x/hotkey"
	"golang.org/x/sys/windows"
)

var (
	w32                     = windows.NewLazyDLL("user32.dll")
	procGetWindowText       = w32.NewProc("GetWindowTextW")
	procGetWindowTextLength = w32.NewProc("GetWindowTextLengthW")
	procGetForegroundWindow = w32.NewProc("GetForegroundWindow")
	procSetWindowPos        = w32.NewProc("SetWindowPos")
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
			test()
			println("I")
		}, hotkey.KeyI, mods))
	}()
	go func() {
		defer wg.Done()
		check(listenHotkey(func() {
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

func test() {
	if hwnd := getWindow(); hwnd != 0 {
		text := GetWindowText(HWND(hwnd))
		fmt.Println("window :", text, "# hwnd:", hwnd)

		setWindowPosition(hwnd, 1000, 500, 1000, 300)
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
