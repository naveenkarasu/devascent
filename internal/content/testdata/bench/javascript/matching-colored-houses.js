function match_colored_houses(left_colors, right_colors) {
  const colorToPos = {};
  for (let i = 0; i < right_colors.length; i++) {
    colorToPos[right_colors[i]] = i + 1;
  }
  return left_colors.map(c => colorToPos[c]);
}
