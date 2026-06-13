function left_rotate(arr: number[], k: number): number[] {
    const n = arr.length;
    if (n === 0) return arr;
    k = k % n;
    return arr.slice(k).concat(arr.slice(0, k));
}
