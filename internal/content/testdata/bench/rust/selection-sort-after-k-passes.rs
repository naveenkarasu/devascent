fn selection_sort_k_passes(arr: Vec<i64>, k: i64) -> Vec<i64> {
    let mut result = arr.clone();
    let n = result.len();
    let passes = k.min(n as i64 - 1);
    if passes <= 0 {
        return result;
    }
    let passes = passes as usize;
    for i in 0..passes {
        let mut min_idx = i;
        for j in (i + 1)..n {
            if result[j] < result[min_idx] {
                min_idx = j;
            }
        }
        result.swap(i, min_idx);
    }
    result
}
