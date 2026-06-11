function subsets_with_dup(nums) {
    const res = [];
    nums.sort((a, b) => a - b);
    function backtrack(start, current) {
        res.push([...current]);
        for (let i = start; i < nums.length; i++) {
            if (i > start && nums[i] === nums[i - 1]) continue;
            current.push(nums[i]);
            backtrack(i + 1, current);
            current.pop();
        }
    }
    backtrack(0, []);
    for (let i = 0; i < res.length; i++) res[i].sort((a, b) => a - b);
    res.sort((a, b) => {
        for (let i = 0; i < Math.min(a.length, b.length); i++) {
            if (a[i] !== b[i]) return a[i] - b[i];
        }
        return a.length - b.length;
    });
    return res;
}
