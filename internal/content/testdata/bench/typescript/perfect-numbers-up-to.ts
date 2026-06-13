function perfect_numbers_up_to(n: number): number[] {
    const result: number[] = [];
    for (let i = 2; i <= n; i++) {
        let s = 1;
        for (let d = 2; d * d <= i; d++) {
            if (i % d === 0) {
                s += d + Math.floor(i / d);
                if (d * d === i) s -= d;
            }
        }
        if (s === i) result.push(i);
    }
    return result;
}
