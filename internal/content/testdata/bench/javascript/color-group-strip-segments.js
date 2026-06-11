function count_segments(groups, strip) {
  const colorToGroup = {};
  for (let idx = 0; idx < groups.length; idx++) {
    for (const color of groups[idx]) {
      colorToGroup[color] = idx;
    }
  }
  if (!strip || strip.length === 0) return 0;
  let segments = 1;
  for (let i = 1; i < strip.length; i++) {
    if (colorToGroup[strip[i - 1]] !== colorToGroup[strip[i]]) {
      segments++;
    }
  }
  return segments;
}
