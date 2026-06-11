using System.Text;

public class Solution {
    public string int_to_roman(long num) {
        var mapping = new (long val, string sym)[] {
            (1000, "M"), (900, "CM"), (500, "D"), (400, "CD"),
            (100, "C"), (90, "XC"), (50, "L"), (40, "XL"),
            (10, "X"), (9, "IX"), (5, "V"), (4, "IV"), (1, "I")
        };
        var sb = new StringBuilder();
        foreach (var (value, numeral) in mapping) {
            while (num >= value) {
                sb.Append(numeral);
                num -= value;
            }
        }
        return sb.ToString();
    }
}
