function majority_element(nums) {
    let count = 0;
    let candidate = null;
    for (const n of nums) {
        if (count === 0) candidate = n;
        count += (n === candidate) ? 1 : -1;
    }
    const freq = nums.filter(x => x === candidate).length;
    if (freq > Math.floor(nums.length / 2)) return candidate;
    return -1;
}
