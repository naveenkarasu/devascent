function contains_element(arr, key) {
    let lo = 0, hi = arr.length - 1;
    while (lo <= hi) {
        let mid = Math.floor((lo + hi) / 2);
        if (arr[mid] === key) return true;
        else if (arr[mid] < key) lo = mid + 1;
        else hi = mid - 1;
    }
    return false;
}
