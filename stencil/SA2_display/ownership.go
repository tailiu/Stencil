package SA2_display

import (
	"stencil/config"
)

func (displayConfig *displayConfig) getADataInOwner(hints []*HintStruct, ownership *config.Ownership) (*HintStruct, error) {

	tag := hints[0].Tag

	pTag := "root"

	procConditions := getProcConditions(displayConfig, tag, pTag, 
		ownership.Conditions)

	return displayConfig.getHintInParentNode(hints, procConditions, pTag)

}

func (displayConfig *displayConfig) getOwner(hints []*HintStruct, ownership *config.Ownership) ([]*HintStruct, error) {
	
	oneDataInOwnerNode, err := displayConfig.getADataInOwner(hints, ownership)
	if err != nil {
		return nil, err
	}

	return displayConfig.GetDataInNodeBasedOnDisplaySetting(oneDataInOwnerNode)
	
}