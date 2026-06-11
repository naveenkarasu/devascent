function set_union(a, b) {
  const seen = new Set([...a, ...b]);
  return Array.from(seen).sort((x, y) => x - y);
}
