function is_anagram(s, t) {
  if (s.length !== t.length) return false;
  const counts = {};
  for (const c of s) {
    counts[c] = (counts[c] || 0) + 1;
  }
  for (const c of t) {
    if (!counts[c]) return false;
    counts[c]--;
  }
  return true;
}
