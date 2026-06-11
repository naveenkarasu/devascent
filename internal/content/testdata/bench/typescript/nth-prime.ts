function nth_prime(n: number): number {
    const primes: number[] = [];
    let i = 2;
    while (primes.length < n) {
        if (primes.every(p => i % p !== 0)) {
            primes.push(i);
        }
        i++;
    }
    return primes[primes.length - 1];
}
