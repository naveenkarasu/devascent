function is_fibonacci(n: number): boolean {
    if (n <= 0) return false;
    let a = 1, b = 1;
    while (a < n) {
        const tmp = a;
        a = b;
        b = tmp + b;
    }
    return a === n;
}
