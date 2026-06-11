class Solution {
    public long trap_node(long n, long[][] edges) {
        long[] degree = new long[(int)n];
        for (long[] edge : edges) {
            degree[(int)edge[0]]++;
            degree[(int)edge[1]]++;
        }
        long max = 0;
        for (long d : degree) {
            if (d > max) max = d;
        }
        return max - 1;
    }
}
