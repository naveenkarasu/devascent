function is_palindrome(s: string): boolean {
    let i = 0, j = s.length - 1;
    while (i < j) {
        while (i < j && !/[a-zA-Z0-9]/.test(s[i])) i++;
        while (i < j && !/[a-zA-Z0-9]/.test(s[j])) j--;
        if (s[i].toLowerCase() !== s[j].toLowerCase()) return false;
        i++;
        j--;
    }
    return true;
}
