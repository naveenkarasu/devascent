type Twitter struct {
	time      int
	tweets    map[int][][2]int // user -> list of [time, tweetId]
	following map[int]map[int]bool
}

func newTwitter() *Twitter {
	return &Twitter{
		tweets:    make(map[int][][2]int),
		following: make(map[int]map[int]bool),
	}
}

func (t *Twitter) postTweet(user, tweet int) {
	t.tweets[user] = append(t.tweets[user], [2]int{t.time, tweet})
	t.time++
}

func (t *Twitter) getNewsFeed(user int) []int {
	type entry struct {
		time    int
		tweetId int
	}
	var all []entry
	// collect from user and following
	add := func(u int) {
		for _, tw := range t.tweets[u] {
			all = append(all, entry{tw[0], tw[1]})
		}
	}
	add(user)
	for f := range t.following[user] {
		add(f)
	}
	// sort descending by time
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].time > all[i].time {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	limit := 10
	if len(all) < limit {
		limit = len(all)
	}
	res := make([]int, limit)
	for i := 0; i < limit; i++ {
		res[i] = all[i].tweetId
	}
	return res
}

func (t *Twitter) follow(a, b int) {
	if t.following[a] == nil {
		t.following[a] = make(map[int]bool)
	}
	t.following[a][b] = true
}

func (t *Twitter) unfollow(a, b int) {
	delete(t.following[a], b)
}

func twitter_ops(operations [][]interface{}) []interface{} {
	tw := newTwitter()
	out := make([]interface{}, len(operations))
	for i, op := range operations {
		switch op[0].(string) {
		case "postTweet":
			tw.postTweet(op[1].(int), op[2].(int))
			out[i] = nil
		case "getNewsFeed":
			out[i] = tw.getNewsFeed(op[1].(int))
		case "follow":
			tw.follow(op[1].(int), op[2].(int))
			out[i] = nil
		case "unfollow":
			tw.unfollow(op[1].(int), op[2].(int))
			out[i] = nil
		}
	}
	return out
}
