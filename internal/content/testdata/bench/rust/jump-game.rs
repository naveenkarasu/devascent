fn can_jump(nums: Vec<i64>) -> bool {
    let mut reach: i64 = 0;
    for (i, &n) in nums.iter().enumerate() {
        let i = i as i64;
        if i > reach {
            return false;
        }
        if i + n > reach {
            reach = i + n;
        }
    }
    true
}
