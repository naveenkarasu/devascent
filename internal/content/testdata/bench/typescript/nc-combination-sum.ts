function combination_sum(candidates: number[], target: number): number[][] {
    const res: number[][] = [];
    candidates.sort((a, b) => a - b);
    function backtrack(start: number, current: number[], remaining: number): void {
        if (remaining === 0) {
            res.push(current.slice());
            return;
        }
        for (let i = start; i < candidates.length; i++) {
            if (candidates[i] > remaining) break;
            current.push(candidates[i]);
            backtrack(i, current, remaining - candidates[i]);
            current.pop();
        }
    }
    backtrack(0, [], target);
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
