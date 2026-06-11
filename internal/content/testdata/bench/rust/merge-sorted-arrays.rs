fn merge_sorted_arrays(arr1: Vec<i64>, arr2: Vec<i64>) -> Vec<i64> {
    let mut result: Vec<i64> = Vec::new();
    let mut i = 0usize;
    let mut j = 0usize;
    while i < arr1.len() && j < arr2.len() {
        if arr1[i] <= arr2[j] {
            result.push(arr1[i]);
            i += 1;
        } else {
            result.push(arr2[j]);
            j += 1;
        }
    }
    while i < arr1.len() {
        result.push(arr1[i]);
        i += 1;
    }
    while j < arr2.len() {
        result.push(arr2[j]);
        j += 1;
    }
    result
}
