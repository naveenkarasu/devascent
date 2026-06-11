function single_number(nums: number[]): number {
    let result = 0;
    for (const n of nums) result ^= n;
    return result;
}
