function permute(nums) {
  const res = [];
  function backtrack(current, remaining) {
    if (remaining.length === 0) {
      res.push(current.slice());
      return;
    }
    for (let i = 0; i < remaining.length; i++) {
      current.push(remaining[i]);
      backtrack(current, remaining.slice(0, i).concat(remaining.slice(i + 1)));
      current.pop();
    }
  }
  backtrack([], nums);
  res.sort((a, b) => {
    for (let i = 0; i < Math.min(a.length, b.length); i++) {
      if (a[i] !== b[i]) return a[i] - b[i];
    }
    return a.length - b.length;
  });
  return res;
}
