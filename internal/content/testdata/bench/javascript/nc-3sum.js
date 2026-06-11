function three_sum(nums) {
    nums.sort((a, b) => a - b);
    const res = [];
    const n = nums.length;
    for (let i = 0; i < n; i++) {
        if (i > 0 && nums[i] === nums[i - 1]) continue;
        let lo = i + 1, hi = n - 1;
        while (lo < hi) {
            const total = nums[i] + nums[lo] + nums[hi];
            if (total < 0) lo++;
            else if (total > 0) hi--;
            else {
                res.push([nums[i], nums[lo], nums[hi]]);
                lo++;
                hi--;
                while (lo < hi && nums[lo] === nums[lo - 1]) lo++;
                while (lo < hi && nums[hi] === nums[hi + 1]) hi--;
            }
        }
    }
    return res;
}
