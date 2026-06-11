function find_celebrity(knows_matrix) {
    const n = knows_matrix.length;
    let candidate = 0;
    for (let i = 1; i < n; i++) {
        if (knows_matrix[candidate][i]) {
            candidate = i;
        }
    }
    for (let i = 0; i < n; i++) {
        if (i !== candidate) {
            if (knows_matrix[candidate][i] || !knows_matrix[i][candidate]) {
                return -1;
            }
        }
    }
    return candidate;
}
