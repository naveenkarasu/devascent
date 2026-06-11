function twitter_ops(operations) {
    let time = 0;
    const tweets = new Map();   // user -> [[time, tweetId], ...]
    const following = new Map(); // user -> Set of followees

    function getTweets(user) {
        return tweets.get(user) || [];
    }
    function getFollowing(user) {
        return following.get(user) || new Set();
    }

    const out = [];
    for (const op of operations) {
        if (op[0] === "postTweet") {
            const user = op[1], tweet = op[2];
            if (!tweets.has(user)) tweets.set(user, []);
            tweets.get(user).push([time++, tweet]);
            out.push(null);
        } else if (op[0] === "getNewsFeed") {
            const user = op[1];
            const people = new Set(getFollowing(user));
            people.add(user);
            const all = [];
            for (const p of people) {
                for (const t of getTweets(p)) all.push(t);
            }
            all.sort((a, b) => b[0] - a[0]);
            out.push(all.slice(0, 10).map(t => t[1]));
        } else if (op[0] === "follow") {
            const a = op[1], b = op[2];
            if (!following.has(a)) following.set(a, new Set());
            following.get(a).add(b);
            out.push(null);
        } else {
            // unfollow
            const a = op[1], b = op[2];
            if (following.has(a)) following.get(a).delete(b);
            out.push(null);
        }
    }
    return out;
}
