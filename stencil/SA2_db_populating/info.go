package SA2_db_populating

/**
 *							100K dataset (stencil_exp_sa2_100k)
 * 
 *							(users, profiles, people, conversation_visibilities, conversations, photos,
 *								messages(0 - 200,000), posts(0 - 200,000))
 * 
 * 		{1, 7}, 		// 1. aspects						(401,924) 
 *		{7, 9},			// 2. comments						(1,348,745) (13) My machine (0 - 200,000)
 *		{9, 10},		// 3. contacts						(519,572) * (11)
 *		{10, 11},		// 4. conversations					(8,151) (done)
 *		{11, 13},		// 5. messages						(540,120) * (12) My machine (200,000 - 400,000)
 *		{13, 14},		// 6. notification_actors			(4,672,224)
 *		{14, 19},		// 7. notifications					(4,672,224)
 *		{19, 20},		// 8. people						(100,481) (done)
 *		{20, 26},		// 9. photos						(369,170) (done)
 *		{26, 27},		// 10. posts						(756,700) * (10) My machine (200,000 - 400,000)
 *		{27, 30},		// 11. profiles						(100,481) (done)
 *		{30, 32}, 		// 12. share_visibilities			(756,700)
 *		{32, 35},		// 13. aspect_visibilities			(1,482,308)
 *		{35, 39},		// 14. users						(100,481) (done) 
 *		{39, 41},		// 15. conversation_visibilities	(16,302) (done)
 *		{41, 52},		// 16. likes						(3,055,543)
 *		{52, 198},		// 17. all other tables
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
**/