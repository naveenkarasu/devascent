function two_sum(nums: number[], target: number): number[] {
    const seen: {[key: number]: number} = {};
    for (let i = 0; i < nums.length; i++) {
        const need = target - nums[i];
        if (need in seen) {
            return [seen[need], i];
        }
        seen[nums[i]] = i;
    }
    return [];
}
