function is_divisible_by_six(s) {
  if (parseInt(s[s.length - 1]) % 2 !== 0) return false;
  let digitSum = 0;
  for (const ch of s) digitSum += parseInt(ch);
  return digitSum % 3 === 0;
}
