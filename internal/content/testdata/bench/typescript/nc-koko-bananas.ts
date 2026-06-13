function min_eating_speed(piles: number[], h: number): number {
    let lo = 1, hi = Math.max(...piles);
    let ans = hi;
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        const hours = piles.reduce((sum, p) => sum + Math.ceil(p / mid), 0);
        if (hours <= h) {
            ans = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    return ans;
}
