function pascal_row(n) {
    let row = [];
    for (let k = 0; k < n; k++) {
        // C(n-1, k) = (n-1)! / (k! * (n-1-k)!)
        let val = 1;
        for (let i = 0; i < k; i++) {
            val = Math.floor(val * (n - 1 - i) / (i + 1));
        }
        row.push(val);
    }
    return row;
}
