fn insertion_sort(arr: Vec<i64>) -> Vec<i64> {
    let mut a = arr.clone();
    let n = a.len();
    for i in 1..n {
        let key = a[i];
        let mut j = i as i64 - 1;
        while j >= 0 && a[j as usize] > key {
            a[(j + 1) as usize] = a[j as usize];
            j -= 1;
        }
        a[(j + 1) as usize] = key;
    }
    a
}
