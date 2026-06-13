function bubble_sort_info(arr) {
    let a = arr.slice();
    let swaps = 0;
    let isSorted = false;
    while (!isSorted) {
        isSorted = true;
        for (let i = 0; i < a.length - 1; i++) {
            if (a[i] > a[i + 1]) {
                let tmp = a[i];
                a[i] = a[i + 1];
                a[i + 1] = tmp;
                swaps++;
                isSorted = false;
            }
        }
    }
    return {swaps: swaps, first: a[0], last: a[a.length - 1]};
}
