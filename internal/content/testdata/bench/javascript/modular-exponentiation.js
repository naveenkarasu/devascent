function modular_exponentiation(x, n, m) {
    let res = 1;
    x = x % m;
    while (n > 0) {
        if (n % 2 !== 0) {
            res = (res * x) % m;
            n -= 1;
        } else {
            x = (x * x) % m;
            n = Math.floor(n / 2);
        }
    }
    return res;
}
