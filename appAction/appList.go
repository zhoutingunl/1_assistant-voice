package appAction

import (
	"encoding/json"
	myredis "github.com/Awaken1119/assistant-voice/cache/redis"
	"github.com/Awaken1119/assistant-voice/db"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type AppInfo struct {
	Name string
	Path string
}

func GetAppList(apps []AppInfo) string {

	rConn := myredis.InitRedis().Get()
	defer rConn.Close()

	appList := make([]string, 0)

	for _, app := range apps {
		appList = append(appList, app.Name)
		_, err := rConn.Do("SET", app.Name, app.Path)
		if err != nil {
			return err.Error()
		}
		err = db.SaveApp(app.Name, app.Path)
		if err != nil {
			return err.Error()
		}
	}
	rConn.Do("SET", "appList", appList)
	result, _ := json.Marshal(appList)
	return string(result)
}

func getInstalledApps() []AppInfo {
	var apps []AppInfo
	roots := []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER}
	paths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for _, root := range roots {
		for _, path := range paths {
			k, err := registry.OpenKey(root, path, registry.READ)
			if err != nil {
				continue
			}
			defer k.Close()

			names, _ := k.ReadSubKeyNames(-1)
			for _, name := range names {
				sub, err := registry.OpenKey(k, name, registry.READ)
				if err != nil {
					continue
				}

				displayName, _, _ := sub.GetStringValue("DisplayName")
				if displayName == "" {
					sub.Close()
					continue
				}

				displayIcon, _, _ := sub.GetStringValue("DisplayIcon")

				exePath := cleanPath(displayIcon)

				if strings.HasSuffix(strings.ToLower(exePath), ".exe") {
					apps = append(apps, AppInfo{
						Name: displayName,
						Path: exePath,
					})
				}

				sub.Close()
			}
		}
	}
	return apps
}

func cleanPath(s string) string {
	// 去掉多余参数，如 `"C:\xxx\app.exe" /uninstall`
	s = strings.Trim(s, `"`)
	if idx := strings.Index(strings.ToLower(s), ".exe"); idx != -1 {
		return s[:idx+4]
	}
	return s
}

func init() {
	apps := getInstalledApps()
	GetAppList(apps)

	//fmt.Println(Apps)
	//fmt.Printf("应用列表：%s", AppNameList)
}
