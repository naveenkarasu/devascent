import java.util.*;

class Solution {
    public long nth_prime(long n) {
        List<Long> primes = new ArrayList<>();
        long i = 2;
        while (primes.size() < n) {
            boolean isPrime = true;
            for (long p : primes) {
                if (p * p > i) break;
                if (i % p == 0) { isPrime = false; break; }
            }
            if (isPrime) primes.add(i);
            i++;
        }
        return primes.get((int)(n - 1));
    }
}
