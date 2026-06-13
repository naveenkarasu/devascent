function combination_sum2(candidates, target) {
    const res = [];
    candidates.sort((a, b) => a - b);
    function backtrack(start, current, remaining) {
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
    for (let i = 0; i < res.length; i++) res[i].sort((a, b) => a - b);
    res.sort((a, b) => {
        for (let i = 0; i < Math.min(a.length, b.length); i++) {
            if (a[i] !== b[i]) return a[i] - b[i];
        }
        return a.length - b.length;
    });
    return res;
}
