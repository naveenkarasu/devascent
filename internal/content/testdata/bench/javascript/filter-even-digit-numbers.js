function filter_even_digit_numbers(nums) {
    function digitCount(n) {
        if (n === 0) return 1;
        let count = 0;
        n = Math.abs(n);
        while (n > 0) {
            count++;
            n = Math.floor(n / 10);
        }
        return count;
    }
    return nums.filter(x => digitCount(x) % 2 === 0);
}
