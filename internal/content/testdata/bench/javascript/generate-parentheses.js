function generate_parentheses(n) {
  const result = [];
  function backtrack(s, openCount, closeCount) {
    if (s.length === 2 * n) {
      result.push(s);
      return;
    }
    if (openCount < n) {
      backtrack(s + '(', openCount + 1, closeCount);
    }
    if (closeCount < openCount) {
      backtrack(s + ')', openCount, closeCount + 1);
    }
  }
  backtrack('', 0, 0);
  return result.sort();
}
