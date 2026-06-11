function find_single(arr) {
    let result = 0;
    for (let x of arr) {
        result ^= x;
    }
    return result;
}
