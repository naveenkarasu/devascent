fn two_sum_ii(numbers: Vec<i64>, target: i64) -> Vec<i64> {
    let mut i = 0i64;
    let mut j = numbers.len() as i64 - 1;
    while i < j {
        let total = numbers[i as usize] + numbers[j as usize];
        if total == target {
            return vec![i + 1, j + 1];
        }
        if total < target {
            i += 1;
        } else {
            j -= 1;
        }
    }
    Vec::new()
}
