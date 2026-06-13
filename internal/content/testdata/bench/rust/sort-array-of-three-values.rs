fn sort_three_values(mut arr: Vec<i64>, low_val: i64, mid_val: i64, high_val: i64) -> Vec<i64> {
    let _ = high_val;
    let mut lo = 0usize;
    let mut mid = 0usize;
    let mut hi = arr.len() as i64 - 1;
    while mid as i64 <= hi {
        if arr[mid] == low_val {
            arr.swap(lo, mid);
            lo += 1;
            mid += 1;
        } else if arr[mid] == mid_val {
            mid += 1;
        } else {
            arr.swap(mid, hi as usize);
            hi -= 1;
        }
    }
    arr
}
