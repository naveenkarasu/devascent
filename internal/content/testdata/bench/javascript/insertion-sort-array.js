function insertion_sort(arr) {
    let a = arr.slice();
    let n = a.length;
    for (let i = 1; i < n; i++) {
        let key = a[i];
        let j = i - 1;
        while (j >= 0 && a[j] > key) {
            a[j + 1] = a[j];
            j--;
        }
        a[j + 1] = key;
    }
    return a;
}
