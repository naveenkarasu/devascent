fn multiply_strings(num1: String, num2: String) -> String {
    let a: Vec<u8> = num1.bytes().collect();
    let b: Vec<u8> = num2.bytes().collect();
    let m = a.len();
    let n = b.len();
    let mut pos = vec![0i64; m + n];
    for i in (0..m).rev() {
        for j in (0..n).rev() {
            let mul = (a[i] - b'0') as i64 * (b[j] - b'0') as i64;
            let p1 = i + j;
            let p2 = i + j + 1;
            let total = mul + pos[p2];
            pos[p2] = total % 10;
            pos[p1] += total / 10;
        }
    }
    let s: String = pos.iter().map(|d| std::char::from_digit(*d as u32, 10).unwrap()).collect();
    let trimmed = s.trim_start_matches('0');
    if trimmed.is_empty() {
        "0".to_string()
    } else {
        trimmed.to_string()
    }
}
