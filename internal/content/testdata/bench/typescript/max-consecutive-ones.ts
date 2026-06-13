function max_consecutive_ones(nums: number[]): number {
    let max_run = 0;
    let current = 0;
    for (const x of nums) {
        if (x === 1) {
            current++;
            if (current > max_run) max_run = current;
        } else {
            current = 0;
        }
    }
    return max_run;
}
