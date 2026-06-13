function largest_rectangle_area(heights: number[]): number {
    let maxArea = 0;
    const stack: [number, number][] = []; // [start_index, height]
    for (let i = 0; i < heights.length; i++) {
        const h = heights[i];
        let start = i;
        while (stack.length > 0 && stack[stack.length - 1][1] > h) {
            const [idx, height] = stack.pop()!;
            maxArea = Math.max(maxArea, height * (i - idx));
            start = idx;
        }
        stack.push([start, h]);
    }
    for (const [idx, height] of stack) {
        maxArea = Math.max(maxArea, height * (heights.length - idx));
    }
    return maxArea;
}
