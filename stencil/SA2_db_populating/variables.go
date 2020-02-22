package SA2_db_populating


const maxRowID = 2147483647

const subPartionNum = 5

// 13 and 14 are the start of ranges
var subPartitionTables = []int {
	13, 14,
}

var ranges = [][]int {
	{1, 7}, 		// 1. aspects						(4,032,432)
	{7, 9},			// 2. comments						(13,481,411)
	{9, 10},		// 3. contacts						(5,191,420)
	{10, 11},		// 4. conversations					(81,119)
	{11, 13},		// 5. messages						(5,400,995)
	{13, 14},		// 6. notification_actors			(46,785,209)
	{14, 19},		// 7. notifications					(46,785,209)
	{19, 20},		// 8. people						(1,008,108)
	{20, 26},		// 9. photos						(3,692,680)
	{26, 27},		// 10. posts						(7,562,681)
	{27, 30},		// 11. profiles						(1,008,108)
	{30, 32}, 		// 12. share_visibilities			(7,562,681)
	{32, 35},		// 13. aspect_visibilities			(14,814,749)
	{35, 39},		// 14. users						(1,008,108)
	{39, 41},		// 15. conversation_visibilities	(162,238)
	{41, 52},		// 16. likes						(30,626,969)
	{52, 198},		// 17. all other tables
}

var tableNameRangeIndexMap = map[string]string{
	"aspects": "1",
	"comments": "2",
	"contacts": "3",
	"conversations": "4",
	"messages": "5",
	"notification_actors": "6",
	"notifications": "7",
	"people": "8",
	"photos": "9",
	"posts": "10",
	"profiles": "11",
	"share_visibilities": "12",
	"aspect_visibilities": "13",
	"users": "14",
	"conversation_visibilities": "15",
	"likes": "16",
}

var oldRanges = [][]int {
	{1, 7}, 		// 1. aspects						(4,032,432) *
	{7, 9},			// 2. comments						(13,481,411)
	{9, 10},		// 3. contacts						(5,191,420) *
	{10, 11},		// 4. conversations					(81,119) 
	{11, 13},		// 5. messages						(5,400,995) *
	{13, 14},		// 6. notification_actors			(46,785,209)
	{14, 19},		// 7. notifications					(46,785,209)
	{19, 20},		// 8. people						(1,008,108)
	{20, 26},		// 9. photos						(3,692,680)
	{26, 27},		// 10. posts						(7,562,681)
	{27, 32},		// 11. profiles	(1,008,108) && share_visibilities (7,562,681) *				
	{32, 35},		// 12. aspect_visibilities			(14,814,749) *
	{35, 39},		// 13. users						(1,008,108)
	{39, 41},		// 14. conversation_visibilities	(162,238)
	{41, 52},		// 15. likes						(30,626,969)
	{52, 198},		// 16. all other tables
}