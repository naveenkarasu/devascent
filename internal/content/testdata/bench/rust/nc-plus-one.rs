fn plus_one(digits: Vec<i64>) -> Vec<i64> {
    let mut digits = digits;
    let n = digits.len();
    for i in (0..n).rev() {
        if digits[i] < 9 {
            digits[i] += 1;
            return digits;
        }
        digits[i] = 0;
    }
    let mut result = vec![1i64];
    result.extend(digits);
    result
}
