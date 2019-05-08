package dependency_handler

import (
	"transaction/config"
	"transaction/display"
	// "transaction/db"
	// "fmt"
	// "errors"
	// "log"
)

func GetDisplaySettingInDependencies(appConfig *config.AppConfig, hint display.HintStruct, pTag string) (string, error) {
	tag, _ := hint.GetTagName(appConfig)
	setting, err := appConfig.GetDepDisplaySetting(tag, pTag)

	if err != nil {
		return "", err
	}

	if setting == "" {
		return "parent_node_complete_displays", nil
	} else {
		return setting, nil
	}
}