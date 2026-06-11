fn bubble_sort_info(arr: Vec<i64>) -> J {
    let mut a = arr.clone();
    let mut swaps: i64 = 0;
    let mut is_sorted = false;
    while !is_sorted {
        is_sorted = true;
        for i in 0..a.len().saturating_sub(1) {
            if a[i] > a[i + 1] {
                a.swap(i, i + 1);
                swaps += 1;
                is_sorted = false;
            }
        }
    }
    let first = a[0];
    let last = a[a.len() - 1];
    J::Obj(vec![
        ("first".to_string(), J::Int(first)),
        ("last".to_string(), J::Int(last)),
        ("swaps".to_string(), J::Int(swaps)),
    ])
}
