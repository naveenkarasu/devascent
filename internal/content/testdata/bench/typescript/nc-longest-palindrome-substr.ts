function longest_palindrome(s: string): string {
    let resStart = 0, resLen = 1;

    function expand(left: number, right: number): void {
        while (left >= 0 && right < s.length && s[left] === s[right]) {
            const length = right - left + 1;
            if (length > resLen) {
                resStart = left;
                resLen = length;
            }
            left--;
            right++;
        }
    }

    for (let i = 0; i < s.length; i++) {
        expand(i, i);       // odd-length
        expand(i, i + 1);   // even-length
    }

    return s.substring(resStart, resStart + resLen);
}
