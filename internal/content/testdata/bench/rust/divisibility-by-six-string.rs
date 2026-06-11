fn is_divisible_by_six(s: String) -> bool {
    let chars: Vec<char> = s.chars().collect();
    if chars.is_empty() {
        return false;
    }
    let last = chars[chars.len() - 1].to_digit(10).unwrap() as i64;
    if last % 2 != 0 {
        return false;
    }
    let mut digit_sum: i64 = 0;
    for ch in &chars {
        digit_sum += ch.to_digit(10).unwrap() as i64;
    }
    digit_sum % 3 == 0
}
