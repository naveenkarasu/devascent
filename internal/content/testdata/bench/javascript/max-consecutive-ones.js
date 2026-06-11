function max_consecutive_ones(nums) {
    let maxRun = 0;
    let current = 0;
    for (let x of nums) {
        if (x === 1) {
            current++;
            if (current > maxRun) maxRun = current;
        } else {
            current = 0;
        }
    }
    return maxRun;
}
