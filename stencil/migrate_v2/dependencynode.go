package migrate_v2

import (
	"errors"
	"fmt"
	"log"
	config "stencil/config/v2"
	"strings"
	"time"

	"github.com/gookit/color"
)

func (self DependencyNode) GetValueForKey(key string) (string, error) {

	// for _, datum := range self.Data {
	if _, ok := self.Data[key]; ok {
		switch v := self.Data[key].(type) {
		case nil:
			return "", nil
		case int, int64:
			val := fmt.Sprintf("%d", v)
			return val, nil
		case string:
			val := fmt.Sprintf("%s", v)
			return val, nil
		case bool:
			val := fmt.Sprintf("%t", v)
			return val, nil
		case time.Time:
			val := v.String()
			return val, nil
		default:
			val := v.(string)
			return val, nil
		}
	}
	// }
	return "", errors.New("No value found for " + key)
}

func (self *DependencyNode) Copy(node DependencyNode) {

	self.Tag = node.Tag
	self.SQL = node.SQL
	self.Data = node.Data
}

func (node *DependencyNode) ResolveParentDependencyConditions(dconditions []config.DCondition, parentTag config.Tag) (string, error) {

	conditionStr := ""
	for _, condition := range dconditions {
		tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
		if err != nil {
			color.Error.Println(err, node.Tag.Name, condition.TagAttr)
			log.Fatal("@ResolveParentDependencyConditions: tagAttr in condition doesn't exist? ", condition.TagAttr)
			break
		}
		if len(condition.Restrictions) > 0 {
			restricted := false
			for _, restriction := range condition.Restrictions {
				if restrictionAttr, err := node.Tag.ResolveTagAttr(restriction["col"]); err == nil {
					if val, ok := node.Data[restrictionAttr]; ok {
						if strings.EqualFold(fmt.Sprint(val), restriction["val"]) {
							restricted = true
						}
					} else {
						color.Error.Println(node.Data)
						log.Fatal("@ResolveParentDependencyConditions:", tagAttr, " doesn't exist in node data? ", node.Tag.Name)
					}
				} else {
					color.Error.Println(err)
					log.Fatal("@ResolveParentDependencyConditions: Col in restrictions doesn't exist? ", restriction["col"])
					break
				}
			}
			if restricted {
				return "", errors.New("Returning empty from restricted. Why?")
			}
		}
		depOnAttr, err := parentTag.ResolveTagAttr(condition.DependsOnAttr)
		if err != nil {
			color.Error.Println(err, parentTag.Name, condition.DependsOnAttr)
			log.Fatal("@ResolveParentDependencyConditions: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
			break
		}
		if val, ok := node.Data[tagAttr]; ok {
			if val == nil {
				return "", errors.New(fmt.Sprintf("trying to assign %s = %s, value is nil in node %s ", tagAttr, depOnAttr, node.Tag.Name))
			}
			if conditionStr != "" {
				conditionStr += " AND "
			}
			conditionStr += fmt.Sprintf("%s = '%v'", depOnAttr, val)
		} else {
			fmt.Println(node.Data)
			log.Fatal("ResolveParentDependencyConditions:", tagAttr, " doesn't exist in node data? ", node.Tag.Name)
		}
	}
	return conditionStr, nil
}

func (node *DependencyNode) ResolveDependencyConditions(SrcAppConfig config.AppConfig, dep config.Dependency, tag config.Tag) (string, error) {

	where := ""
	for _, depOn := range dep.DependsOn {
		if depOnTag, err := SrcAppConfig.GetTag(depOn.Tag); err == nil {
			if strings.EqualFold(depOnTag.Name, node.Tag.Name) {
				for _, condition := range depOn.Conditions {
					conditionStr := ""
					tagAttr, err := tag.ResolveTagAttr(condition.TagAttr)
					if err != nil {
						color.Error.Println(err, tag.Name, condition.TagAttr)
						log.Fatal("Stop and Check Dependencies!")
						break
					}
					depOnAttr, err := depOnTag.ResolveTagAttr(condition.DependsOnAttr)
					if err != nil {
						color.Error.Println(err, depOnTag.Name, condition.DependsOnAttr)
						log.Fatal("Stop and Check Dependencies!")
						break
					}
					if nodeVal, ok := node.Data[depOnAttr]; ok {
						if nodeVal == nil {
							return "", errors.New(color.FgLightRed.Render(fmt.Sprintf("trying to assign %s = %s, value is nil in node %s ", tagAttr, depOnAttr, node.Tag.Name)))
						}
						if conditionStr != "" || where != "" {
							conditionStr += " AND "
						}
						conditionStr += fmt.Sprintf("%s = '%v'", tagAttr, nodeVal)
					} else {
						color.Warn.Println(depOnTag)
						log.Fatal("ResolveDependencyConditions:", depOnAttr, " doesn't exist in ", depOnTag.Name)
					}
					if len(condition.Restrictions) > 0 {
						restrictions := ""
						for _, restriction := range condition.Restrictions {
							if restrictionAttr, err := tag.ResolveTagAttr(restriction["col"]); err == nil {
								if restrictions != "" {
									restrictions += " OR "
								}
								restrictions += fmt.Sprintf(" %s = '%s' ", restrictionAttr, restriction["val"])
							}

						}
						if restrictions == "" {
							color.Error.Println("Restrictions not resolved?")
							log.Fatal(condition.Restrictions)
						}
						conditionStr += fmt.Sprintf(" AND (%s) ", restrictions)
					}
					where += conditionStr
				}
			}
		}
	}
	return where, nil
}

func (node *DependencyNode) ResolveParentOwnershipConditions(own *config.Ownership, depOnTag config.Tag) (string, error) {

	where := ""
	for _, condition := range own.Conditions {
		conditionStr := ""
		tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
		if err != nil {
			fmt.Println("@ResolveParentOwnershipConditions > data1", node.Data)
			log.Fatal(err, node.Tag.Name, condition.TagAttr)
			break
		}
		depOnAttr, err := depOnTag.ResolveTagAttr(condition.DependsOnAttr)
		if err != nil {
			fmt.Println("@ResolveParentOwnershipConditions > data2", node.Data)
			log.Fatal(err, depOnTag.Name, condition.DependsOnAttr)
			break
		}
		if nodeVal, ok := node.Data[tagAttr]; ok {
			if nodeVal == nil {
				return "", errors.New(fmt.Sprintf("@ResolveParentOwnershipConditions > trying to assign %s = %s, value is nil in node %s ", depOnAttr, tagAttr, node.Tag.Name))
			}
			if conditionStr != "" || where != "" {
				conditionStr += " AND "
			}
			conditionStr += fmt.Sprintf("%s = '%v'", depOnAttr, node.Data[tagAttr])
		} else {
			fmt.Println("@ResolveParentOwnershipConditions > data3", node.Data)
			log.Fatal("@ResolveParentOwnershipConditions > ", tagAttr, " doesn't exist in ", node.Tag.Name)
		}
		where += conditionStr
	}
	return where, nil
}

func (root *DependencyNode) ResolveOwnershipConditions(own config.Ownership, tag config.Tag) (string, error) {

	where := ""
	for _, condition := range own.Conditions {
		conditionStr := ""
		tagAttr, err := tag.ResolveTagAttr(condition.TagAttr)
		if err != nil {
			fmt.Println("data1", root.Data)
			log.Fatal(err, tag.Name, condition.TagAttr)
			break
		}
		depOnAttr, err := root.Tag.ResolveTagAttr(condition.DependsOnAttr)
		if err != nil {
			fmt.Println("data2", root.Data)
			log.Fatal(err, tag.Name, condition.DependsOnAttr)
			break
		}
		if nodeVal, ok := root.Data[depOnAttr]; ok {
			if nodeVal == nil {
				return "", errors.New(fmt.Sprintf("trying to assign %s = %s, value is nil in node %s ", tagAttr, depOnAttr, root.Tag.Name))
			}
			if conditionStr != "" || where != "" {
				conditionStr += " AND "
			}
			conditionStr += fmt.Sprintf("%s = '%v'", tagAttr, root.Data[depOnAttr])
		} else {
			fmt.Println("data3", root.Data)
			log.Fatal("ResolveOwnershipConditions:", depOnAttr, " doesn't exist in ", tag.Name)
		}
		where += conditionStr
	}
	return where, nil
}

func (node *DependencyNode) DeleteMappedDataFromNode(mappedMemberData []MappedMemberData) {

	for _, mappedMemberDatum := range mappedMemberData {
		node.DeleteMappedDatumFromNode(mappedMemberDatum)
	}
}

func (node *DependencyNode) DeleteMappedDatumFromNode(mappedMemberDatum MappedMemberData) {

	var deletedCols []string
	for _, mmv := range mappedMemberDatum.Data {
		for key := range node.Data {
			if strings.Contains(key, ".id") {
				continue
			}
			if strings.EqualFold(key, mmv.GetMemberAttr()) {
				delete(node.Data, key)
				deletedCols = append(deletedCols, key)
			}
		}
	}
	if len(deletedCols) > 0 {
		fmt.Printf("Deleted Cols From Node | %v \n", deletedCols)
	}
}

func (node *DependencyNode) DeleteNulls() {

	for key, val := range node.Data {
		if val == nil {
			delete(node.Data, key)
		}
	}
}

func (node *DependencyNode) IsEmptyExcept() bool {
	return node.Data.IsEmptyExcept()
}
