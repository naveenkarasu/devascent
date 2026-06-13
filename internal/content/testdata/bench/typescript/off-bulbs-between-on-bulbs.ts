function off_bulbs_between_lit(bulbs: number[]): number {
    const n = bulbs.length;
    let first_on = -1;
    let last_on = -1;
    for (let i = 0; i < n; i++) {
        if (bulbs[i] === 1) {
            first_on = i;
            break;
        }
    }
    for (let i = n - 1; i >= 0; i--) {
        if (bulbs[i] === 1) {
            last_on = i;
            break;
        }
    }
    if (first_on === -1) return 0;
    let count = 0;
    for (let j = first_on; j <= last_on; j++) {
        if (bulbs[j] === 0) count++;
    }
    return count;
}
