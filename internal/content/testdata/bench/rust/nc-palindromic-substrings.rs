fn count_substrings(s: String) -> i64 {
    let b = s.as_bytes();
    let n = b.len();
    let mut count: i64 = 0;

    let mut expand = |mut left: i64, mut right: i64| {
        while left >= 0 && (right as usize) < n && b[left as usize] == b[right as usize] {
            count += 1;
            left -= 1;
            right += 1;
        }
    };

    for i in 0..n {
        expand(i as i64, i as i64);
        expand(i as i64, i as i64 + 1);
    }

    count
}
