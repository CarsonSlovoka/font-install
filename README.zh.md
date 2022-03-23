## Windows

此份文件主要針對windows所撰寫

字體安裝步驟

1. 寫註冊檔，又分兩類，一種安裝在機器上(所有使用者都能用)、另一種只針對特定使用者安裝
  - 安裝在機器上:
    - `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`
    - `字體名稱 (TrueType)`   `REG_SZ` 、 `basename(字型檔.{oft,ttf})`
        > 由於他已經認定在fonts<sup>`path.Join(os.Getenv("WINDIR"), "Fonts")`</sup>資料夾，因此值只要寫basename(含附檔名)即可
  - 安裝在某使用者:
    - `HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`
    - `字體名稱 (TrueType)`   `REG_SZ` 、 `filepath.Join(os.Getenv("USERPROFILE"), "AppData/Local/Microsoft/Windows/Fonts/[YOUR.{ttf, otf}]")`
      - (大致上和localMachine相同，差別在於一般而言會把它存放到`%USERPROFILE%/AppData/Local/Microsoft/Windows/Fonts/`之下，並且需要提供完整路徑)
2. 將字型檔放到指定的位子
  - 安裝在機器: `%WINDIR%/Fonts/`
  - 特定使用者: `%USERPROFILE%/AppData/Local/Microsoft/Windows/Fonts/`
3. 呼叫 `Gdi32.dll`的`AddFontResource("C://example.ttf")`
    > 此步驟如果沒有，在其他應用程式仍然看不到字體
5. (可選) 呼叫 `user32.dll`的`SendMessage(HWND_BROADCAST, WM_FONTCHANGE, 0, 0)`

## 測試

```
go build -o test.exe
test.exe -fromFile="install-font.lst"
```

其中: `install-font.lst` 放置所有您想要安裝字體的路徑 (可以是相對路徑或絕對路徑)

> ❗ 注意，確保您在admin的權限下<sup>(寫機碼需要權限)</sup>，執行程式。

output
```
2022/03/23 17:41:42 Installing "Steel City Comic" to C:\WINDOWS/Fonts/scb.ttf
```

其中 Installing 後面所接的就是您在各個應用程式上會看到的名稱，以上面的範例而言就是`Steel City Comic`

## 參考資料

- [win32/api/wingdi/nf-wingdi-addfontresourcew](https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-addfontresourcew)
- [windows/win32/gdi/wm-fontchange](https://docs.microsoft.com/en-us/windows/win32/gdi/wm-fontchange)
- [Installing font and making Windows aware](https://stackoverflow.com/a/59112398/9935654)
- https://www.digitaltrends.com/computing/how-to-install-fonts-in-windows-10/
