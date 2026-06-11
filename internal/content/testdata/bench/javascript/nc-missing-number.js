function missing_number(nums) {
    const n = nums.length;
    const expected = n * (n + 1) / 2;
    return expected - nums.reduce((a, b) => a + b, 0);
}
