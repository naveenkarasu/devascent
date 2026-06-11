function max_and_below_k(n: number, k: number): number {
    if (((k - 1) | k) <= n) {
        return k - 1;
    }
    return k - 2;
}
