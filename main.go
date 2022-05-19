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
	procEnumWindows         = w32.NewProc("EnumWindows")
	procIsWindowVisible     = w32.NewProc("IsWindowVisible")
	procIsIconic            = w32.NewProc("IsIconic")
	procGetClientRect       = w32.NewProc("GetClientRect")

	windowList           = []uintptr{}
	screenWidth          = float32(win.GetSystemMetrics(win.SM_CXSCREEN))
	screenHeight         = float32(win.GetSystemMetrics(win.SM_CYSCREEN))
	multi                = .97
	relativeScreenHeight = float32(multi * float64(screenHeight))
)

type LONG int32

type RECT struct {
	left   LONG
	top    LONG
	right  LONG
	bottom LONG
}

func main() {
	ShowNotification("Hello!")
	getAllWindows()

	wg := sync.WaitGroup{}
	wg.Add(2)
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

func getAllWindows() {
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		text := getWindowText(uintptr(h))

		rect := RECT{}

		r1, _, _ := procIsWindowVisible.Call(uintptr(h))
		procGetClientRect.Call(uintptr(h), uintptr(unsafe.Pointer(&rect)))

		if r1 == 1 && isValidRect(rect) && isValidTitle(text) {
			println(text)
		}
		return 1
	})
	enumWindows(cb, 0)
}

func isValidTitle(text string) bool {
	if text == "" {
		return false
	}

	// These seem to always exist
	if text == "Settings" || text ==  "Movies & TV" || text == "MainWindow" {
		return false
	}

	return true
}

func isValidRect(rect RECT) bool {
	if rect.bottom == 0 && rect.right == 0 {
		return false
	}

	if rect.right == LONG(screenWidth) && rect.bottom == LONG(screenHeight) {
		return false
	}

	return true
}

func enumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
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
