import java.util.*;

class Solution {
    public String insert_to_minimize(String a) {
        if (a.charAt(0) == '1') {
            return a.charAt(0) + "0" + a.substring(1);
        }
        return "1" + a;
    }
}
