package evaluation

func dstDanglingDataSystem(evalConfig *EvalConfig) map[string]int64 {
	danglingDataStats := make(map[string]int64)

	getDanglingLikesNumSystem(evalConfig, danglingDataStats)
	getDanglingCommentsNumSystem(evalConfig, danglingDataStats)
	getDanglingMessagesNumSystem(evalConfig, danglingDataStats)

	return danglingDataStats
}