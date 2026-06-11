class Solution {
    public long[] best_team_coverage(long[][] topics) {
        int n = topics.length;
        int k = topics[0].length;
        long maxTopics = 0, bestCount = 0;
        for (int i = 0; i < n; i++) {
            for (int j = i + 1; j < n; j++) {
                long covered = 0;
                for (int t = 0; t < k; t++) {
                    if (topics[i][t] == 1 || topics[j][t] == 1) covered++;
                }
                if (covered > maxTopics) {
                    maxTopics = covered;
                    bestCount = 1;
                } else if (covered == maxTopics) {
                    bestCount++;
                }
            }
        }
        return new long[]{maxTopics, bestCount};
    }
}
