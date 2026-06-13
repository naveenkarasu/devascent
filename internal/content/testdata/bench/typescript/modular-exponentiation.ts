function modular_exponentiation(x: number, n: number, m: number): number {
    let res = 1;
    x = x % m;
    while (n > 0) {
        if (n % 2 !== 0) {
            res = (res * x) % m;
            n--;
        } else {
            x = (x * x) % m;
            n = Math.floor(n / 2);
        }
    }
    return res;
}
