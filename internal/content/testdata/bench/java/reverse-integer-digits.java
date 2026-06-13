class Solution {
    public long reverse_integer(long n) {
        long reversed = 0;
        while (n > 0) {
            reversed = reversed * 10 + n % 10;
            n /= 10;
        }
        return reversed;
    }
}
