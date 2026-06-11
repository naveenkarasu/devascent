fn car_fleet(target: i64, position: Vec<i64>, speed: Vec<i64>) -> i64 {
    let mut pairs: Vec<(i64, i64)> = position
        .iter()
        .cloned()
        .zip(speed.iter().cloned())
        .collect();
    // sort by position descending (reverse of ascending zip sort)
    pairs.sort_by(|a, b| b.cmp(a));
    let mut stack: Vec<f64> = Vec::new();
    for &(pos, spd) in pairs.iter() {
        let time = (target - pos) as f64 / spd as f64;
        if stack.is_empty() || time > *stack.last().unwrap() {
            stack.push(time);
        }
    }
    stack.len() as i64
}
