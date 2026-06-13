function subsets(nums: number[]): number[][] {
    const res: number[][] = [];
    function backtrack(start: number, current: number[]): void {
        res.push(current.slice());
        for (let i = start; i < nums.length; i++) {
            current.push(nums[i]);
            backtrack(i + 1, current);
            current.pop();
        }
    }
    backtrack(0, []);
    for (let i = 0; i < res.length; i++) res[i].sort((a, b) => a - b);
    res.sort((a, b) => {
        const minLen = Math.min(a.length, b.length);
        for (let i = 0; i < minLen; i++) {
            if (a[i] !== b[i]) return a[i] - b[i];
        }
        return a.length - b.length;
    });
    return res;
}
