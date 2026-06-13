function longest_ones_run(n) {
    let binary = n.toString(2);
    let runs = binary.split('0');
    let max = 0;
    for (let r of runs) {
        if (r.length > max) max = r.length;
    }
    return max;
}
