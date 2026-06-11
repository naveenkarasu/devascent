fn daily_temperatures(temps: Vec<i64>) -> Vec<i64> {
    let mut result = vec![0i64; temps.len()];
    let mut stack: Vec<(i64, usize)> = Vec::new(); // (temperature, index)
    for (i, &t) in temps.iter().enumerate() {
        while let Some(&(top_t, j)) = stack.last() {
            if t > top_t {
                stack.pop();
                result[j] = (i - j) as i64;
            } else {
                break;
            }
        }
        stack.push((t, i));
    }
    result
}
