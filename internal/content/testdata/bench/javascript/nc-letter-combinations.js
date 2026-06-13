function letter_combinations(digits) {
  if (!digits) return [];
  const mapping = {
    '2': 'abc', '3': 'def', '4': 'ghi', '5': 'jkl',
    '6': 'mno', '7': 'pqrs', '8': 'tuv', '9': 'wxyz'
  };
  const res = [];
  function backtrack(index, current) {
    if (index === digits.length) {
      res.push(current);
      return;
    }
    for (const ch of mapping[digits[index]]) {
      backtrack(index + 1, current + ch);
    }
  }
  backtrack(0, '');
  res.sort();
  return res;
}
