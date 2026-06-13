fn contains_element(arr: Vec<i64>, key: i64) -> bool {
    if arr.is_empty() {
        return false;
    }
    let mut lo: i64 = 0;
    let mut hi: i64 = arr.len() as i64 - 1;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        let v = arr[mid as usize];
        if v == key {
            return true;
        } else if v < key {
            lo = mid + 1;
        } else {
            hi = mid - 1;
        }
    }
    false
}
