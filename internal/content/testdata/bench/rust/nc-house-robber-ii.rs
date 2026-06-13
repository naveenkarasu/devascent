fn rob_circular(nums: Vec<i64>) -> i64 {
    if nums.len() == 1 {
        return nums[0];
    }

    fn rob_linear(houses: &[i64]) -> i64 {
        if houses.is_empty() {
            return 0;
        }
        if houses.len() == 1 {
            return houses[0];
        }
        let mut prev2 = houses[0];
        let mut prev1 = houses[0].max(houses[1]);
        for i in 2..houses.len() {
            let curr = prev1.max(prev2 + houses[i]);
            prev2 = prev1;
            prev1 = curr;
        }
        prev1
    }

    let n = nums.len();
    rob_linear(&nums[..n - 1]).max(rob_linear(&nums[1..]))
}
