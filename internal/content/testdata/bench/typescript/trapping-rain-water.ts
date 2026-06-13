function trap_water(height: number[]): number {
    let l = 0, r = height.length - 1, level = 0, water = 0;
    while (l < r) {
        const lower = height[l] < height[r] ? height[l] : height[r];
        if (height[l] < height[r]) l++;
        else r--;
        level = Math.max(level, lower);
        water += level - lower;
    }
    return water;
}
