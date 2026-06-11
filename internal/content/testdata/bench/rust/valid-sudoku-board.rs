fn is_valid_sudoku(board: Vec<Vec<String>>) -> bool {
    let mut rows = [[false; 9]; 9];
    let mut cols = [[false; 9]; 9];
    let mut boxes = [[false; 9]; 9];
    for r in 0..9 {
        for c in 0..9 {
            let val = &board[r][c];
            if val == "." {
                continue;
            }
            let d = val.chars().next().unwrap();
            let idx = (d as usize) - ('1' as usize);
            if idx >= 9 {
                continue;
            }
            let box_idx = (r / 3) * 3 + (c / 3);
            if rows[r][idx] || cols[c][idx] || boxes[box_idx][idx] {
                return false;
            }
            rows[r][idx] = true;
            cols[c][idx] = true;
            boxes[box_idx][idx] = true;
        }
    }
    true
}
