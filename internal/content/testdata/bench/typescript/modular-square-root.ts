function modular_sqrt(n: number, m: number): number {
    for (let x = 0; x < m; x++) {
        if ((x * x) % m === n) {
            return x;
        }
    }
    return -1;
}
