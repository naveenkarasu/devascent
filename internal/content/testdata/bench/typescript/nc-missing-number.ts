function missing_number(nums: number[]): number {
    const n = nums.length;
    const expected = (n * (n + 1)) / 2;
    return expected - nums.reduce((sum, x) => sum + x, 0);
}
