function sort_three_values(arr: number[], low_val: number = 0, mid_val: number = 1, high_val: number = 2): number[] {
    let lo = 0, mid = 0, hi = arr.length - 1;
    while (mid <= hi) {
        if (arr[mid] === low_val) {
            [arr[lo], arr[mid]] = [arr[mid], arr[lo]];
            lo++;
            mid++;
        } else if (arr[mid] === mid_val) {
            mid++;
        } else {
            [arr[mid], arr[hi]] = [arr[hi], arr[mid]];
            hi--;
        }
    }
    return arr;
}
