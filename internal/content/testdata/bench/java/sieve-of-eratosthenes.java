class Solution {
    public long[] primes_up_to_n(long n) {
        if (n < 2) return new long[0];
        boolean[] isComposite = new boolean[(int)n + 1];
        for (long i = 2; i * i <= n; i++) {
            if (!isComposite[(int)i]) {
                for (long j = i * i; j <= n; j += i) {
                    isComposite[(int)j] = true;
                }
            }
        }
        java.util.List<Long> primes = new java.util.ArrayList<>();
        for (int i = 2; i <= (int)n; i++) {
            if (!isComposite[i]) primes.add((long)i);
        }
        long[] out = new long[primes.size()];
        for (int i = 0; i < out.length; i++) out[i] = primes.get(i);
        return out;
    }
}
