function max_area(heights) {
    let i = 0, j = heights.length - 1;
    let best = 0;
    while (i < j) {
        const area = Math.min(heights[i], heights[j]) * (j - i);
        best = Math.max(best, area);
        if (heights[i] < heights[j]) i++;
        else j--;
    }
    return best;
}
