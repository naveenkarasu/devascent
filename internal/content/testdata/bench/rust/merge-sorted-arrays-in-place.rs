fn merge_sorted_in_place(arr1: Vec<i64>, arr2: Vec<i64>) -> Vec<i64> {
    let m = arr1.len() as i64;
    let n = arr2.len() as i64;
    let mut result: Vec<i64> = vec![0; (m + n) as usize];
    let mut l = m - 1;
    let mut r = n - 1;
    let mut k = m + n - 1;
    while l >= 0 && r >= 0 {
        if arr1[l as usize] > arr2[r as usize] {
            result[k as usize] = arr1[l as usize];
            l -= 1;
        } else {
            result[k as usize] = arr2[r as usize];
            r -= 1;
        }
        k -= 1;
    }
    while l >= 0 {
        result[k as usize] = arr1[l as usize];
        l -= 1;
        k -= 1;
    }
    while r >= 0 {
        result[k as usize] = arr2[r as usize];
        r -= 1;
        k -= 1;
    }
    result
}
