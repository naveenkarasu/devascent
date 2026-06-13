function max_product(nums: number[]): number {
    let max_prod = nums[0];
    let min_prod = nums[0];
    let res = nums[0];
    for (let i = 1; i < nums.length; i++) {
        const n = nums[i];
        const candidates = [n, max_prod * n, min_prod * n];
        max_prod = Math.max(...candidates);
        min_prod = Math.min(...candidates);
        res = Math.max(res, max_prod);
    }
    return res;
}
