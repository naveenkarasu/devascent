fn left_rotate(arr: Vec<i64>, k: i64) -> Vec<i64> {
    let n = arr.len();
    if n == 0 {
        return arr;
    }
    let k = (k.rem_euclid(n as i64)) as usize;
    let mut result = Vec::with_capacity(n);
    result.extend_from_slice(&arr[k..]);
    result.extend_from_slice(&arr[..k]);
    result
}
