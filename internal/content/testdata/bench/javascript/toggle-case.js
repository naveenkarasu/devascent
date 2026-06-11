function toggle_case(s) {
    let result = [];
    for (let ch of s) {
        if (ch >= 'A' && ch <= 'Z') {
            result.push(ch.toLowerCase());
        } else if (ch >= 'a' && ch <= 'z') {
            result.push(ch.toUpperCase());
        } else {
            result.push(ch);
        }
    }
    return result.join('');
}
