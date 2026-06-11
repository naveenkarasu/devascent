function toggle_case(s: string): string {
    let result: string[] = [];
    for (const ch of s) {
        if (ch === ch.toUpperCase() && ch !== ch.toLowerCase()) {
            result.push(ch.toLowerCase());
        } else if (ch === ch.toLowerCase() && ch !== ch.toUpperCase()) {
            result.push(ch.toUpperCase());
        } else {
            result.push(ch);
        }
    }
    return result.join('');
}
