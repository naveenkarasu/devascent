function merge_sorted_in_place(arr1, arr2) {
  const m = arr1.length;
  const n = arr2.length;
  const result = new Array(m + n).fill(0);
  let l = m - 1, r = n - 1, k = m + n - 1;
  while (l >= 0 && r >= 0) {
    if (arr1[l] > arr2[r]) {
      result[k] = arr1[l];
      l--;
    } else {
      result[k] = arr2[r];
      r--;
    }
    k--;
  }
  while (l >= 0) {
    result[k] = arr1[l];
    l--;
    k--;
  }
  while (r >= 0) {
    result[k] = arr2[r];
    r--;
    k--;
  }
  return result;
}
