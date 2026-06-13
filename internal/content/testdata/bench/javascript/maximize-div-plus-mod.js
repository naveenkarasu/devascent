function max_div_plus_mod(l, r, a) {
  let best = Math.floor(r / a) + (r % a);
  const multiple = Math.floor(r / a) * a;
  const candidate = multiple - 1;
  if (l <= candidate && candidate <= r) {
    const val = Math.floor(candidate / a) + (candidate % a);
    if (val > best) best = val;
  }
  return best;
}
