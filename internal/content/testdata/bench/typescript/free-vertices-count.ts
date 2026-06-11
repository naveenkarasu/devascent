function free_vertices(n: number, m: number): number {
    let k = 0;
    while (k * (k - 1) / 2 < m) {
        k++;
    }
    return n - k;
}
