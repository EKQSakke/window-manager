package main

import (
	"sync"

	"golang.design/x/hotkey"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	mods := []hotkey.Modifier{ hotkey.ModShift, hotkey.ModCtrl }
	go func() {
		defer wg.Done()

		check(listenHotkey(func(){
			println("I")
		}, hotkey.KeyI, mods))
	}()
	go func() {
		defer wg.Done()

		check(listenHotkey(func(){
			println("K")
		}, hotkey.KeyK, mods))
	}()
	wg.Wait()
}

func listenHotkey(onKeyDown func(), key hotkey.Key, mods []hotkey.Modifier,) (err error) {
	ms := []hotkey.Modifier{}
	ms = append(ms, mods...)
	hk := hotkey.New(ms, key)

	check(hk.Register())
	
	// Blocks until the hokey is triggered.
	for {
		<-hk.Keydown()
		onKeyDown()
	}
}