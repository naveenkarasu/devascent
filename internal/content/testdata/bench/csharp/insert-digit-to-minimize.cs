public class Solution {
    public string insert_to_minimize(string a) {
        if (a[0] == '1') {
            return a[0] + "0" + a.Substring(1);
        }
        return "1" + a;
    }
}
