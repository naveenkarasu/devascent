fn edit_distance(str1: String, str2: String) -> i64 {
    let a: Vec<char> = str1.chars().collect();
    let b: Vec<char> = str2.chars().collect();
    let m = a.len();
    let n = b.len();
    let mut cur: Vec<i64> = (0..=n as i64).collect();
    for i in 1..=m {
        let mut pre = cur[0];
        cur[0] = i as i64;
        for j in 1..=n {
            let temp = cur[j];
            if a[i - 1] == b[j - 1] {
                cur[j] = pre;
            } else {
                let mn = pre.min(cur[j - 1]).min(cur[j]);
                cur[j] = mn + 1;
            }
            pre = temp;
        }
    }
    cur[n]
}
