function reverse_integer(x: number): number {
    const sign = x < 0 ? -1 : 1;
    x = Math.abs(x);
    let rev = 0;
    while (x > 0) {
        rev = rev * 10 + (x % 10);
        x = Math.floor(x / 10);
    }
    rev *= sign;
    if (rev > 2147483647 || rev < -2147483648) return 0;
    return rev;
}
