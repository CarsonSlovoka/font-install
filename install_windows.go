package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"syscall"
)

func registerFont(rootKey registry.Key, name, value string) error {
	key, err := registry.OpenKey(rootKey, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`, registry.WRITE)
	if err != nil {
		return err
	}
	defer key.Close()

	if err = key.SetStringValue(
		fmt.Sprintf("%v (TrueType)", name), // it's "ok" to mark an OpenType font as "TrueType"
		value,
	); err != nil {
		return err
	}
	return nil
}

func notifySysFontChange(fontpath string) (err error) {
	var gdi32dll = syscall.NewLazyDLL("Gdi32.dll")
	var procAddFontResource = gdi32dll.NewProc("AddFontResourceW")

	rtnValue, _, _ := procAddFontResource.Call(LPCWSTR(fontpath))

	if rtnValue == 0 {
		return fmt.Errorf("AddFontResourceW error: \n%s", err)
	}

	user32dll := syscall.NewLazyDLL("user32.dll")
	// procSendMessage := user32dll.NewProc("SendMessageW") // 這個會等待所有程式都回應，所以你會覺得程式好像掛掉了都不動 // https://social.msdn.microsoft.com/Forums/en-US/6900f74f-6ece-47da-88fc-f9c8bcd40206/sendmessage-api-slow?forum=wpf
	// 也可以考慮 SendMessageTimeout(HWND_BROADCAST, WM_FONTCHANGE, 0, 0, SMTO_ABORTIFHUNG, 500) // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessagetimeouta
	procPostMessage := user32dll.NewProc("PostMessageW")
	_, _, err = procPostMessage.Call(HWND_BROADCAST, WM_FONTCHANGE, 0, 0)
	if err != syscall.Errno(0x0) {
		return err
	}
	return nil
}

func installOnHKCU(fontData *FontData) (err error) {
	// Reg
	userFontPath := filepath.Join(UserFontsDir, fontData.FileName)
	if err = registerFont(registry.CURRENT_USER, fontData.Name, userFontPath); err != nil {
		return err
	}

	// File
	err = ioutil.WriteFile(userFontPath, fontData.Data, 0644)
	if err != nil {
		return err
	}
	log.Printf("Installing \"%v\" to %v", fontData.Name, userFontPath)

	// tell system font change
	return notifySysFontChange(userFontPath)
}

func installOnHKLM(fontData *FontData) (err error) {
	if err = registerFont(registry.LOCAL_MACHINE, fontData.Name, fontData.FileName); err != nil { // 如果是安裝在local machine只要寫檔名即可，他會認定已經存在在FontsDir的路徑之中
		return err
	}
	sysFontPath := path.Join(FontsDir, fontData.FileName)
	err = ioutil.WriteFile(sysFontPath, fontData.Data, 0644)
	if err != nil {
		return err
	}
	log.Printf("Installing \"%v\" to %v", fontData.Name, sysFontPath)
	return notifySysFontChange(sysFontPath)
}

func platformDependentInstall(fontData *FontData) (err error) {
	/*
		To install a font on Windows:
		1. Create a registry entry for the font
		2. Copy the file to the fonts directory
		3. Call Gdi32.dll.AddFontResource
		4. Call user32.dll.PostMessage
	*/

	if err = installOnHKLM(fontData); err != nil {
		return err
	}

	return nil
}
