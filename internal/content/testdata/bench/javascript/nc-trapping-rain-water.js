function trap(heights) {
    if (!heights || heights.length === 0) return 0;
    let i = 0, j = heights.length - 1;
    let left_max = heights[i], right_max = heights[j];
    let total = 0;
    while (i < j) {
        if (left_max < right_max) {
            i++;
            left_max = Math.max(left_max, heights[i]);
            total += left_max - heights[i];
        } else {
            j--;
            right_max = Math.max(right_max, heights[j]);
            total += right_max - heights[j];
        }
    }
    return total;
}
