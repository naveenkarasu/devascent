function sort_three_values(arr, low_val = 0, mid_val = 1, high_val = 2) {
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
