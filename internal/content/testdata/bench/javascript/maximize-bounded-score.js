function maximize_score(arr, k) {
  let maxScore = 0;
  for (const x of arr) {
    const score = Math.floor((2 * k) / (1 + 2 * Math.abs(x - k)));
    if (score > maxScore) maxScore = score;
  }
  return maxScore;
}
