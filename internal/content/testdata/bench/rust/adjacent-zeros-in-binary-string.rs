fn has_adjacent_zeros(s: String) -> bool {
    let chars: Vec<char> = s.chars().collect();
    for i in 1..chars.len() {
        if chars[i] == '0' && chars[i - 1] == '0' {
            return true;
        }
    }
    false
}
