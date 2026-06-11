function rob_circular(nums: number[]): number {
    if (nums.length === 1) return nums[0];

    function rob_linear(houses: number[]): number {
        if (houses.length === 0) return 0;
        if (houses.length === 1) return houses[0];
        let prev2 = houses[0], prev1 = Math.max(houses[0], houses[1]);
        for (let i = 2; i < houses.length; i++) {
            const curr = Math.max(prev1, prev2 + houses[i]);
            prev2 = prev1;
            prev1 = curr;
        }
        return prev1;
    }

    return Math.max(rob_linear(nums.slice(0, -1)), rob_linear(nums.slice(1)));
}
