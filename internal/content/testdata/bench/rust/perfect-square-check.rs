fn is_perfect_square(num: i64) -> bool {
    if num < 1 {
        return false;
    }
    let mut lo: i64 = 1;
    let mut hi: i64 = num;
    while lo <= hi {
        let mid = lo + (hi - lo) / 2;
        // Compare mid*mid to num without overflow.
        if mid > num / mid {
            hi = mid - 1;
        } else if mid * mid == num {
            return true;
        } else {
            lo = mid + 1;
        }
    }
    false
}
