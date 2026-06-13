fn has_cycle(values: Vec<i64>, pos: i64) -> bool {
    let n = values.len();
    if n == 0 {
        return false;
    }
    // next[i] = index of next node, or -1 for null
    let mut next: Vec<i64> = vec![-1; n];
    for i in 0..n - 1 {
        next[i] = (i + 1) as i64;
    }
    if pos >= 0 {
        next[n - 1] = pos;
    }
    // Floyd's on indices, head = 0
    let mut slow: i64 = 0;
    let mut fast: i64 = 0;
    // fast != null && next[fast] != null
    while fast != -1 && next[fast as usize] != -1 {
        slow = next[slow as usize];
        fast = next[next[fast as usize] as usize];
        if slow == fast {
            return true;
        }
    }
    false
}
