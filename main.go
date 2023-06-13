package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"

	"github.com/sevlyar/go-daemon"
)

func main() {
	home := os.Getenv("HOME")
	cntxt := &daemon.Context{
		PidFileName: home+"/docky.pid",
		PidFilePerm: 0644,
		LogFileName: home+"/docky.log",
		LogFilePerm: 0640,
		WorkDir:     home,
		Umask:       027,
		Args:        []string{"[go-daemon docky]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")


	runtime.LockOSThread()

	cocoa.TerminateAfterWindowsClose = false
	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle("docky")

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
