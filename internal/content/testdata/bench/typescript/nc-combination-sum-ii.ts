function combination_sum2(candidates: number[], target: number): number[][] {
    function cmpArr(a: number[], b: number[]): number {
        const n = Math.min(a.length, b.length);
        for (let i = 0; i < n; i++) {
            if (a[i] < b[i]) return -1;
            if (a[i] > b[i]) return 1;
        }
        return a.length - b.length;
    }

    const res: number[][] = [];
    candidates.sort((a, b) => a - b);

    function backtrack(start: number, current: number[], remaining: number): void {
        if (remaining === 0) {
            res.push([...current]);
            return;
        }
        for (let i = start; i < candidates.length; i++) {
            if (candidates[i] > remaining) break;
            if (i > start && candidates[i] === candidates[i - 1]) continue;
            current.push(candidates[i]);
            backtrack(i + 1, current, remaining - candidates[i]);
            current.pop();
        }
    }

    backtrack(0, [], target);
    const sorted = res.map(c => [...c].sort((a, b) => a - b));
    sorted.sort(cmpArr);
    return sorted;
}
