fn max_area(heights: Vec<i64>) -> i64 {
    let mut i = 0i64;
    let mut j = heights.len() as i64 - 1;
    let mut best = 0i64;
    while i < j {
        let area = heights[i as usize].min(heights[j as usize]) * (j - i);
        best = best.max(area);
        if heights[i as usize] < heights[j as usize] {
            i += 1;
        } else {
            j -= 1;
        }
    }
    best
}
