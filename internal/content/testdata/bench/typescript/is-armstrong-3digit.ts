function is_armstrong_3digit(n: number): boolean {
    if (n < 100 || n > 999) return false;
    let s = 0;
    let temp = n;
    while (temp > 0) {
        const r = temp % 10;
        s += r * r * r;
        temp = Math.floor(temp / 10);
    }
    return s === n;
}
