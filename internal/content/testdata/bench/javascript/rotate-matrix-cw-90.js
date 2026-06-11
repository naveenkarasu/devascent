function rotate_matrix_cw(matrix) {
  const n = matrix.length;
  const result = [];
  for (let i = 0; i < n; i++) {
    const row = [];
    for (let j = 0; j < n; j++) {
      row.push(matrix[n - 1 - j][i]);
    }
    result.push(row);
  }
  return result;
}
