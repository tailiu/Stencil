package dependency_handler

const getOneDataFromParentNodeAttemptTimes = 10

type DataInDependencyNode struct {
	Table 	string
	Data	map[string]string
}
