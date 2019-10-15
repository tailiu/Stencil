package app_dependency_handler

import (
	"stencil/config"
	"stencil/app_display"
)

func GetDisplaySettingInDependencies(appConfig *config.AppConfig, hint *app_display.HintStruct, pTag string) (string, error) {
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
