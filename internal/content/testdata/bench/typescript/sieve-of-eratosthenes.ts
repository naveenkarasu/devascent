function primes_up_to_n(n: number): number[] {
    if (n < 2) return [];
    const is_prime: boolean[] = new Array(n + 1).fill(true);
    is_prime[0] = false;
    is_prime[1] = false;
    let i = 2;
    while (i * i <= n) {
        if (is_prime[i]) {
            for (let j = i * i; j <= n; j += i) {
                is_prime[j] = false;
            }
        }
        i++;
    }
    const result: number[] = [];
    for (let i = 2; i <= n; i++) {
        if (is_prime[i]) result.push(i);
    }
    return result;
}
