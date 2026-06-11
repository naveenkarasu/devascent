function min_eating_speed(piles, h) {
    let lo = 1, hi = Math.max(...piles);
    let ans = hi;
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        let hours = 0;
        for (const p of piles) hours += Math.ceil(p / mid);
        if (hours <= h) {
            ans = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    return ans;
}
