function all_permutations(nums: number[]): number[][] {
    if (nums.length === 0) return [[]];
    const result: number[][] = [];
    for (let i = 0; i < nums.length; i++) {
        const n = nums[i];
        const rest = nums.slice(0, i).concat(nums.slice(i + 1));
        for (const perm of all_permutations(rest)) {
            result.push([n, ...perm]);
        }
    }
    return result;
}
