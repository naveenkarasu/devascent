fn longest_palindrome(s: String) -> String {
    let b = s.as_bytes();
    let n = b.len();
    if n == 0 {
        return String::new();
    }
    let mut res_start: usize = 0;
    let mut res_len: usize = 1;

    let expand = |mut left: i64, mut right: i64, res_start: &mut usize, res_len: &mut usize| {
        while left >= 0 && (right as usize) < n && b[left as usize] == b[right as usize] {
            let length = (right - left + 1) as usize;
            if length > *res_len {
                *res_start = left as usize;
                *res_len = length;
            }
            left -= 1;
            right += 1;
        }
    };

    for i in 0..n {
        expand(i as i64, i as i64, &mut res_start, &mut res_len);
        expand(i as i64, i as i64 + 1, &mut res_start, &mut res_len);
    }

    String::from_utf8(b[res_start..res_start + res_len].to_vec()).unwrap()
}
