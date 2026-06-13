function find_unpaired(arr: number[]): number {
    let result = 0;
    for (const x of arr) {
        result ^= x;
    }
    return result;
}
