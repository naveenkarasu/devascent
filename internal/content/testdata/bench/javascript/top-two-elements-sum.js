function top_two_sum(nums) {
  let first = 0, second = 0;
  for (const x of nums) {
    if (x >= first) {
      second = first;
      first = x;
    } else if (x > second) {
      second = x;
    }
  }
  return first + second;
}
