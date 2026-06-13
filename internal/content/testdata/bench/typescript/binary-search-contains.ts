function contains_element(arr: number[], key: number): boolean {
    let lo = 0;
    let hi = arr.length - 1;
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        if (arr[mid] === key) return true;
        else if (arr[mid] < key) lo = mid + 1;
        else hi = mid - 1;
    }
    return false;
}
