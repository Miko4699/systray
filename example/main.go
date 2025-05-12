package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Miko4699/systray"
	"github.com/Miko4699/systray/icon"
)
func main() {
	MainRun()
}
var start func()
var end func()
func MainRun() {
	onExit := func() {
		now := time.Now()
		fmt.Println("Exit at", now.String())
	}
	systray.Run(onReady, onExit)
}
func addQuitItem() {
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.Enable()
	mQuit.Click(func() {
		fmt.Println("Requesting quit")
		systray.Quit()
		//systray.Quit()// macos error
		//end() // macos error
		fmt.Println("Finished quitting")
	})
}
func onReady() {
	fmt.Println("systray.onReady")
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Energy Sys Tray")
	systray.SetTooltip("Energy tooltip")
	systray.SetOnClick(func(menu systray.IMenu) {
		if menu != nil { // menu for linux nil
			menu.ShowMenu()
		}
		fmt.Println("SetOnClick")
	})
	systray.SetOnDClick(func(menu systray.IMenu) {
		if menu != nil { // menu for linux nil
			menu.ShowMenu()
		}
		fmt.Println("SetOnDClick")
	})
	// OnRClick linux not impl
	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
		fmt.Println("SetOnRClick")
	})
	systray.CreateMenu()
	addQuitItem()
	systray.SetTemplateIcon(icon.Data, icon.Data)
	
	// 添加显示通知的菜单项
	mNotify := systray.AddMenuItem("Show Notification", "Display a system tray notification")
	mNotify.Click(func() {
		fmt.Println("Showing notification")
		err := systray.ShowNotification("Systray Notification", "这是一条来自系统托盘的通知!")
		if err != nil {
			fmt.Println("Failed to show notification:", err)
		}
	})
	
	mChange := systray.AddMenuItem("Change Me", "Change Me")
	mChecked := systray.AddMenuItemCheckbox("Checked", "Check Me", true)
	mEnabled := systray.AddMenuItem("Enabled", "Enabled")
	// Sets the icon of a menu item. Only available on Mac.
	mEnabled.SetTemplateIcon(icon.Data, icon.Data)
	systray.AddMenuItem("Ignored", "Ignored")
	subMenuTop := systray.AddMenuItem("SubMenuTop", "SubMenu Test (top)")
	subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
	subMenuBottom := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
	subMenuBottom2 := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")
	subMenuBottom2.SetIcon(icon.Data)
	systray.AddSeparator()
	mToggle := systray.AddMenuItem("Toggle", "Toggle some menu items")
	shown := true
	toggle := func() {
		if shown {
			subMenuBottom.Check()
			subMenuBottom2.Hide()
			mEnabled.Hide()
			shown = false
			mEnabled.Disable()
		} else {
			subMenuBottom.Uncheck()
			subMenuBottom2.Show()
			mEnabled.Show()
			mEnabled.Enable()
			shown = true
		}
	}
	mReset := systray.AddMenuItem("Reset", "Reset all items")
	mChange.Click(func() {
		mChange.SetTitle("I've Changed")
		// 在更改标题时显示通知
		systray.ShowNotification("Title Changed", "菜单项的标题已更改!")
	})
	mChecked.Click(func() {
		if mChecked.Checked() {
			mChecked.Uncheck()
			mChecked.SetTitle("Unchecked")
			systray.ShowNotification("Checkbox Changed", "复选框已取消选中!")
		} else {
			mChecked.Check()
			mChecked.SetTitle("Checked")
			systray.ShowNotification("Checkbox Changed", "复选框已选中!")
		}
	})
	mEnabled.Click(func() {
		mEnabled.SetTitle("Disabled")
		fmt.Println("mEnabled.Disabled()", mEnabled.Disabled())
		mEnabled.Disable()
		systray.ShowNotification("Item Disabled", "菜单项已禁用!")
	})
	subMenuBottom2.Click(func() {
		systray.ShowNotification("Warning", "即将触发 panic!")
		time.Sleep(time.Second)
		panic("panic button pressed")
	})
	subMenuBottom.Click(func() {
		toggle()
		systray.ShowNotification("Toggle Changed", "切换状态已更改!")
	})
	mReset.Click(func() {
		systray.ResetMenu()
		addQuitItem()
		systray.ShowNotification("Menu Reset", "菜单已重置!")
	})
	mToggle.Click(func() {
		toggle()
		systray.ShowNotification("Items Toggled", "菜单项已切换显示状态!")
	})
	// tray icon switch
	go func() {
		var b bool
		// demo: to png full path
		wd, _ := os.Getwd()
		wd = strings.Replace(wd, "example", "", -1)
		wd = filepath.Join(wd, "icon")
		fmt.Println("wd", wd) // /to/icon/path/icon.png, logo.png
		var ext = ".png"
		if runtime.GOOS == "windows" {
			ext = ".ico" // windows .ico
		}
		icoData, _ := ioutil.ReadFile(filepath.Join(wd, "icon"+ext))
		logoData, _ := ioutil.ReadFile(filepath.Join(wd, "logo"+ext))
		for true {
			time.Sleep(time.Second * 1)
			b = !b
			if b {
				systray.SetIcon(logoData)
			} else {
				systray.SetIcon(icoData)
			}
		}
	}()
}