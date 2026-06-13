import java.util.*;

class Solution {
    public long[] palindrome_primes_in_range(long a, long b) {
        List<Long> result = new ArrayList<>();
        for (long i = a; i <= b; i++) {
            if (isPrime(i) && isPalindrome(i)) result.add(i);
        }
        long[] arr = new long[result.size()];
        for (int i = 0; i < result.size(); i++) arr[i] = result.get(i);
        return arr;
    }

    private boolean isPrime(long n) {
        if (n < 2) return false;
        for (long i = 2; i * i <= n; i++) {
            if (n % i == 0) return false;
        }
        return true;
    }

    private boolean isPalindrome(long n) {
        String s = Long.toString(n);
        int l = 0, r = s.length() - 1;
        while (l < r) {
            if (s.charAt(l) != s.charAt(r)) return false;
            l++; r--;
        }
        return true;
    }
}
