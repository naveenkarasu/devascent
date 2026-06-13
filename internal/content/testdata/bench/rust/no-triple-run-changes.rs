fn min_changes_no_triple(s: String) -> i64 {
    let mut chars: Vec<char> = s.chars().collect();
    let mut ans: i64 = 0;
    let mut same: i64 = 1;
    for i in 1..chars.len() {
        if chars[i] == chars[i - 1] {
            same += 1;
        } else {
            same = 1;
        }
        if same == 3 {
            ans += 1;
            chars[i] = '@';
            same = 1;
        }
    }
    ans
}
