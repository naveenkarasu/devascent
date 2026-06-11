using System.Collections.Generic;

public class Solution {
    public long[] palindrome_primes_in_range(long a, long b) {
        bool isPrime(long n) {
            if (n < 2) return false;
            for (long i = 2; i * i <= n; i++) if (n % i == 0) return false;
            return true;
        }
        bool isPalindrome(long n) {
            string s = n.ToString();
            int l = 0, r = s.Length - 1;
            while (l < r) { if (s[l] != s[r]) return false; l++; r--; }
            return true;
        }
        var result = new List<long>();
        for (long i = a; i <= b; i++) {
            if (isPrime(i) && isPalindrome(i)) result.Add(i);
        }
        return result.ToArray();
    }
}
