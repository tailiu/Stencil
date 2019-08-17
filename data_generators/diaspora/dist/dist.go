package dist

import (
	"diaspora/helper"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Range struct {
	lower_bound, upper_bound int
	pshare                   float64
}

func GetTweetRanges() []Range {
	var ranges []Range
	ranges = append(ranges, Range{1, 10, 0.2})
	ranges = append(ranges, Range{11, 50, 0.28})
	ranges = append(ranges, Range{51, 150, 0.37})
	ranges = append(ranges, Range{151, 1000, 0.15})
	return ranges
}

func GetFriendRanges() []Range {
	// Buckets: 32% <100, 15% 100-200, 13% 200-300, 8% 300-400, 7% 400-500, 25% 500+
	var ranges []Range
	ranges = append(ranges, Range{1, 100, 0.32})
	ranges = append(ranges, Range{100, 200, 0.15})
	ranges = append(ranges, Range{200, 300, 0.13})
	ranges = append(ranges, Range{300, 400, 0.08})
	ranges = append(ranges, Range{400, 500, 0.07})
	ranges = append(ranges, Range{500, 600, 0.25})
	return ranges
}

func GetBoundsForBuckets(bucketsGenerated, bucketsRequired map[string]int) (int, int) {
	lower_bound, upper_bound := -1, -1
	if bucketsGenerated["500+"] < bucketsRequired["500+"] {
		lower_bound, upper_bound = 500, 600
	} else if bucketsGenerated["400-500"] < bucketsRequired["400-500"] {
		lower_bound, upper_bound = 400, 499
	} else if bucketsGenerated["300-400"] < bucketsRequired["300-400"] {
		lower_bound, upper_bound = 300, 399
	} else if bucketsGenerated["200-300"] < bucketsRequired["200-300"] {
		lower_bound, upper_bound = 200, 299
	} else if bucketsGenerated["100-200"] < bucketsRequired["100-200"] {
		lower_bound, upper_bound = 100, 199
	} else if bucketsGenerated["<100"] < bucketsRequired["<100"] {
		lower_bound, upper_bound = 1, 99
	}
	return lower_bound, upper_bound
}

func CheckFriendsDistribution(friendcount map[int]int) {

	o, i, j, k, l, m, n := 0, 0, 0, 0, 0, 0, 0

	for _, fnum := range friendcount {
		if fnum < 1 {
			o++
		} else if fnum >= 1 && fnum < 100 {
			i++
		} else if fnum >= 100 && fnum < 200 {
			j++
		} else if fnum >= 200 && fnum < 300 {
			k++
		} else if fnum >= 300 && fnum < 400 {
			l++
		} else if fnum >= 400 && fnum < 500 {
			m++
		} else if fnum >= 500 {
			n++
		}
	}
	fmt.Println(fmt.Sprintf("friendcount <100: %3d; 100-200: %3d; 200-300: %3d; 300-400: %3d; 400-500: %3d; 500-600: %3d; | 0: %3d", i, j, k, l, m, n, o))
}

func VerifyFriendsDistribution(friendlist map[int][]int) bool {

	o, i, j, k, l, m, n := 0, 0, 0, 0, 0, 0, 0

	for _, friends := range friendlist {
		fnum := len(friends)
		if fnum < 1 {
			o++
		} else if fnum >= 1 && fnum < 100 {
			i++
		} else if fnum >= 100 && fnum < 200 {
			j++
		} else if fnum >= 200 && fnum < 300 {
			k++
		} else if fnum >= 300 && fnum < 400 {
			l++
		} else if fnum >= 400 && fnum < 500 {
			m++
		} else if fnum >= 500 {
			n++
		}
	}
	fmt.Println(fmt.Sprintf("friendlist  <100: %3d; 100-200: %3d; 200-300: %3d; 300-400: %3d; 400-500: %3d; 500-600: %3d; | 0: %3d", i, j, k, l, m, n, o))

	return i > 0 && j > 0 && k > 0 && l > 0 && m > 0 && n > 0
}

func GenFriendsCountForUsers(totalUsers int) map[int]int {
	friendcount := make(map[int]int)
	for fRange, unum := range GetFriendsBuckets(totalUsers) {
		rTokens := strings.Split(fRange, "-")
		lower_bound, _ := strconv.Atoi(rTokens[0])
		upper_bound, _ := strconv.Atoi(rTokens[1])
		ucount := 0
		for uid := 0; uid < totalUsers; uid++ {
			if _, ok := friendcount[uid]; !ok {
				friendcount[uid] = helper.RandomNumber(lower_bound, upper_bound-1)
				ucount++
			}
			if ucount == unum {
				break
			}
		}
	}
	return friendcount
}

func AssignFriendsToUsers(totalUsers int) map[int][]int {
	friendcount := GenFriendsCountForUsers(totalUsers)
	friendlist := make(map[int][]int)
	for uid := 0; uid < totalUsers; uid++ {
		for _, fid := range rand.Perm(totalUsers) {
			if uid == fid {
				continue
			}
			// check if friend already has enough friends of their own
			if friendIDs, ok := friendlist[fid]; ok && len(friendIDs) >= friendcount[fid] {
				continue
			}
			if friendIDs, ok := friendlist[uid]; ok && len(friendIDs) >= friendcount[uid] {
				break
			} else {
				if !IdExistsInList(fid, friendlist[uid]) && !IdExistsInList(uid, friendlist[fid]) {
					friendlist[uid] = append(friendlist[uid], fid)
					friendlist[fid] = append(friendlist[fid], uid)
				}
			}
		}
	}

	// checkFriendsDistribution(friendcount)
	// verifyFriendsDistribution(friendlist)

	return friendlist
}

func GetFriendsBuckets(totalUsers int) map[string]int {
	buckets := make(map[string]int)
	for _, fRange := range GetFriendRanges() {
		key := fmt.Sprintf("%d-%d", fRange.lower_bound, fRange.upper_bound)
		buckets[key] = int(fRange.pshare * float64(totalUsers))
	}
	return buckets
}

func AssignPostsToUsers(totalPosts int) map[int]int {
	users := make(map[int]int)
	for _, r := range GetTweetRanges() {
		t, numOfPosts := 0, int(math.Ceil(float64(totalPosts)*r.pshare))
		for {
			for i := r.lower_bound - 1; i < r.upper_bound && t < numOfPosts; i++ {
				if _, ok := users[i]; !ok {
					users[i] = 0
				}
				users[i]++
				t++
			}
			if t >= numOfPosts {
				break
			}
		}
	}
	return users
}

func IdExistsInList(id int, list []int) bool {
	for _, pID := range list {
		if pID == id {
			return true
		}
	}
	return false
}
