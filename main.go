package main

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
)

func main() {
	runtime.LockOSThread()

	cocoa.TerminateAfterWindowsClose = false
	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle("Docky")

		moveClicked := make(chan bool)
		go func() {
			for {
				select {
				case <-moveClicked:
					log.Printf("clicked Move")

					cmd1 := exec.Command("defaults", "write", "com.apple.dock", "orientation", "bottom")
					log.Printf("Running command \"defaults write com.apple.dock orientation bottom\"...")
					err1 := cmd1.Run()
					log.Printf("Command finished with error: %v", err1)

					cmd2 := exec.Command("killall", "Dock")
					log.Printf("Running command \"killall Dock\"...")
					err2 := cmd2.Run()
					log.Printf("Command finished with error: %v", err2)
				}
			}
		}()

		itemMove := cocoa.NSMenuItem_New()
		itemMove.SetTitle("Move dock to laptop screen")
		itemMove.SetAction(objc.Sel("moveClicked:"))
		cocoa.DefaultDelegateClass.AddMethod("moveClicked:", func(_ objc.Object) {
			moveClicked <- true
		})

		itemQuit := cocoa.NSMenuItem_New()
		itemQuit.SetTitle("Quit")
		itemQuit.SetAction(objc.Sel("terminate:"))

		menu := cocoa.NSMenu_New()
		menu.AddItem(itemMove)
		menu.AddItem(itemQuit)
		obj.SetMenu(menu)

	})
	app.Run()
}
