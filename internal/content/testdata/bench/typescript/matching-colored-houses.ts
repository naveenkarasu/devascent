function match_colored_houses(left_colors: string[], right_colors: string[]): number[] {
    const color_to_pos: {[key: string]: number} = {};
    for (let i = 0; i < right_colors.length; i++) {
        color_to_pos[right_colors[i]] = i + 1;
    }
    return left_colors.map(c => color_to_pos[c]);
}
