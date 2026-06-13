fn compressed_length(s: String) -> i64 {
    let chars: Vec<char> = s.chars().collect();
    let n = chars.len();
    let mut count = 0i64;
    let mut i = 0usize;
    while i < n {
        if i + 1 < n && chars[i] == chars[i + 1] {
            count += 1;
            i += 2;
        } else {
            count += 1;
            i += 1;
        }
    }
    count
}
