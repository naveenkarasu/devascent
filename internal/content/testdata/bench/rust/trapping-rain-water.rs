fn trap_water(height: Vec<i64>) -> i64 {
    let mut l = 0i64;
    let mut r = height.len() as i64 - 1;
    let mut level = 0i64;
    let mut water = 0i64;
    while l < r {
        let lower = if height[l as usize] < height[r as usize] {
            height[l as usize]
        } else {
            height[r as usize]
        };
        if height[l as usize] < height[r as usize] {
            l += 1;
        } else {
            r -= 1;
        }
        level = level.max(lower);
        water += level - lower;
    }
    water
}
