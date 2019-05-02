package dependency_handler 

import (
	// "fmt"
	"transaction/config"
	"transaction/display"
)

func GetTobeCheckedDataInNode(appConfig *config.AppConfig, hint display.HintStruct) ([]display.HintStruct, error) {
	var data []display.HintStruct

	tagName, err := hint.GetTagName(appConfig.Tags)
	if err != nil {
		return data, err
	}
	// fmt.Println(tagName)

	displaySetting, _ := appConfig.GetTagDisplaySetting(tagName)
	if displaySetting == "default_display_setting" {
		if completeData, complete := CheckNodeComplete(appConfig, hint); complete {
			return completeData, nil
		} else {
			return data, nil
		}
	} else {

	}

	return data, nil
}