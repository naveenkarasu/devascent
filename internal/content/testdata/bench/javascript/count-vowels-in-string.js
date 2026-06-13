function count_vowels(s) {
  const vowels = new Set('aeiouAEIOU');
  let count = 0;
  for (const ch of s) {
    if (vowels.has(ch)) count++;
  }
  return count;
}
