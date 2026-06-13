function last_remaining(n: number): number {
    if (n === 1) return 1;
    return 2 * (1 + Math.floor(n / 2) - last_remaining(Math.floor(n / 2)));
}
