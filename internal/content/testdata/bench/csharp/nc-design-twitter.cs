using System.Collections.Generic;

public class Solution {
    public object[] twitter_ops(object[][] operations) {
        long time = 0;
        var tweets = new Dictionary<long, List<long[]>>();  // user -> list of {time, tweetId}
        var following = new Dictionary<long, HashSet<long>>();
        var out_ = new object[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            var op = operations[i];
            string name = (string)op[0];
            if (name == "postTweet") {
                long u = (long)op[1];
                long t = (long)op[2];
                if (!tweets.ContainsKey(u)) tweets[u] = new List<long[]>();
                tweets[u].Add(new long[] { time++, t });
                out_[i] = null;
            } else if (name == "getNewsFeed") {
                long u = (long)op[1];
                var people = new HashSet<long> { u };
                if (following.ContainsKey(u)) foreach (long f in following[u]) people.Add(f);
                var all = new List<long[]>();
                foreach (long p in people) {
                    if (tweets.ContainsKey(p)) all.AddRange(tweets[p]);
                }
                all.Sort((a, b) => b[0].CompareTo(a[0]));
                int n = Math.Min(10, all.Count);
                var feed = new object[n];
                for (int j = 0; j < n; j++) feed[j] = all[j][1];
                out_[i] = feed;
            } else if (name == "follow") {
                long a = (long)op[1];
                long b = (long)op[2];
                if (!following.ContainsKey(a)) following[a] = new HashSet<long>();
                following[a].Add(b);
                out_[i] = null;
            } else { // unfollow
                long a = (long)op[1];
                long b = (long)op[2];
                if (following.ContainsKey(a)) following[a].Remove(b);
                out_[i] = null;
            }
        }
        return out_;
    }
}
