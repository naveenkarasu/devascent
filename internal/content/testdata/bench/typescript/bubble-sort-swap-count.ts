function bubble_sort_info(arr: number[]): {swaps: number, first: number, last: number} {
    const a = arr.slice();
    let swaps = 0;
    let is_sorted = false;
    while (!is_sorted) {
        is_sorted = true;
        for (let i = 0; i < a.length - 1; i++) {
            if (a[i] > a[i + 1]) {
                const tmp = a[i];
                a[i] = a[i + 1];
                a[i + 1] = tmp;
                swaps++;
                is_sorted = false;
            }
        }
    }
    return {swaps, first: a[0], last: a[a.length - 1]};
}
