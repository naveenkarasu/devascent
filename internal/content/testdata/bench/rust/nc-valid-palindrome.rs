fn is_palindrome(s: String) -> bool {
    let chars: Vec<char> = s.chars().collect();
    let mut i = 0i64;
    let mut j = chars.len() as i64 - 1;
    while i < j {
        while i < j && !chars[i as usize].is_alphanumeric() {
            i += 1;
        }
        while i < j && !chars[j as usize].is_alphanumeric() {
            j -= 1;
        }
        if chars[i as usize].to_ascii_lowercase() != chars[j as usize].to_ascii_lowercase() {
            return false;
        }
        i += 1;
        j -= 1;
    }
    true
}
