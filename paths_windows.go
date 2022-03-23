package main

import (
	"os"
	"path"
	"path/filepath"
)

// FontsDir denotes the path to the user's fonts directory on Linux.
// Windows doesn't have the concept of a permanent, per-user collection
// of fonts, meaning that all fonts are stored in the system-level fonts
// directory, which is %WINDIR%\Fonts by default.
var (
	FontsDir     = path.Join(os.Getenv("WINDIR"), "Fonts")
	UserFontsDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData/Local/Microsoft/Windows/Fonts")
)
