function single_number(nums) {
    let result = 0;
    for (const n of nums) result ^= n;
    return result;
}
