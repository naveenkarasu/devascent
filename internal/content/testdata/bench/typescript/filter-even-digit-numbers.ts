function filter_even_digit_numbers(nums: number[]): number[] {
    function digit_count(n: number): number {
        if (n === 0) return 1;
        let count = 0;
        while (n > 0) {
            count++;
            n = Math.floor(n / 10);
        }
        return count;
    }
    return nums.filter(x => digit_count(x) % 2 === 0);
}
