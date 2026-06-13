class Solution {
    public long factorial(long n) {
        if (n <= 1) return 1;
        return n * factorial(n - 1);
    }
}
