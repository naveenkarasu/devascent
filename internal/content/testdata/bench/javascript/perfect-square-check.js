function is_perfect_square(num) {
    if (num < 1) return false;
    let lo = 1, hi = num;
    while (lo <= hi) {
        let mid = Math.floor((lo + hi) / 2);
        let sq = mid * mid;
        if (sq === num) return true;
        else if (sq < num) lo = mid + 1;
        else hi = mid - 1;
    }
    return false;
}
