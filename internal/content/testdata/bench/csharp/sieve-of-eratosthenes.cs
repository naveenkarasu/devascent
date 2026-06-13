public class Solution {
    public long[] primes_up_to_n(long n) {
        if (n < 2) return new long[0];
        bool[] isPrime = new bool[n + 1];
        for (long i = 2; i <= n; i++) isPrime[i] = true;
        for (long i = 2; i * i <= n; i++) {
            if (isPrime[i]) {
                for (long j = i * i; j <= n; j += i) isPrime[j] = false;
            }
        }
        var result = new List<long>();
        for (long i = 2; i <= n; i++) {
            if (isPrime[i]) result.Add(i);
        }
        return result.ToArray();
    }
}
