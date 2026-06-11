function is_factorion(n: number): boolean {
    const original = n;
    const factorials = [1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880];
    let s = 0;
    let temp = n;
    while (temp > 0) {
        const d = temp % 10;
        s += factorials[d];
        temp = Math.floor(temp / 10);
    }
    return s === original;
}
