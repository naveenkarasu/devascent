function primes_up_to_n(n) {
    if (n < 2) return [];
    let isPrime = new Array(n + 1).fill(true);
    isPrime[0] = isPrime[1] = false;
    for (let i = 2; i * i <= n; i++) {
        if (isPrime[i]) {
            for (let j = i * i; j <= n; j += i) {
                isPrime[j] = false;
            }
        }
    }
    let result = [];
    for (let i = 2; i <= n; i++) {
        if (isPrime[i]) result.push(i);
    }
    return result;
}
