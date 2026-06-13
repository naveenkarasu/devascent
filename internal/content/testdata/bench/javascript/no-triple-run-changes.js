function min_changes_no_triple(s) {
  const arr = s.split('');
  let ans = 0;
  let same = 1;
  for (let i = 1; i < arr.length; i++) {
    if (arr[i] === arr[i - 1]) {
      same++;
    } else {
      same = 1;
    }
    if (same === 3) {
      ans++;
      arr[i] = '@';
      same = 1;
    }
  }
  return ans;
}
