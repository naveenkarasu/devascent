public class Solution {
    public long nth_prime(long n) {
        var primes = new System.Collections.Generic.List<long>();
        long i = 2;
        while (primes.Count < n) {
            bool isPrime = true;
            foreach (long p in primes) {
                if (p * p > i) break;
                if (i % p == 0) { isPrime = false; break; }
            }
            if (isPrime) primes.Add(i);
            i++;
        }
        return primes[(int)(n - 1)];
    }
}
