function peak_occupancy(initial, deltas) {
  let current = initial;
  let peak = initial;
  for (const d of deltas) {
    current += d;
    if (current > peak) peak = current;
  }
  return peak;
}
