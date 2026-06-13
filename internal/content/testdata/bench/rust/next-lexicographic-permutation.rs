fn next_permutation(mut nums: Vec<i64>) -> Vec<i64> {
    let n = nums.len();
    if n == 0 {
        return nums;
    }
    let mut pivot: i64 = -1;
    for i in (0..n - 1).rev() {
        if nums[i] < nums[i + 1] {
            pivot = i as i64;
            break;
        }
    }
    if pivot == -1 {
        nums.reverse();
        return nums;
    }
    let p = pivot as usize;
    for r in (p + 1..n).rev() {
        if nums[r] > nums[p] {
            nums.swap(p, r);
            break;
        }
    }
    nums[p + 1..].reverse();
    nums
}
