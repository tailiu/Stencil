package SA2_display

import (
	"stencil/config"
)

func (display *display) getADataInOwner(hints []*HintStruct, ownership *config.Ownership) (*HintStruct, error) {

	tag := hints[0].Tag

	pTag := "root"

	procConditions := getProcConditions(display, tag, pTag, ownership.Conditions)

	return display.getHintInParentNode(hints, procConditions, pTag)

}

func (display *display) getOwner(hints []*HintStruct, ownership *config.Ownership) ([]*HintStruct, error) {
	
	oneDataInOwnerNode, err := display.getADataInOwner(hints, ownership)
	if err != nil {
		return nil, err
	}

	return display.GetDataInNodeBasedOnDisplaySetting(oneDataInOwnerNode)
	
}