function subsets_with_dup(nums: number[]): number[][] {
    function cmpArr(a: number[], b: number[]): number {
        const n = Math.min(a.length, b.length);
        for (let i = 0; i < n; i++) {
            if (a[i] < b[i]) return -1;
            if (a[i] > b[i]) return 1;
        }
        return a.length - b.length;
    }

    const res: number[][] = [];
    nums.sort((a, b) => a - b);

    function backtrack(start: number, current: number[]): void {
        res.push([...current]);
        for (let i = start; i < nums.length; i++) {
            if (i > start && nums[i] === nums[i - 1]) continue;
            current.push(nums[i]);
            backtrack(i + 1, current);
            current.pop();
        }
    }

    backtrack(0, []);
    const sorted = res.map(s => [...s].sort((a, b) => a - b));
    sorted.sort(cmpArr);
    return sorted;
}
