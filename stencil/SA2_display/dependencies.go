package SA2_display

import (
	"stencil/config"
)

func GetDisplaySettingInDependencies(appConfig *config.AppConfig, hint HintStruct, pTag string) (string, error) {
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
