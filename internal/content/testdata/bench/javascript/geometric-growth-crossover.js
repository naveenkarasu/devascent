function steps_to_overtake(a, b) {
  let steps = 0;
  while (a <= b) {
    a *= 3;
    b *= 2;
    steps++;
  }
  return steps;
}
