function majority_element(nums: number[]): number {
    let count = 0;
    let candidate: number | null = null;
    for (const n of nums) {
        if (count === 0) candidate = n;
        count += (n === candidate) ? 1 : -1;
    }
    const cnt = nums.filter(x => x === candidate).length;
    if (cnt > Math.floor(nums.length / 2)) return candidate as number;
    return -1;
}
