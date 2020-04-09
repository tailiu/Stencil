package SA2_display

import (
	"stencil/config"
)

func getADataInOwner(displayConfig *displayConfig, hints []*HintStruct,
	ownership *config.Ownership) (*HintStruct, error) {

	tag := hints[0].Tag

	pTag := "root"

	procConditions := getProcConditions(displayConfig, tag, pTag, 
		ownership.Conditions)

	return getHintInParentNode(displayConfig, hints, procConditions, pTag)

}

func getOwner(displayConfig *displayConfig, hints []*HintStruct,
	ownership *config.Ownership) ([]*HintStruct, error) {
	
	oneDataInOwnerNode, err := getADataInOwner(displayConfig, hints, ownership)
	if err != nil {
		return nil, err
	}

	return displayConfig.GetDataInNodeBasedOnDisplaySetting(oneDataInOwnerNode)
	
}