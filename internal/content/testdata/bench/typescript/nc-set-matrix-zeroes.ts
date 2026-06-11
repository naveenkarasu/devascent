function set_zeroes(matrix: number[][]): number[][] {
    const zeroRows = new Set<number>();
    const zeroCols = new Set<number>();
    for (let r = 0; r < matrix.length; r++) {
        for (let c = 0; c < matrix[0].length; c++) {
            if (matrix[r][c] === 0) {
                zeroRows.add(r);
                zeroCols.add(c);
            }
        }
    }
    for (let r = 0; r < matrix.length; r++) {
        for (let c = 0; c < matrix[0].length; c++) {
            if (zeroRows.has(r) || zeroCols.has(c)) {
                matrix[r][c] = 0;
            }
        }
    }
    return matrix;
}
