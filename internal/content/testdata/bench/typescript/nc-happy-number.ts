function is_happy(n: number): boolean {
    const seen = new Set<number>();
    while (n !== 1) {
        if (seen.has(n)) return false;
        seen.add(n);
        let total = 0;
        while (n > 0) {
            const digit = n % 10;
            total += digit * digit;
            n = Math.floor(n / 10);
        }
        n = total;
    }
    return true;
}
