function is_perfect_square(num: number): boolean {
    if (num < 1) return false;
    let lo = 1;
    let hi = num;
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        const sq = mid * mid;
        if (sq === num) return true;
        else if (sq < num) lo = mid + 1;
        else hi = mid - 1;
    }
    return false;
}
