function find_majority_element(nums: number[]): number {
    let candidate: number = 0;
    let count = 0;
    for (const x of nums) {
        if (count === 0) candidate = x;
        count += x === candidate ? 1 : -1;
    }
    const freq = nums.filter(x => x === candidate).length;
    if (freq > Math.floor(nums.length / 2)) return candidate;
    return -1;
}
