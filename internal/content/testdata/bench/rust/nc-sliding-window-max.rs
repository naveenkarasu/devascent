use std::collections::VecDeque;

fn max_sliding_window(nums: Vec<i64>, k: i64) -> Vec<i64> {
    let n = nums.len() as i64;
    let k = k;
    let mut result: Vec<i64> = Vec::new();
    let mut dq: VecDeque<usize> = VecDeque::new(); // indices, decreasing by value
    for i in 0..n as usize {
        while let Some(&front) = dq.front() {
            if (front as i64) < i as i64 - k + 1 {
                dq.pop_front();
            } else {
                break;
            }
        }
        while let Some(&back) = dq.back() {
            if nums[back] < nums[i] {
                dq.pop_back();
            } else {
                break;
            }
        }
        dq.push_back(i);
        if i as i64 >= k - 1 {
            result.push(nums[*dq.front().unwrap()]);
        }
    }
    result
}
