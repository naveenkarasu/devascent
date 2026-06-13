function net_salary(gross) {
    let tax = 0.0;
    if (gross <= 300000) {
        tax = 0.0;
    } else if (gross <= 500000) {
        tax = (gross - 300000) * 0.05;
    } else if (gross <= 1000000) {
        tax = 200000 * 0.05 + (gross - 500000) * 0.20;
    } else {
        tax = 200000 * 0.05 + 500000 * 0.20 + (gross - 1000000) * 0.30;
    }
    return Math.round((gross - tax) * 100) / 100;
}
