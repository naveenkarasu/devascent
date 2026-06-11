use std::collections::HashSet;

fn is_valid_sudoku(board: Vec<Vec<String>>) -> bool {
    let mut rows: Vec<HashSet<String>> = (0..9).map(|_| HashSet::new()).collect();
    let mut cols: Vec<HashSet<String>> = (0..9).map(|_| HashSet::new()).collect();
    let mut boxes: Vec<HashSet<String>> = (0..9).map(|_| HashSet::new()).collect();
    for r in 0..9 {
        for c in 0..9 {
            let val = &board[r][c];
            if val == "." {
                continue;
            }
            let box_idx = (r / 3) * 3 + (c / 3);
            if rows[r].contains(val) || cols[c].contains(val) || boxes[box_idx].contains(val) {
                return false;
            }
            rows[r].insert(val.clone());
            cols[c].insert(val.clone());
            boxes[box_idx].insert(val.clone());
        }
    }
    true
}
