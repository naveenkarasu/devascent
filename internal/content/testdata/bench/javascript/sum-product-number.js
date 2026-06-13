function is_sum_product_number(n) {
  let digitSum = 0;
  let digitProd = 1;
  let temp = n;
  while (temp > 0) {
    const d = temp % 10;
    digitSum += d;
    digitProd *= d;
    temp = Math.floor(temp / 10);
  }
  return digitSum * digitProd === n;
}
