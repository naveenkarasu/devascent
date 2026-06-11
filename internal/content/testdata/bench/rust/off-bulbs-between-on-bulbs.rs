fn off_bulbs_between_lit(bulbs: Vec<i64>) -> i64 {
    let n = bulbs.len();
    let mut first_on: i64 = -1;
    let mut last_on: i64 = -1;
    for i in 0..n {
        if bulbs[i] == 1 {
            first_on = i as i64;
            break;
        }
    }
    for i in (0..n).rev() {
        if bulbs[i] == 1 {
            last_on = i as i64;
            break;
        }
    }
    if first_on == -1 {
        return 0;
    }
    let mut count: i64 = 0;
    for j in first_on..=last_on {
        if bulbs[j as usize] == 0 {
            count += 1;
        }
    }
    count
}
