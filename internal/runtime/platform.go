package runtime

import "runtime"

func GetPlatform() string {
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}

func IsPosix() bool {
	goos := runtime.GOOS
	return goos == "linux" || goos == "darwin"
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}
