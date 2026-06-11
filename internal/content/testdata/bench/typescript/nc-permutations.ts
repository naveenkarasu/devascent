function permute(nums: number[]): number[][] {
    const res: number[][] = [];
    function backtrack(current: number[], remaining: number[]): void {
        if (remaining.length === 0) {
            res.push(current.slice());
            return;
        }
        for (let i = 0; i < remaining.length; i++) {
            current.push(remaining[i]);
            backtrack(current, remaining.slice(0, i).concat(remaining.slice(i + 1)));
            current.pop();
        }
    }
    backtrack([], nums);
    res.sort((a, b) => {
        for (let i = 0; i < a.length; i++) {
            if (a[i] !== b[i]) return a[i] - b[i];
        }
        return 0;
    });
    return res;
}
