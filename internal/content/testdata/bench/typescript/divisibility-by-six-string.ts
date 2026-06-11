function is_divisible_by_six(s: string): boolean {
    const lastDigit = parseInt(s[s.length - 1]);
    if (lastDigit % 2 !== 0) return false;
    const digitSum = s.split('').reduce((acc, ch) => acc + parseInt(ch), 0);
    return digitSum % 3 === 0;
}
