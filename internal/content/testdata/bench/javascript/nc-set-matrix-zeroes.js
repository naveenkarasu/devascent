function set_zeroes(matrix) {
    const zeroRows = new Set();
    const zeroCols = new Set();
    const rows = matrix.length;
    const cols = matrix[0].length;
    for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
            if (matrix[r][c] === 0) {
                zeroRows.add(r);
                zeroCols.add(c);
            }
        }
    }
    for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
            if (zeroRows.has(r) || zeroCols.has(c)) {
                matrix[r][c] = 0;
            }
        }
    }
    return matrix;
}
