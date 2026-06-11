fn largest_rectangle_area(heights: Vec<i64>) -> i64 {
    let mut max_area = 0i64;
    let mut stack: Vec<(usize, i64)> = Vec::new(); // (start_index, height)
    for (i, &h) in heights.iter().enumerate() {
        let mut start = i;
        while let Some(&(idx, height)) = stack.last() {
            if height > h {
                stack.pop();
                let area = height * (i - idx) as i64;
                if area > max_area {
                    max_area = area;
                }
                start = idx;
            } else {
                break;
            }
        }
        stack.push((start, h));
    }
    let n = heights.len();
    for &(idx, height) in stack.iter() {
        let area = height * (n - idx) as i64;
        if area > max_area {
            max_area = area;
        }
    }
    max_area
}
