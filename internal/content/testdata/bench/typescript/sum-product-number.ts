function is_sum_product_number(n: number): boolean {
    let digit_sum = 0;
    let digit_prod = 1;
    let temp = n;
    while (temp > 0) {
        const d = temp % 10;
        digit_sum += d;
        digit_prod *= d;
        temp = Math.floor(temp / 10);
    }
    return digit_sum * digit_prod === n;
}
