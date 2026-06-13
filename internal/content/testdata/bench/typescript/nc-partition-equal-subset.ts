function can_partition(nums: number[]): boolean {
    const total = nums.reduce((a, b) => a + b, 0);
    if (total % 2 !== 0) return false;
    const target = total / 2;
    let dp = new Set<number>([0]);
    for (const n of nums) {
        const snapshot = [...dp];
        for (const s of snapshot) {
            dp.add(s + n);
        }
        if (dp.has(target)) return true;
    }
    return false;
}
