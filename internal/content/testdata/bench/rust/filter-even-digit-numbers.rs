fn filter_even_digit_numbers(nums: Vec<i64>) -> Vec<i64> {
    fn digit_count(mut n: i64) -> i64 {
        if n == 0 {
            return 1;
        }
        let mut count = 0;
        while n > 0 {
            count += 1;
            n /= 10;
        }
        count
    }
    nums.into_iter().filter(|&x| digit_count(x) % 2 == 0).collect()
}
