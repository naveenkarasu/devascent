fn top_two_sum(nums: Vec<i64>) -> i64 {
    let mut first: i64 = 0;
    let mut second: i64 = 0;
    for x in nums {
        if x >= first {
            second = first;
            first = x;
        } else if x > second {
            second = x;
        }
    }
    first + second
}
