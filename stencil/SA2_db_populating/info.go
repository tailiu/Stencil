package SA2_db_populating

/**
 *
 *						 1K dataset (stencil_exp_sa2_1k) (Total Row Count: 393,308) (Time used: 2h6m1.215291933s)
 *
 * 			{1, 7}, 		// 1. aspects						(4,000) (done) 
 *			{7, 9},			// 2. comments						(13,477) (done)
 *			{9, 10},		// 3. contacts						(35,502) (done) 
 *			{10, 11},		// 4. conversations					(2,402) (done)
 *			{11, 13},		// 5. messages						(39,646) (done)
 *			{13, 14},		// 6. notification_actors			(96,809) (done)
 *			{14, 19},		// 7. notifications					(96,809) (done)
 *			{19, 20},		// 8. people						(1,000) (done)
 *			{20, 26},		// 9. photos						(3,693) (done)
 *			{26, 27},		// 10. posts						(7,550) (done)
 *			{27, 30},		// 11. profiles						(1,000) (1,000)
 *			{30, 32}, 		// 12. share_visibilities			(7,550) (done) 
 *			{32, 35},		// 13. aspect_visibilities			(14,889) (done)
 *			{35, 39},		// 14. users						(1,000) (done) 
 *			{39, 41},		// 15. conversation_visibilities	(4,804) (done)
 *			{41, 52},		// 16. likes						(63,177) (done)
 *			{52, 198},		// 17. all other tables
 *
 *						10K dataset (stencil_exp_sa2_10k) (Total Row Count: 3,577,133, Unique: 3577128)
 *										 (Time used: 11m26.24976412s)
 *
 * 			{1, 7}, 		// 1. aspects						(40,000) (done) 
 *			{7, 9},			// 2. comments						(134,896) (done)
 *			{9, 10},		// 3. contacts						(357,506) (done) 
 *			{10, 11},		// 4. conversations					(16,601) (done)
 *			{11, 13},		// 5. messages						(35,853) (done)
 *			{13, 14},		// 6. notification_actors			(976,348) (done)
 *			{14, 19},		// 7. notifications					(976,348) (done)
 *			{19, 20},		// 8. people						(10,000) (done)
 *			{20, 26},		// 9. photos						(36,934) (done)
 *			{26, 27},		// 10. posts						(75,606) (done)
 *			{27, 30},		// 11. profiles						(10,000) (1,000)
 *			{30, 32}, 		// 12. share_visibilities			(75,606) (done) 
 *			{32, 35},		// 13. aspect_visibilities			(148,312) (done)
 *			{35, 39},		// 14. users						(10,000) (done) 
 *			{39, 41},		// 15. conversation_visibilities	(33,202) (done)
 *			{41, 52},		// 16. likes						(639,921) (done)
 *			{52, 198},		// 17. all other tables
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *										100K dataset (stencil_exp_sa2_100k)
 *
 * 
 * 		{1, 7}, 		// 1. aspects						(401,924) (done) 
 *		{7, 9},			// 2. comments						(1,348,745) (done)
 *		{9, 10},		// 3. contacts						(519,572) (done) 
 *		{10, 11},		// 4. conversations					(8,151) (done)
 *		{11, 13},		// 5. messages						(540,120) (done)
 *		{13, 14},		// 6. notification_actors			(4,672,224) (done) My machine 
 *		{14, 19},		// 7. notifications					(4,672,224) (done) My machine
 *		{19, 20},		// 8. people						(100,481) (done)
 *		{20, 26},		// 9. photos						(369,170) (done)
 *		{26, 27},		// 10. posts						(756,700) (done)
 *		{27, 30},		// 11. profiles						(100,481) (done)
 *		{30, 32}, 		// 12. share_visibilities			(756,700) (done) 
 *		{32, 35},		// 13. aspect_visibilities			(1,482,308) (done) My machine
 *		{35, 39},		// 14. users						(100,481) (done) 
 *		{39, 41},		// 15. conversation_visibilities	(16,302) (done)
 *		{41, 52},		// 16. likes						(3,055,543) (done) (11) My machine
 *		{52, 198},		// 17. all other tables
 *
 *							
 * 
 * 						   
 


							 1M dataset (stencil_exp_sa2_3)

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
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
 *
**/