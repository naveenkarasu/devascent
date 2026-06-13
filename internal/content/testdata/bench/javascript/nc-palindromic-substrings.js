function count_substrings(s) {
    let count = 0;

    function expand(left, right) {
        while (left >= 0 && right < s.length && s[left] === s[right]) {
            count++;
            left--;
            right++;
        }
    }

    for (let i = 0; i < s.length; i++) {
        expand(i, i);     // odd-length
        expand(i, i + 1); // even-length
    }

    return count;
}
