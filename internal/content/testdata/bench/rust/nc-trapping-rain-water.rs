fn trap(heights: Vec<i64>) -> i64 {
    if heights.is_empty() {
        return 0;
    }
    let mut i = 0i64;
    let mut j = heights.len() as i64 - 1;
    let mut left_max = heights[i as usize];
    let mut right_max = heights[j as usize];
    let mut total = 0i64;
    while i < j {
        if left_max < right_max {
            i += 1;
            left_max = left_max.max(heights[i as usize]);
            total += left_max - heights[i as usize];
        } else {
            j -= 1;
            right_max = right_max.max(heights[j as usize]);
            total += right_max - heights[j as usize];
        }
    }
    total
}
