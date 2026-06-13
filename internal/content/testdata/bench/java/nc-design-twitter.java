import java.util.*;

class Solution {
    public Object[] twitter_ops(Object[][] operations) {
        long[] time = {0};
        Map<Long, List<long[]>> tweets = new HashMap<>();     // user -> list of {time, tweetId}
        Map<Long, Set<Long>> following = new HashMap<>();      // user -> followees
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            Object[] op = operations[i];
            String name = (String) op[0];
            if (name.equals("postTweet")) {
                long u = ((Number) op[1]).longValue();
                long t = ((Number) op[2]).longValue();
                tweets.computeIfAbsent(u, k -> new ArrayList<>()).add(new long[]{time[0]++, t});
                out[i] = null;
            } else if (name.equals("getNewsFeed")) {
                long u = ((Number) op[1]).longValue();
                Set<Long> people = new HashSet<>();
                people.add(u);
                if (following.containsKey(u)) people.addAll(following.get(u));
                List<long[]> all = new ArrayList<>();
                for (long p : people) {
                    if (tweets.containsKey(p)) all.addAll(tweets.get(p));
                }
                all.sort((a, b) -> Long.compare(b[0], a[0])); // newest first
                int n = Math.min(10, all.size());
                Object[] feed = new Object[n];
                for (int j = 0; j < n; j++) feed[j] = all.get(j)[1];
                out[i] = feed;
            } else if (name.equals("follow")) {
                long a = ((Number) op[1]).longValue();
                long b = ((Number) op[2]).longValue();
                following.computeIfAbsent(a, k -> new HashSet<>()).add(b);
                out[i] = null;
            } else { // unfollow
                long a = ((Number) op[1]).longValue();
                long b = ((Number) op[2]).longValue();
                if (following.containsKey(a)) following.get(a).remove(b);
                out[i] = null;
            }
        }
        return out;
    }
}
